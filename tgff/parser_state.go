package tgff

import (
	"strconv"
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

func parControlState(p *parser) parState {
	control, err := p.receive(controlToken)

	if err != nil {
		return parErrorState(err)
	}

	number, err := p.receive(numberToken)

	if err != nil {
		return parErrorState(err)
	}

	switch control.value {
	case "HYPERPERIOD":
		value, _ := strconv.ParseUint(number.value, 10, 32)
		p.result.hyperPeriod = uint(value)
	}

	return parCompleteState
}
