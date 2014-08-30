package tgff

import (
	"io"
)

func Parse(reader io.Reader) *Result {
	lexer, stream := newLexer(reader)
	parser, done := newParser(stream)

	go lexer.run()
	go parser.run()

	return <- done
}
