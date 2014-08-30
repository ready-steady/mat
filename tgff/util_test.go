package tgff

import (
	"testing"
)

func TestIsNamely(t *testing.T) {
	scenarios := []struct{
		char byte
		result bool
	}{
		{'A', true},
		{'a', true},
		{'0', true},
		{'_', true},
		{'#', false},
	}

	for _, s := range scenarios {
		if isNamely(s.char) != s.result {
			t.Errorf("%v failed", s)
		}
	}
}
