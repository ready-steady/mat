package tgff

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	bufferCapacity = 42
	whitespaceChars = " \t\n\r"
)


type lexer struct {
	reader *bufio.Reader
	stream chan token
	buffer []byte
}

func newLexer(reader io.Reader) (*lexer, <-chan token) {
	stream := make(chan token)

	lexer := &lexer{
		reader: bufio.NewReader(reader),
		stream: stream,
		buffer: make([]byte, 0, bufferCapacity),
	}

	return lexer, stream
}

func (l *lexer) run() {
	for state := controlState; state != nil; {
		state = state(l)
	}
	close(l.stream)
}

func (l *lexer) flush() string {
	value := string(l.buffer)

	l.buffer = l.buffer[0:1]

	return value
}

func (l *lexer) read(accept func(uint, byte) bool) error {
	for {
		c, err := l.reader.ReadByte()

		if err != nil {
			return err
		}

		if !accept(uint(len(l.buffer)), c) {
			l.reader.UnreadByte()
			break
		}

		l.buffer = append(l.buffer, c)
	}

	return nil
}

func (l *lexer) readChars(chars string) error {
	return l.read(func(_ uint, c byte) bool {
		return strings.IndexByte(chars, c) >= 0
	})
}

func (l *lexer) readName() error {
	return l.read(func(i uint, c byte) bool {
		if i == 0 {
			return isAlpha(c)
		} else {
			return isNamely(c)
		}
	})
}

func (l *lexer) skip(chars string) error {
	return l.readChars(chars)
}

func (l *lexer) skipWhitespace() error {
	return l.readChars(whitespaceChars)
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

func (l *lexer) emit(kind tokenKind, more ...interface{}) {
	l.stream <- token{kind, l.flush(), more}
}
