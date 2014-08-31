package tgff

import (
	"io"
)

func Parse(reader io.Reader) (*Result, error) {
	done := make(chan bool, 2)

	lexer, stream := newLexer(reader, done)
	parser, success, failure := newParser(stream, done)

	go lexer.run()
	go parser.run()

	select {
	case result := <-success:
		return result, nil
	case err := <-failure:
		return nil, err
	}
}
