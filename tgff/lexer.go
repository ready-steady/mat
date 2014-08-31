package tgff

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const (
	bufferCapacity = 42
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
	for state := lexUncertainState; state != nil; {
		state = state(l)
	}
	close(l.stream)
}

func (l *lexer) length() int {
	return len(l.buffer)
}

func (l *lexer) value() string {
	return string(l.buffer)
}

func (l *lexer) flush() {
	l.buffer = l.buffer[0:0]
}

func (l *lexer) set(value string) {
	l.buffer = append(l.buffer[0:0], []byte(value)...)
}

func (l *lexer) peek() (byte, error) {
	c, err := l.reader.ReadByte()

	if err == nil {
		l.reader.UnreadByte()
	}

	return c, err
}

func (l *lexer) read(accept func(int, byte) bool) error {
	for {
		c, err := l.reader.ReadByte()

		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}

		if !accept(len(l.buffer), c) {
			l.reader.UnreadByte()
			return nil
		}

		l.buffer = append(l.buffer, c)
	}
}

func (l *lexer) readAny(groups ...string) error {
	for _, chars := range groups {
		err := l.read(func(_ int, c byte) bool {
			return isMember(chars, c)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (l *lexer) readSequence(chars string) error {
	size := len(l.buffer)
	count := len(chars)

	err := l.read(func(i int, c byte) bool {
		return i < count && chars[i] == c
	})

	if err != nil {
		return err
	}

	if len(l.buffer) != size + count {
		return errors.New(fmt.Sprintf("expected '%v'", chars))
	}

	return nil
}

func (l *lexer) readOne(char byte) error {
	return l.readSequence(string(char))
}

func (l *lexer) readName() error {
	return l.read(func(i int, c byte) bool {
		if i == 0 {
			return isAlpha(c)
		} else {
			return isNamely(c)
		}
	})
}

func (l *lexer) skipAny(groups ...string) error {
	size := len(l.buffer)

	err := l.readAny(groups...)

	if err != nil {
		return err;
	}

	l.buffer = l.buffer[0:size]

	return nil
}

func (l *lexer) skipSequence(chars string) error {
	err := l.readSequence(chars)

	if err != nil {
		return err
	}

	l.buffer = l.buffer[0 : len(l.buffer)-1]

	return nil
}

func (l *lexer) emit(kind tokenKind, more ...interface{}) {
	defer l.flush()
	l.stream <- token{kind, l.value(), more}
}
