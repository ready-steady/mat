package tgff

type Result struct {
	HyperPeriod uint

	Graphs []*Graph
	Tables []*Table
}

type Graph struct {
	Name      string
	Number    uint
	Period    uint
	Tasks     []*Task
	Arcs      []*Arc
	Deadlines []*Deadline
}

type Arc struct {
	Name string
	From string
	To   string
	Type uint
}

type Task struct {
	Name string
	Type uint
}

type Deadline struct {
	Name string
	On   string
	At   uint
}

type Table struct {
	Name   string
	Number uint
}
