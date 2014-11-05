package tgff

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const (
	lexStreamCap = 0
	lexBufferCap = 42
)

type lexer struct {
	reader *bufio.Reader
	stream chan<- *token
	abort  <-chan bool
	buffer []byte
}

func newLexer(reader io.Reader, abort <-chan bool) (*lexer, <-chan *token) {
	stream := make(chan *token, lexStreamCap)

	lexer := &lexer{
		reader: bufio.NewReader(reader),
		stream: stream,
		abort:  abort,
		buffer: make([]byte, 0, lexBufferCap),
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
	l.buffer = l.buffer[:0]
}

func (l *lexer) set(value string) {
	l.buffer = append(l.buffer[:0], []byte(value)...)
}

func (l *lexer) peek() (byte, error) {
	c, err := l.reader.ReadByte()

	if err == nil {
		l.reader.UnreadByte()
	}

	return c, err
}

func (l *lexer) read(accept func(int, byte) bool) error {
	k := 0

	for {
		c, err := l.reader.ReadByte()

		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}

		if !accept(k, c) {
			l.reader.UnreadByte()
			return nil
		}

		l.buffer = append(l.buffer, c)
		k++
	}
}

func (l *lexer) readSomething(accept func(int, byte) bool) error {
	size := len(l.buffer)

	if err := l.read(accept); err != nil {
		return err
	}

	if size == len(l.buffer) {
		return errors.New("expected to read something")
	}

	return nil
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

func (l *lexer) readAnyOneOf(chars string) error {
	err := l.read(func(i int, c byte) bool {
		return i == 0 && isMember(chars, c)
	})

	if err != nil {
		return err
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

	if len(l.buffer) != size+count {
		return errors.New(fmt.Sprintf("expected '%v'", chars))
	}

	return nil
}

func (l *lexer) readOne(char byte) error {
	return l.readSequence(string(char))
}

func (l *lexer) readOneOf(chars string) error {
	return l.readSomething(func(i int, c byte) bool {
		return i == 0 && isMember(chars, c)
	})
}

func (l *lexer) readIdent() error {
	return l.readSomething(func(_ int, c byte) bool {
		return isIdently(c)
	})
}

func (l *lexer) readName() error {
	return l.readSomething(func(i int, c byte) bool {
		return isNamely(c)
	})
}

func (l *lexer) skip(yield func() error) error {
	size := len(l.buffer)

	err := yield()

	if err == nil {
		l.buffer = l.buffer[0:size]
	}

	return err
}

func (l *lexer) skipAny(groups ...string) error {
	return l.skip(func() error {
		return l.readAny(groups...)
	})
}

func (l *lexer) skipSequence(chars string) error {
	return l.skip(func() error {
		return l.readSequence(chars)
	})
}

func (l *lexer) skipOne(char byte) error {
	return l.skipSequence(string(char))
}

func (l *lexer) skipLine() error {
	return l.skip(func() error {
		return l.read(func(i int, c byte) bool {
			return !isMember(newLine, c)
		})
	})
}

func (l *lexer) send(kind tokenKind) bool {
	select {
	case l.stream <- &token{kind, l.value()}:
		l.flush()
		return true
	case <-l.abort:
		return false
	}
}
