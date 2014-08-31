package tgff

import (
	"io"
)

type parState func(p *parser) parState

const (
	hyperPeriodName = "HYPERPERIOD"
)

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

func parBlockOpenState(p *parser) parState {
	if err := p.receiveOne(blockOpenToken); err != nil {
		return parErrorState(err)
	}

	if err := p.receiveOneOf(identToken, titleToken); err != nil {
		return parErrorState(err)
	}

	if p.unreceive().kind == identToken {
		return parGraphState
	} else {
		return parTableState
	}
}

func parControlState(p *parser) parState {
	if err := p.receiveOne(controlToken); err != nil {
		return parEndOrErrorState(err)
	}

	if err := p.receiveOne(numberToken); err != nil {
		return parErrorState(err)
	}

	if err := p.receiveOneOf(blockOpenToken, controlToken); err != nil {
		return parErrorState(err)
	}

	if p.unreceive().kind == blockOpenToken {
		return parBlockOpenState
	} else {
		if err := p.commitParameter(); err != nil {
			return parErrorState(err)
		}
		return parControlState
	}
}

func parGraphState(_ *parser) parState {
	return parCompleteState
}

func parTableState(_ *parser) parState {
	return parCompleteState
}
