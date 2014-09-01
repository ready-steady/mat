package tgff

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	assertEqual(result.HyperPeriod, uint32(1180), t)
	assertEqual(len(result.Graphs), 5, t)
	assertEqual(len(result.Tables), 3, t)
}

func TestParseFailure(t *testing.T) {
	reader := strings.NewReader("  @ garbage")

	_, err := Parse(reader)

	assertFailure(err, t)
}

func BenchmarkParse(b *testing.B) {
	data := readFixture("simple")

	for i := 0; i < b.N; i++ {
		Parse(bytes.NewReader(data))
	}
}

func readFixture(name string) []byte {
	path := path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
	data, _ := ioutil.ReadFile(path)

	return data
}

func openFixture(name string) *os.File {
	path := path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
	file, _ := os.Open(path)

	return file
}
