package tgff

const (
	taskBufferCap = 50
	arcBufferCap  = 50
)

// Result is a representation of a TGFF file.
type Result struct {
	Period uint16 // Hyperperiod

	Graphs []Graph
	Tables []Table
}

// Graph represents a graph in a TGFF file.
type Graph struct {
	Name      string
	Number    uint16
	Period    uint16
	Tasks     []Task
	Arcs      []Arc
	Deadlines []Deadline // Hard deadlines
}

// Task is a TASK entry of a graph.
type Task struct {
	Name string
	Type uint16
}

// Arc is a ARC entry of a graph.
type Arc struct {
	Name string
	From string
	To   string
	Type uint16
}

// Deadline is a HARD_DEADLINE entry of a graph.
type Deadline struct {
	Name string
	On   string
	At   uint16
}

// Table represents a table in a TGFF file.
type Table struct {
	Name       string
	Number     uint16
	Attributes map[string]float64
	Columns    []Column
}

// Column is a column of a table.
type Column struct {
	Name string
	Data []float64
}

func (r *Result) addGraph(name string, number uint16) *Graph {
	r.Graphs = append(r.Graphs, Graph{
		Name:   name,
		Number: number,
		Tasks:  make([]Task, 0, taskBufferCap),
		Arcs:   make([]Arc, 0, arcBufferCap),
	})

	return &r.Graphs[len(r.Graphs)-1]
}

func (r *Result) addTable(name string, number uint16) *Table {
	r.Tables = append(r.Tables, Table{
		Name:   name,
		Number: number,
	})

	return &r.Tables[len(r.Tables)-1]
}

func (g *Graph) addTask() *Task {
	size := len(g.Tasks)

	if size == cap(g.Tasks) {
		temp := make([]Task, 2*size)
		copy(temp, g.Tasks)
		g.Tasks = temp
	}

	g.Tasks = g.Tasks[:size+1]

	return &g.Tasks[size]
}

func (g *Graph) addArc() *Arc {
	size := len(g.Arcs)

	if size == cap(g.Arcs) {
		temp := make([]Arc, 2*size)
		copy(temp, g.Arcs)
		g.Arcs = temp
	}

	g.Arcs = g.Arcs[:size+1]

	return &g.Arcs[size]
}

func (g *Graph) addDeadline() *Deadline {
	g.Deadlines = append(g.Deadlines, Deadline{})

	return &g.Deadlines[len(g.Deadlines)-1]
}
