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

	r, err := Parse(file)

	assertSuccess(err, t)

	assertEqual(r.HyperPeriod, uint32(1180), t)
	assertEqual(len(r.Graphs), 5, t)
	assertEqual(len(r.Tables), 3, t)

	graphs := []struct {
		period    uint32
		tasks     int
		arcs      int
		deadlines int
	}{
		{590, 12, 19, 1},
		{1180, 20, 25, 6},
		{1180, 24, 28, 8},
		{590, 8, 7, 3},
		{1180, 20, 24, 6},
	}

	for i, graph := range graphs {
		assertEqual(r.Graphs[i].Name, "TASK_GRAPH", t)
		assertEqual(r.Graphs[i].Number, uint32(i), t)
		assertEqual(r.Graphs[i].Period, graph.period, t)

		assertEqual(len(r.Graphs[i].Tasks), graph.tasks, t)
		assertEqual(len(r.Graphs[i].Arcs), graph.arcs, t)
		assertEqual(len(r.Graphs[i].Deadlines), graph.deadlines, t)
	}

	tables := []struct {
		price float64
	}{
		{70.1121},
		{71.4235},
		{80.491},
	}

	for i, table := range tables {
		assertEqual(r.Tables[i].Name, "COMMUN", t)
		assertEqual(r.Tables[i].Number, uint32(i), t)
		assertEqual(r.Tables[i].Attributes["price"], table.price, t)
	}
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
