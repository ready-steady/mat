package tgff

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
)

type lexer struct {
	reader *bufio.Reader
}

var whitespaceList = []byte{' ', '\t', '\n', '\r'}

func newLexer(reader io.Reader) *lexer {
	return &lexer{
		reader: bufio.NewReader(reader),
	}
}

func (l *lexer) accept(chars ...byte) error {
	for {
		c, err := l.reader.ReadByte()

		if err != nil {
			return err
		}

		if !isMember(c, chars) {
			l.reader.UnreadByte()
			return nil
		}
	}

	return nil
}

func (l *lexer) acceptWhitespace() error {
	return l.accept(whitespaceList...)
}

func (l *lexer) require(char byte) error {
	c, err := l.reader.ReadByte()

	if err != nil {
		return err
	}

	if c != char {
		return errors.New(fmt.Sprintf("got %v instead of %v", c, char))
	}

	return nil
}
