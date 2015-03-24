package tgff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ready-steady/assert"
)

const (
	fixturePath = "fixtures"
)

func TestParseSuccess_simple(t *testing.T) {
	file := openFixture("simple")
	defer file.Close()

	r, err := Parse(file)

	assert.Success(err, t)

	assert.Equal(r.Period, uint(1180), t)
	assert.Equal(len(r.Graphs), 5, t)
	assert.Equal(len(r.Tables), 3, t)

	graphs := []struct {
		period    uint
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
		assert.Equal(r.Graphs[i].Name, "TASK_GRAPH", t)
		assert.Equal(r.Graphs[i].ID, uint(i), t)
		assert.Equal(r.Graphs[i].Period, graph.period, t)

		assert.Equal(len(r.Graphs[i].Tasks), graph.tasks, t)
		assert.Equal(len(r.Graphs[i].Arcs), graph.arcs, t)
		assert.Equal(len(r.Graphs[i].Deadlines), graph.deadlines, t)
	}

	tables := []struct {
		price float64
	}{
		{70.1121},
		{71.4235},
		{80.491},
	}

	for i, table := range tables {
		assert.Equal(r.Tables[i].Name, "COMMUN", t)
		assert.Equal(r.Tables[i].ID, uint(i), t)
		assert.Equal(r.Tables[i].Attributes["price"], table.price, t)

		assert.Equal(len(r.Tables[i].Columns), 2, t)
		assert.Equal(r.Tables[i].Columns[0].Name, "type", t)
		assert.Equal(r.Tables[i].Columns[1].Name, "exec_time", t)
	}

	assert.Equal(r.Tables[1].Columns[1].Data, fixtureSimpleTable1Column1, t)
}

func TestParseSuccess_032_640(t *testing.T) {
	file := openFixture("032_640")
	defer file.Close()

	r, err := Parse(file)

	assert.Success(err, t)

	assert.Equal(len(r.Graphs), 1, t)
	assert.Equal(len(r.Graphs[0].Tasks), 640, t)
	assert.Equal(len(r.Graphs[0].Arcs), 848, t)
	assert.Equal(len(r.Graphs[0].Deadlines), 259, t)

	assert.Equal(len(r.Tables), 32, t)
	for _, table := range r.Tables {
		assert.Equal(len(table.Attributes), 1, t)
		assert.Equal(len(table.Columns), 4, t)
		for _, column := range table.Columns {
			assert.Equal(len(column.Data), 320, t)
		}
	}
}

func TestParseFailure(t *testing.T) {
	reader := strings.NewReader("  @ garbage")

	_, err := Parse(reader)

	assert.Failure(err, t)
}

func BenchmarkParse_simple(b *testing.B) {
	data := readFixture("simple")

	for i := 0; i < b.N; i++ {
		Parse(bytes.NewReader(data))
	}
}

func BenchmarkParse_032_640(b *testing.B) {
	data := readFixture("032_640")

	for i := 0; i < b.N; i++ {
		Parse(bytes.NewReader(data))
	}
}

func ExampleParseFile() {
	result, _ := ParseFile("fixtures/simple.tgff")

	fmt.Println("Task graphs:", len(result.Graphs))
	fmt.Println("Data tables:", len(result.Tables))

	// Output:
	// Task graphs: 5
	// Data tables: 3
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

var fixtureSimpleTable1Column1 = []float64{
	48.5893,
	33.4384,
	34.2468,
	51.2027,
	51.3571,
	30.3827,
	43.3982,
	60.9097,
	36.0322,
	34.7446,
	45.3479,
	31.7221,
	49.6842,
	52.0635,
	44.7690,
	37.7183,
	54.7523,
	58.4432,
	33.1266,
	48.2143,
	31.2946,
	45.9168,
	36.4521,
	61.6448,
	49.4966,
	37.1130,
	40.1642,
	38.9454,
	41.6213,
	42.1084,
	42.4186,
	42.5145,
	34.4180,
	33.4178,
	32.4243,
	63.7925,
	50.3810,
	51.9030,
	46.4714,
	35.0566,
	41.8399,
	30.1513,
	31.7449,
	57.3263,
	61.2321,
	44.9932,
	32.0830,
	37.9489,
	62.4774,
	39.2500,
}
