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

func parControlState(p *parser) parState {
	name, err := p.receiveOne(controlToken)
	if err != nil {
		return parEndOrErrorState(err)
	}

	number, err := p.receiveOne(numberToken)
	if err != nil {
		return parErrorState(err)
	}

	if name.value == "HYPERPERIOD" {
		p.result.HyperPeriod = number.Uint32()

		return parControlState
	}

	if _, err := p.receiveOne(blockOpenToken); err != nil {
		return parErrorState(err)
	}

	if token, err := p.peekOneOf(identToken, titleToken); err != nil {
		return parErrorState(err)
	} else if token.kind == identToken {
		graph := &Graph{
			Name:   name.value,
			Number: number.Uint32(),
		}

		p.result.Graphs = append(p.result.Graphs, graph)

		return parGraphState(graph)
	} else {
		table := &Table{
			Name:       name.value,
			Number:     number.Uint32(),
			Attributes: make(map[string]float64, 10),
		}

		p.result.Tables = append(p.result.Tables, table)

		return parTableState(table)
	}
}

// parGraphState processes the body of a task graph declaration including
// PERIOD, TASK, ARC, and HARD_DEADLINE declarations.
func parGraphState(graph *Graph) parState {
	return func(p *parser) parState {
		token, err := p.receiveOneOf(identToken, blockCloseToken)
		if err != nil {
			return parErrorState(err)
		}

		if token.kind == blockCloseToken {
			return parControlState
		}

		switch token.value {
		case "PERIOD":
			if token, err := p.receiveOneOf(numberToken); err != nil {
				return parErrorState(err)
			} else {
				graph.Period = token.Uint32()
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

// parArcState processes an ARC declaration in the following format:
//
//     ARC <name> FROM <task name> TO <task name> TYPE <number as unsigned int>
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
			arc.Type = token.Uint32()
		}

		graph.Arcs = append(graph.Arcs, arc)

		return parGraphState(graph)
	}
}

// parTaskState processes a TASK declaration in the following format:
//
//     TASK <name> TYPE <number as unsigned int>
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
			task.Type = token.Uint32()
		}

		graph.Tasks = append(graph.Tasks, task)

		return parGraphState(graph)
	}
}

// parDeadlineState processes a HARD_DEADLINE declaration in the following
// format:
//
//     HARD_DEADLINE <deadline name> ON <task name> AT <time as unsigned int>
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
			deadline.At = token.Uint32()
		}

		graph.Deadlines = append(graph.Deadlines, deadline)

		return parGraphState(graph)
	}
}

// parTableState processes the body of a table declaration including two
// headers, one for the attributes of the table and one for the actual data,
// and the corresponding content.
func parTableState(table *Table) parState {
	return func(p *parser) parState {
		names, err := p.receiveAny(titleToken)
		if err != nil {
			return parErrorState(err)
		}

		values, err := p.receiveAny(numberToken)
		if err != nil {
			return parErrorState(err)
		}

		if len(names) != len(values) {
			return parErrorState(errors.New(fmt.Sprintf("the attribute header of %v is invalid", table)))
		}

		for i := range names {
			table.Attributes[names[i].value] = values[i].Float64()
		}

		names, err = p.receiveAny(titleToken)
		if err != nil {
			return parErrorState(err)
		}

		cols := len(names)

		values, err = p.receiveAny(numberToken)
		if err != nil {
			return parErrorState(err)
		}

		if len(values) % cols != 0 {
			return parErrorState(errors.New(fmt.Sprintf("the data header of %v is invalid", table)))
		}

		// rows := len(values) / cols

		if _, err := p.receiveOne(blockCloseToken); err != nil {
			return parErrorState(err)
		}

		return parControlState
	}
}
