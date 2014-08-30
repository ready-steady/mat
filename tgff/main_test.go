package tgff

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

const (
	fixturePath = "fixtures"
)

func TestParseSuccess(t *testing.T) {
	file := openFixture("simple")
	defer file.Close()

	result, err := Parse(file)

	assertSuccess(err, t)
	assertEqual(result.graphCount, 5, t)
	assertEqual(result.tableCount, 3, t)
}

func TestParseFailure(t *testing.T) {
	reader := strings.NewReader("garbage")

	_, err := Parse(reader)

	assertFailure(err, t)
}

func openFixture(name string) *os.File {
	path := path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
	file, _ := os.Open(path)

	return file
}
