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

func (p *parser) pop() *token {
	size := len(p.buffer)

	token := p.buffer[size-1]
	p.buffer = p.buffer[0 : size-1]

	return token
}

func (p *parser) discard() {
	_ = p.pop()
}

func (p *parser) unreceive() {
	p.unreceived <- p.pop()
}

func (p *parser) receive(accept func(*token) bool) (*token, error) {
	var token *token

	select {
	case token = <-p.unreceived:
	default:
		select {
		case token = <-p.stream:
		case <-p.done:
			return nil, io.EOF
		}
	}

	if token.kind == errorToken {
		return nil, errors.New(token.value)
	}

	if !accept(token) {
		return nil, errors.New(fmt.Sprintf("rejected %v", token))
	}

	p.buffer = append(p.buffer, token)

	return token, nil
}

func (p *parser) receiveOne(kind tokenKind) (*token, error) {
	return p.receive(func(token *token) bool {
		return token.kind == kind
	})
}

func (p *parser) receiveOneWith(kind tokenKind, value string) (*token, error) {
	return p.receive(func(token *token) bool {
		return token.kind == kind && token.value == value
	})
}

func (p *parser) receiveOneOf(kinds ...tokenKind) (*token, error) {
	return p.receive(func(token *token) bool {
		for _, kind := range kinds {
			if token.kind == kind {
				return true
			}
		}
		return false
	})
}

func (p *parser) peekOneOf(kinds ...tokenKind) (*token, error) {
	if token, err := p.receiveOneOf(kinds...); err != nil {
		return nil, err
	} else {
		p.unreceive()
		return token, nil
	}
}
