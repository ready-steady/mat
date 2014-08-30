package tgff

import (
	"testing"
	"strings"
)

func assertSuccess(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("got an error '%v'", err)
	}
}

func assertFailure(err error, t *testing.T) {
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func assertAt(lexer *lexer, char byte, t *testing.T) {
	if c, _ := lexer.reader.ReadByte(); c != char {
		t.Fatalf("at '%v' instead of '%v'", char)
	} else {
		lexer.reader.UnreadByte()
	}
}

func TestAccept(t *testing.T) {
	lexer := newLexer(strings.NewReader("abbbaacdefg"))

	err := lexer.accept('a', 'b')

	assertSuccess(err, t)
	assertAt(lexer, 'c', t)
}

func TestAcceptWhitespace(t *testing.T) {
	lexer := newLexer(strings.NewReader("  \t  \t  \n  abc"))

	err := lexer.acceptWhitespace()

	assertSuccess(err, t)
	assertAt(lexer, 'a', t)
}

func TestRequire(t *testing.T) {
	lexer := newLexer(strings.NewReader("abcde"))

	err := lexer.require('a')

	assertSuccess(err, t)
	assertAt(lexer, 'b', t)

	err = lexer.require('c')

	assertFailure(err, t)
}
