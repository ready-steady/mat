package tgff

import (
	"reflect"
	"testing"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if actual != expected {
		t.Fatalf("got %v (%T) instead of %v (%T)", actual, actual, expected, expected)
	}
}

func assertEqualSlices(actual, expected interface{}, t *testing.T) {
	if reflect.DeepEqual(actual, expected) {
		t.Fatalf("got %v instead of %v", actual, expected)
	}
}

func assertSuccess(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("got an error '%v'", err)
	}
}

func assertFailure(err error, t *testing.T) {
	if err == nil {
		t.Fatalf("expected an error")
	}
}
