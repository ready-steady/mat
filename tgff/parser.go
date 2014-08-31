package tgff

import (
	"errors"
	"fmt"
)

type parser struct {
	stream  <-chan *token
	done    chan bool
	success chan<- *Result
	failure chan<- error
	result  *Result
}

func newParser(stream <-chan *token, done chan bool) (*parser, <-chan *Result, <-chan error) {
	success := make(chan *Result)
	failure := make(chan error)

	parser := &parser{
		stream:  stream,
		done:    done,
		success: success,
		failure: failure,
		result:  &Result{},
	}

	return parser, success, failure
}

func (p *parser) run() {
	for state := parControlState; state != nil; {
		state = state(p)
	}
	p.done <- true
}

func (p *parser) receive(kind tokenKind) (*token, error) {
	select {
	case token := <-p.stream:
		if token.kind == errorToken {
			return nil, errors.New(fmt.Sprintf("got an error '%v'", token.value))
		}

		if token.kind != kind {
			return nil, errors.New(fmt.Sprintf("got %v instead of %v",
				token.kind, kind))
		}

		return token, nil
	case <-p.done:
		return nil, errors.New(fmt.Sprintf("expected to receive %v", kind))
	}
}
