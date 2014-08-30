package tgff

import (
	"testing"
)

func TestIsMember(t *testing.T) {
	scenarios := []struct{
		item byte
		set []byte
		result bool
	}{
		{'a', []byte{'a', 'b', 'c'}, true},
		{'b', []byte{'a', 'b', 'c'}, true},
		{'z', []byte{'a', 'b', 'c'}, false},
	}

	for _, s := range scenarios {
		if isMember(s.item, s.set) != s.result {
			t.Errorf("%v failed", s)
		}
	}
}

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
