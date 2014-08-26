package mat

import (
	"testing"
	"reflect"
)

const (
	fixturePath = "fixtures/data.mat"
)

func assertEqual(expected, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Got", actual, "instead of", expected)
	}
}

func TestOpen(t *testing.T) {
	file, _ := Open(fixturePath)
	assertEqual(fixturePath, file.path, t)
}
