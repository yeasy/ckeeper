package util

import (
	"testing"
)

func TestListHasString(t *testing.T) {
	value := "hello"
	list := []string{value, "world"}

	if ListHasString(value, list) != true {
		t.Error("Error to find existing element in list")
	}

	list = []string{"world"}

	if ListHasString(value, list) == true {
		t.Error("Error to find non-existing element in list")
	}
}
