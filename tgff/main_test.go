package tgff

import (
	"fmt"
	"os"
	"path"
	"testing"
)

const (
	fixturePath = "fixtures"
)

func TestParse(t *testing.T) {
	file := openFixture("simple")
	defer file.Close()

	result := Parse(file)

	assertEqual(result.graphCount, 5, t)
	assertEqual(result.tableCount, 3, t)
}

func openFixture(name string) *os.File {
	path := path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
	file, _ := os.Open(path)

	return file
}
