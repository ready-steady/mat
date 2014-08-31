package tgff

import (
	"testing"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if actual != expected {
		t.Fatalf("got %v (%T) instead of %v (%T)", actual, actual, expected, expected)
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
