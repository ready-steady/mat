package tgff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/goesd/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestParseSuccess(t *testing.T) {
	file := openFixture("simple")
	defer file.Close()

	r, err := Parse(file)

	assert.Success(err, t)

	assert.Equal(r.Period, uint32(1180), t)
	assert.Equal(len(r.Graphs), 5, t)
	assert.Equal(len(r.Tables), 3, t)

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
		assert.Equal(r.Graphs[i].Name, "TASK_GRAPH", t)
		assert.Equal(r.Graphs[i].Number, uint32(i), t)
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
		assert.Equal(r.Tables[i].Number, uint32(i), t)
		assert.Equal(r.Tables[i].Attributes["price"], table.price, t)

		assert.Equal(len(r.Tables[i].Columns), 2, t)
		assert.Equal(r.Tables[i].Columns[0], "type", t)
		assert.Equal(r.Tables[i].Columns[1], "exec_time", t)
	}

	assert.DeepEqual(r.Tables[2].Data, fixtureSimpleTableData2, t)
}

func TestParseFailure(t *testing.T) {
	reader := strings.NewReader("  @ garbage")

	_, err := Parse(reader)

	assert.Failure(err, t)
}

func BenchmarkParseSimple(b *testing.B) {
	data := readFixture("simple")

	for i := 0; i < b.N; i++ {
		Parse(bytes.NewReader(data))
	}
}

func BenchmarkParseComplex(b *testing.B) {
	data := readFixture("complex")

	for i := 0; i < b.N; i++ {
		Parse(bytes.NewReader(data))
	}
}

func ExampleParse() {
	file, _ := os.Open("fixtures/simple.tgff")
	defer file.Close()

	result, _ := Parse(file)

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

var fixtureSimpleTableData2 = []float64{
	+0, 48.5893,
	+1, 33.4384,
	+2, 34.2468,
	+3, 51.2027,
	+4, 51.3571,
	+5, 30.3827,
	+6, 43.3982,
	+7, 60.9097,
	+8, 36.0322,
	+9, 34.7446,
	10, 45.3479,
	11, 31.7221,
	12, 49.6842,
	13, 52.0635,
	14, 44.7690,
	15, 37.7183,
	16, 54.7523,
	17, 58.4432,
	18, 33.1266,
	19, 48.2143,
	20, 31.2946,
	21, 45.9168,
	22, 36.4521,
	23, 61.6448,
	24, 49.4966,
	25, 37.1130,
	26, 40.1642,
	27, 38.9454,
	28, 41.6213,
	29, 42.1084,
	30, 42.4186,
	31, 42.5145,
	32, 34.4180,
	33, 33.4178,
	34, 32.4243,
	35, 63.7925,
	36, 50.3810,
	37, 51.9030,
	38, 46.4714,
	39, 35.0566,
	40, 41.8399,
	41, 30.1513,
	42, 31.7449,
	43, 57.3263,
	44, 61.2321,
	45, 44.9932,
	46, 32.0830,
	47, 37.9489,
	48, 62.4774,
	49, 39.2500,
}
