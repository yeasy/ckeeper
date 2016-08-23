package engine

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/yeasy/ckeeper/util"
	"strings"
	"bytes"
)

var logger = logging.MustGetLogger("ckeeper.engine")

type Handler struct {
	ruleset    *RuleSet
	containers []docker.APIContainers
}

func NewHanlder() *Handler {
	handler := Handler{}
	handler.ruleset = NewRuleSet()
	return &handler
}

// Load will read in all rules from config
func (h *Handler) Load() {
	r := make(map[string]interface{})
	for name, _ := range viper.GetStringMap("rules") {
		err := viper.UnmarshalKey("rules."+name, &r)
		if err != nil {
			logger.Errorf("unable to decode into struct, %+v", err)
			return
		}
		filter := r["filter"].(map[interface{}]interface{})
		options := docker.ListContainersOptions{}
		for k, v := range filter {
			util.SetField(&options, k.(string), v)
		}
		condition := r["condition"].(string)
		action := r["action"].(string)
		h.ruleset.AddRule(name, Rule{options, condition, action})
	}
	logger.Debugf("Loaded %d rules...", len(h.ruleset.GetRules()))
	logger.Debugf("%+v", h.ruleset.GetRules())
}

func (h *Handler) Process() {
	endpoint := viper.GetString("host.daemon")
	client, _ := docker.NewClient(endpoint)

	done := make(chan int)
	// Process each rule
	for name, r := range h.ruleset.GetRules() {
		logger.Debugf("Processing rule %s:%+v", name, r)
		h.containers, _ = client.ListContainers(r.filter)

		for _, c := range h.containers {
			go execCmd(client, c, []string{r.condition}, done)
			logger.Infof("Run %s on container %s", r.condition, c.ID)
		}

		number := 0
		for s := range done {
			logger.Debugf("cmd resp=%+v", s)
			number++
			if number >= len(h.containers) {
				break
			}
		}
	}
}


func execCmd(client *docker.Client, container docker.APIContainers, cmd []string, done chan int) error {
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
		RawTerminal:  true,
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
