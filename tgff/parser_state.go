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
		p.result.Period = number.Uint()

		return parControlState
	}

	if _, err := p.receiveOne(blockOpenToken); err != nil {
		return parErrorState(err)
	}

	if token, err := p.peekOneOf(identToken, titleToken); err != nil {
		return parErrorState(err)
	} else if token.kind == identToken {
		return parGraphState(p.result.addGraph(name.value, number.Uint()))
	} else {
		return parTableState(p.result.addTable(name.value, number.Uint()))
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

// parTaskState processes a TASK declaration in the following format:
//
//     TASK <name> TYPE <number as unsigned int>
//
// The leading TASK is assumed to be already consumed.
func parTaskState(graph *Graph) parState {
	return func(p *parser) parState {
		task := graph.addTask()

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

		return parGraphState(graph)
	}
}

// parArcState processes an ARC declaration in the following format:
//
//     ARC <name> FROM <task name> TO <task name> TYPE <number as unsigned int>
//
// The leading ARC is assumed to be already consumed.
func parArcState(graph *Graph) parState {
	return func(p *parser) parState {
		arc := graph.addArc()

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
		deadline := graph.addDeadline()

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

		cols := len(names)

		if cols != len(values) {
			return parErrorState(errors.New(fmt.Sprintf(
				"the attribute header of %v %v is invalid", table.Name, table.ID)))
		}

		table.Attributes = make(map[string]float64, cols)
		for i := range names {
			table.Attributes[names[i].value] = values[i].Float64()
		}

		names, err = p.receiveAny(titleToken)
		if err != nil {
			return parErrorState(err)
		}

		cols = len(names)

		values, err = p.receiveAny(numberToken)
		if err != nil {
			return parErrorState(err)
		}

		size := len(values)

		if size%cols != 0 {
			return parErrorState(errors.New(fmt.Sprintf(
				"the data header of %v %v is invalid", table.Name, table.ID)))
		}

		rows := size / cols

		table.Columns = make([]Column, cols)
		for i, name := range names {
			table.Columns[i].Name = name.value
			table.Columns[i].Data = make([]float64, rows)
			for j := 0; j < rows; j++ {
				table.Columns[i].Data[j] = values[j*cols+i].Float64()
			}
		}

		if _, err := p.receiveOne(blockCloseToken); err != nil {
			return parErrorState(err)
		}

		return parControlState
	}
}
