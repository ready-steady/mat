package tgff

const (
	taskBufferCap = 50
	arcBufferCap  = 50
)

type Result struct {
	Period uint32

	Graphs []Graph
	Tables []Table
}

type Graph struct {
	Name      string
	Number    uint32
	Period    uint32
	Tasks     []Task
	Arcs      []Arc
	Deadlines []Deadline
}

type Task struct {
	Name string
	Type uint32
}

type Arc struct {
	Name string
	From string
	To   string
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
	Columns    []Column
}

type Column struct {
	Name string
	Data []float64
}

func (r *Result) addGraph(name string, number uint32) *Graph {
	r.Graphs = append(r.Graphs, Graph{
		Name:   name,
		Number: number,
		Tasks:  make([]Task, 0, taskBufferCap),
		Arcs:   make([]Arc, 0, arcBufferCap),
	})

	return &r.Graphs[len(r.Graphs)-1]
}

func (r *Result) addTable(name string, number uint32) *Table {
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
