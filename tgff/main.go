package tgff

import (
	"io"
)

func Parse(reader io.Reader) (*Result, error) {
	lexer, stream := newLexer(reader)
	parser, success, failure := newParser(stream)

	go lexer.run()
	go parser.run()

	select {
	case result := <-success:
		return result, nil
	case err := <-failure:
		return nil, err
	}
}
