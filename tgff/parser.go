package tgff

import (
	"errors"
	"fmt"
	"io"
)

const (
	parBufferCapacity = 10
)

type parser struct {
	stream <-chan *token
	done   chan bool

	buffer     []*token
	unreceived chan *token

	result *Result

	success chan<- *Result
	failure chan<- error
}

func newParser(stream <-chan *token, done chan bool) (*parser, <-chan *Result, <-chan error) {
	success := make(chan *Result)
	failure := make(chan error)

	parser := &parser{
		stream: stream,
		done:   done,

		buffer:     make([]*token, 0, parBufferCapacity),
		unreceived: make(chan *token, 1),

		result: &Result{},

		success: success,
		failure: failure,
	}

	return parser, success, failure
}

func (p *parser) run() {
	for state := parControlState; state != nil; {
		state = state(p)
	}
	p.done <- true
}

func (p *parser) last() *token {
	return p.buffer[len(p.buffer)-1]
}

func (p *parser) pop() *token {
	size := len(p.buffer)

	token := p.buffer[size-1]
	p.buffer = p.buffer[0 : size-1]

	return token
}

func (p *parser) receive(accept func(*token) bool) error {
	var token *token

	select {
	case token = <-p.unreceived:
	default:
		select {
		case token = <-p.stream:
		case <-p.done:
			return io.EOF
		}
	}

	if token.kind == errorToken {
		return errors.New(fmt.Sprintf("got an error '%v'", token.value))
	}

	if !accept(token) {
		return errors.New(fmt.Sprintf("rejected %v", token.kind))
	}

	p.buffer = append(p.buffer, token)

	return nil
}

func (p *parser) receiveOne(kind tokenKind) error {
	return p.receive(func(token *token) bool {
		return token.kind == kind
	})
}

func (p *parser) receiveOneOf(kinds ...tokenKind) error {
	return p.receive(func(token *token) bool {
		for _, kind := range kinds {
			if token.kind == kind {
				return true
			}
		}
		return false
	})
}

func (p *parser) unreceive() *token {
	token := p.pop()

	p.unreceived <- token

	return token
}

func (p *parser) commitParameter() error {
	value := p.pop()
	name := p.pop()

	switch name.value {
	case hyperPeriodName:
		p.result.hyperPeriod = value.Uint()
		return nil
	default:
		return errors.New(fmt.Sprintf("unknown parameter '%v'", name.value))
	}
}
