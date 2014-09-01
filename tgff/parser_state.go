package tgff

import (
	"errors"
	"fmt"
	"io"
)

type parState func(p *parser) parState

func parCompleteState(p *parser) parState {
	p.success <- p.result

	return nil
}

func parErrorState(err error) parState {
	return func(p *parser) parState {
		p.failure <- err

		return nil
	}
}

func parEndOrErrorState(err error) parState {
	if err == io.EOF {
		return parCompleteState
	} else {
		return parErrorState(err)
	}
}

func parBlockState(p *parser) parState {
	number, name := p.pop(), p.pop()

	if _, err := p.receiveOne(blockOpenToken); err != nil {
		return parErrorState(err)
	}

	if token, err := p.peekOneOf(identToken, titleToken); err != nil {
		return parErrorState(err)
	} else if token.kind == identToken {
		graph := &Graph{
			Name:   name.value,
			Number: number.Uint(),
		}

		p.result.Graphs = append(p.result.Graphs, graph)

		return parGraphState(graph)
	} else {
		table := &Table{
			Name:   name.value,
			Number: number.Uint(),
		}

		p.result.Tables = append(p.result.Tables, table)

		return parTableState(table)
	}
}

func parControlState(p *parser) parState {
	name, err := p.receiveOne(controlToken)
	if err != nil {
		return parEndOrErrorState(err)
	}

	value, err := p.receiveOne(numberToken)
	if err != nil {
		return parErrorState(err)
	}

	if name.value != "HYPERPERIOD" {
		return parBlockState
	}

	p.discard()
	p.discard()

	p.result.HyperPeriod = value.Uint()

	return parControlState
}

func parGraphState(graph *Graph) parState {
	return func(p *parser) parState {
		token, err := p.receiveOneOf(identToken, blockCloseToken)
		if err != nil {
			return parErrorState(err)
		}

		if token.kind == blockCloseToken {
			p.discard()

			return parControlState
		}

		switch token.value {
		case "PERIOD":
			if token, err := p.receiveOneOf(numberToken); err != nil {
				return parErrorState(err)
			} else {
				graph.Period = token.Uint()
				return parGraphState(graph)
			}
		case "ARC":
			return parArcState(graph)
		case "TASK":
			return parTaskState(graph)
		case "HARD_DEADLINE":
			return parDeadlineState(graph)
		default:
			return parErrorState(errors.New(fmt.Sprintf("unexpected %v", token)))
		}
	}
}

// parArcState processes ARC declarations in the following format:
//
//     ARC <name> FROM <task name> TO <task name> TYPE <number as uint>
//
// The leading ARC is assumed to be already consumed.
func parArcState(graph *Graph) parState {
	return func(p *parser) parState {
		arc := &Arc{}

		if token, err := p.receiveOne(nameToken); err != nil {
			return parErrorState(err)
		} else {
			arc.Name = token.value
		}

		if _, err := p.receiveOneWith(identToken, "FROM"); err != nil {
			return parErrorState(err)
		}

		if token, err := p.receiveOne(nameToken); err != nil {
			return parErrorState(err)
		} else {
			arc.From = token.value
		}

		if _, err := p.receiveOneWith(identToken, "TO"); err != nil {
			return parErrorState(err)
		}

		if token, err := p.receiveOne(nameToken); err != nil {
			return parErrorState(err)
		} else {
			arc.To = token.value
		}

		if _, err := p.receiveOneWith(identToken, "TYPE"); err != nil {
			return parErrorState(err)
		}

		if token, err := p.receiveOne(numberToken); err != nil {
			return parErrorState(err)
		} else {
			arc.Type = token.Uint()
		}

		graph.Arcs = append(graph.Arcs, arc)

		return parGraphState(graph)
	}
}

// parTaskState processes with TASK declarations in the following format:
//
//     TASK <name> TYPE <number as uint>
//
// The leading TASK is assumed to be already consumed.
func parTaskState(graph *Graph) parState {
	return func(p *parser) parState {
		task := &Task{}

		if token, err := p.receiveOne(nameToken); err != nil {
			return parErrorState(err)
		} else {
			task.Name = token.value
		}

		if _, err := p.receiveOneWith(identToken, "TYPE"); err != nil {
			return parErrorState(err)
		}

		if token, err := p.receiveOne(numberToken); err != nil {
			return parErrorState(err)
		} else {
			task.Type = token.Uint()
		}

		graph.Tasks = append(graph.Tasks, task)

		return parGraphState(graph)
	}
}

// parDeadlineState processes with HARD_DEADLINE declarations in the following
// format:
//
//     HARD_DEADLINE <deadline name> ON <task name> AT <time as uint>
//
// The leading HARD_DEADLINE is assumed to be already consumed.
func parDeadlineState(graph *Graph) parState {
	return func(p *parser) parState {
		deadline := &Deadline{}

		if token, err := p.receiveOne(nameToken); err != nil {
			return parErrorState(err)
		} else {
			deadline.Name = token.value
		}

		if _, err := p.receiveOneWith(identToken, "ON"); err != nil {
			return parErrorState(err)
		}

		if token, err := p.receiveOne(nameToken); err != nil {
			return parErrorState(err)
		} else {
			deadline.On = token.value
		}

		if _, err := p.receiveOneWith(identToken, "AT"); err != nil {
			return parErrorState(err)
		}

		if token, err := p.receiveOne(numberToken); err != nil {
			return parErrorState(err)
		} else {
			deadline.At = token.Uint()
		}

		graph.Deadlines = append(graph.Deadlines, deadline)

		return parGraphState(graph)
	}
}

func parTableState(table *Table) parState {
	return func(p *parser) parState {
		return parCompleteState
	}
}
