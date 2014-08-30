package tgff

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"
)

const (
	fixturePath = "fixtures"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if actual != expected {
		t.Fatalf("got %v instead of %v", actual, expected)
	}
}

func TestParse(t *testing.T) {
	result := Parse(fixture("simple"))

	assertEqual(result.graphCount, 5, t)
	assertEqual(result.tableCount, 3, t)
}

func fixture(name string) io.Reader {
	path := path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
	file, _ := os.Open(path)
	return file
}
