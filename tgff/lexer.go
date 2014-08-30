package tgff

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const (
	bufferCapacity = 50
)

var whitespaceChars = []byte{' ', '\t', '\n', '\r'}

type lexer struct {
	reader *bufio.Reader
	stream chan token
}

func newLexer(reader io.Reader) (*lexer, <-chan token) {
	stream := make(chan token)

	lexer := &lexer{
		reader: bufio.NewReader(reader),
		stream: stream,
	}

	return lexer, stream
}

func (l *lexer) run() {
	for state := stripState; state != nil; {
		state = state(l)
	}
	close(l.stream)
}

func (l *lexer) read(accept func(uint, byte) bool) (string, error) {
	buffer := make([]byte, 0, bufferCapacity)

	for {
		c, err := l.reader.ReadByte()

		if err != nil {
			return "", err
		}

		if !accept(uint(len(buffer)), c) {
			l.reader.UnreadByte()
			break
		}

		buffer = append(buffer, c)
	}

	return string(buffer), nil
}

func (l *lexer) readChars(chars ...byte) (string, error) {
	return l.read(func(_ uint, c byte) bool {
		return isMember(c, chars)
	})
}

func (l *lexer) readName() (string, error) {
	return l.read(func(i uint, c byte) bool {
		if i == 0 {
			return isAlpha(c)
		} else {
			return isNamely(c)
		}
	})
}

func (l *lexer) accept(chars ...byte) error {
	_, err := l.readChars(chars...)
	return err
}

func (l *lexer) acceptWhitespace() error {
	_, err := l.readChars(whitespaceChars...)
	return err
}

func (l *lexer) requireChar(char byte) error {
	c, err := l.reader.ReadByte()

	if err != nil {
		return err
	}

	if c != char {
		return errors.New(fmt.Sprintf("got %v instead of %v", c, char))
	}

	return nil
}

func (l *lexer) emit(token token) {
	l.stream <- token
}
