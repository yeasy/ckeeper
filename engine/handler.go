package engine

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/yeasy/ckeeper/util"
)

var logger = logging.MustGetLogger("ckeeper")

// Handler is the processor on the ruleset
type Handler struct {
	ruleset    *RuleSet
	containers []docker.APIContainers
}

// NewHandler return a initialized Handler object
func NewHandler() *Handler {
	handler := Handler{}
	handler.ruleset = NewRuleSet()
	return &handler
}

// Load will read the rules from config
func (h *Handler) Load() error {
	r := Rule{}
	option := docker.ListContainersOptions{}
	for name := range viper.GetStringMap("rules") {
		err := viper.UnmarshalKey("rules."+name+".option", &option)
		if err != nil {
			logger.Errorf("unable to decode into struct, %+v", err)
			return err
		}
		if util.ListHasString("rules."+name+".option.Filters", viper.AllKeys()) {
			option.Filters = viper.GetStringMapStringSlice("rules." + name + ".option.Filters")
		}
		r.name = name
		r.option = option
		r.target = viper.GetString("rules." + name + ".target")
		r.action = viper.GetString("rules." + name + ".action")
		logger.Debugf("%s=%+v", name, r)
		h.ruleset.AddRule(r)
	}
	logger.Infof("Loaded rules: %d", len(h.ruleset.GetRules()))
	logger.Debugf("%+v", h.ruleset.GetRules())
	return nil
}

// Process will check each rule and run it
func (h *Handler) Process() {
	endpoint := viper.GetString("host.daemon")
	client, _ := docker.NewClient(endpoint)
	done := make(chan int)
	triggered := 0

	// Process each rule
	for name, rule := range h.ruleset.GetRules() {
		logger.Debugf("Processing rule %+v", rule)
		h.containers, _ = client.ListContainers(rule.option)
		if len(h.containers) <= 0 {
			logger.Infof("No container matched rule %s", name)
			continue
		}

		for _, container := range h.containers {
			logger.Debugf("Rule %s matched container %s", name, container.ID)
			go execCmd(container, rule, client, done)
		}

		i := 0
		for s := range done {
			triggered += s
			i++
			if i >= len(h.containers) {
				break
			}
		}
	}
	logger.Infof("Process rules: %d, triggered actions: %d", len(h.ruleset.GetRules()), triggered)
}

func execCmd(container docker.APIContainers, rule Rule, client *docker.Client, done chan int) error {
	bashCmd := strings.Replace(rule.target, "CONTAINER", util.GetContainerIP(container), -1)
	if bashCmd != "" {
		for i := 0; i < viper.GetInt("check.retries"); i++ {
			_, err := exec.Command("bash", "-c", bashCmd).Output()
			if err == nil { // target run successfully, just return
				logger.Debugf("Do nothing on container %s", container.ID)
				done <- 0
				return nil
			}
		}
	}
	logger.Debugf("Trigger action %s on container %s", rule.action, container.ID)
	switch rule.action {
	case "restart":
		client.RestartContainer(container.ID, 5)
	case "start":
		client.StartContainer(container.ID, nil)
	}
	done <- 1
	return nil
}

func execCmdT(client *docker.Client, container docker.APIContainers, cmd []string, done chan int) error {
	success := make(chan struct{})
	var stdout, stderr bytes.Buffer
	var reader = strings.NewReader("send value")
	createOptions := docker.CreateExecOptions{
		AttachStdout: true,
		Tty:          true,
		Cmd:          cmd,
		Container:    container.ID,
	}
	startOptions := docker.StartExecOptions{
		OutputStream: &stdout,
		ErrorStream:  &stderr,
		InputStream:  reader,
		//Tty:true,
		RawTerminal: true,
		//Success:      success,
	}
	exec, err := client.CreateExec(createOptions)
	logger.Debug("start exec")
	err = client.StartExec(exec.ID, startOptions)
	if err != nil {
		logger.Errorf("Cannot start exec")
		logger.Error(err)
		done <- -1
		return err
	}
	logger.Debug("read exec")
	v := make([]byte, 5000)
	n, err := stdout.Read(v)
	if err != nil {
		logger.Errorf("Cannot parse cmd read\n")
		logger.Error(err)
		done <- -1
		return err
	}
	logger.Infof("input[%d]=%s, success=%+v", n, v, <-success)
	done <- 0
	return nil
}
