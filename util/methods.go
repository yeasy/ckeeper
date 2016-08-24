package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fsouza/go-dockerclient"
)


// GetContainerIP will return the ip address of the given container
func GetContainerIP(container docker.APIContainers) string {
	return container.Networks.Networks["bridge"].IPAddress
}

// ListHasString will check whether a list has the given string as element
func ListHasString(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
