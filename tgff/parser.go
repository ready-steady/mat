package tgff

import (
	"errors"
	"fmt"
	"io"
)

const (
	parBufferCap = 100
)

type parser struct {
	stream <-chan *token
	buffer chan *token

	abort   chan<- bool
	success chan<- Result
	failure chan<- error

	result Result
}

func newParser(stream <-chan *token, abort chan<- bool) (*parser, <-chan Result, <-chan error) {
	success := make(chan Result)
	failure := make(chan error)

	parser := &parser{
		stream: stream,
		buffer: make(chan *token, 1),

		abort:   abort,
		success: success,
		failure: failure,

		result: Result{},
	}

	return parser, success, failure
}

func (p *parser) run() {
	for state := parControlState; state != nil; {
		state = state(p)
	}
	p.abort <- true
}

func (p *parser) unreceive(token *token) {
	p.buffer <- token
}

func (p *parser) receive(accept func(*token) bool) (*token, error) {
	var token *token
	var ok bool

	select {
	case token = <-p.buffer:
	default:
		if token, ok = <-p.stream; !ok {
			return nil, io.EOF
		}
	}

	if token.kind == errorToken {
		return nil, errors.New(token.value)
	}

	if !accept(token) {
		return nil, errors.New(fmt.Sprintf("rejected %v", token))
	}

	return token, nil
}

func (p *parser) receiveWhile(accept func(*token) bool) ([]*token, error) {
	tokens := make([]*token, 0, parBufferCap)

	extendIfNeeded := func() {
		size := len(tokens)

		if size < cap(tokens) {
			return
		}

		newTokens := make([]*token, size, 2*size)
		copy(newTokens, tokens)
		tokens = newTokens
	}

	var token *token
	var ok bool

	for {
		select {
		case token = <-p.buffer:
		default:
			if token, ok = <-p.stream; !ok {
				return tokens, nil
			}
		}

		if token.kind == errorToken {
			return tokens, errors.New(token.value)
		}

		if !accept(token) {
			p.unreceive(token)
			return tokens, nil
		}

		extendIfNeeded()

		tokens = append(tokens, token)
	}
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
		p.unreceive(token)
		return token, nil
	}
}

func (p *parser) receiveAny(kind tokenKind) ([]*token, error) {
	return p.receiveWhile(func(token *token) bool {
		return token.kind == kind
	})
}
