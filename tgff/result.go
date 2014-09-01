package tgff

import (
	"fmt"
)

type Result struct {
	HyperPeriod uint32

	Graphs []*Graph
	Tables []*Table
}

type Graph struct {
	Name      string
	Number    uint32
	Period    uint32
	Tasks     []*Task
	Arcs      []*Arc
	Deadlines []*Deadline
}

type Arc struct {
	Name string
	From string
	To   string
	Type uint32
}

type Task struct {
	Name string
	Type uint32
}

type Deadline struct {
	Name string
	On   string
	At   uint32
}

type Table struct {
	Name       string
	Number     uint32
	Attributes map[string]float64
	Columns    []string
	Data       []float64
}

func (t *Table) String() string {
	return fmt.Sprintf("%v (%v)", t.Name, t.Number)
}
