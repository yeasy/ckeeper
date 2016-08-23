package util

import (
	"errors"
	"fmt"
	"reflect"
	"github.com/fsouza/go-dockerclient"
)

//SetField will set a obj's variable with given value
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}


func GetContainerIP(container docker.APIContainers) string {
	return container.Networks.Networks["bridge"].IPAddress
}

func ListHasString(value string, list []string) bool {
    for _, v := range list {
        if v == value {
            return true
        }
    }
    return false
}