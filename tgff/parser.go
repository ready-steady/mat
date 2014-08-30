package tgff

import (
	"log"
)

type parser struct {
	stream <-chan token
	done chan<- *Result
}

func newParser(stream <-chan token) (*parser, <-chan *Result) {
	done := make(chan *Result)

	parser := &parser{
		stream: stream,
		done: done,
	}

	return parser, done
}

func (p *parser) run() {
	result := &Result{}

	for token := range p.stream {
		log.Printf("%T: %v\n", token, token)
	}

	p.done <- result
}
