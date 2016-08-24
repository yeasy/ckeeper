package engine

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
)

type FakeDockerClient struct {
}

func (client *FakeDockerClient) ListContainers(interface{}) ([]docker.APIContainers, error) {

	return []docker.APIContainers{
		{
			ID:    "test_container_id",
			Names: []string{"test_container_name"},
		},
	}, nil
}

func TestDockerList(t *testing.T) {
	//endpoint := "unix:///var/run/docker.sock"
	//client, _ := docker.NewClient(endpoint)
	client := FakeDockerClient{}
	ListOption := docker.ListContainersOptions{
		All: true,
	}
	ListOption.Filters = make(map[string][]string)
	ListOption.Filters["status"] = []string{"exited"}
	logger.Debug("Exited Containers")
	containers, _ := client.ListContainers(ListOption)
	for _, container := range containers {
		logger.Debugf("name=%s, id=%s", container.Names[0], container.ID)
	}

	ListOption.Filters["status"] = []string{"running"}
	logger.Debug("Running Containers")
	containers, _ = client.ListContainers(ListOption)
	for _, container := range containers {
		logger.Debugf("name=%s, id=%s", container.Names[0], container.ID)
	}

	ListOption.Filters["status"] = []string{"paused"}
	logger.Debug("Paused Containers")
	containers, _ = client.ListContainers(ListOption)
	for _, container := range containers {
		logger.Debugf("name=%s, id=%s", container.Names[0], container.ID)
	}

	ListOption.Filters["status"] = []string{"paused", "exited"}
	logger.Debug("Paused and Exited Containers")
	containers, _ = client.ListContainers(ListOption)
	for _, container := range containers {
		logger.Debugf("name=%s, id=%s", container.Names[0], container.ID)
	}
}
