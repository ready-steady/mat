package tgff

import (
	"log"
)

type parser struct {
	stream  <-chan token
	success chan<- *Result
	failure chan<- error
}

func newParser(stream <-chan token) (*parser, <-chan *Result, <-chan error) {
	success := make(chan *Result)
	failure := make(chan error)

	parser := &parser{
		stream:  stream,
		success: success,
		failure: failure,
	}

	return parser, success, failure
}

func (p *parser) run() {
	result := &Result{}

	for token := range p.stream {
		log.Printf("%T: %v\n", token, token)

		switch token.kind {
		case errorToken:
			p.failure <- token.more[0].(error)
			return
		}
	}

	p.success <- result
}
