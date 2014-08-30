package tgff

import (
	"testing"
	"strings"
)

func assertAt(lexer *lexer, char byte, t *testing.T) {
	if c, _ := lexer.reader.ReadByte(); c != char {
		t.Fatalf("at '%v' instead of '%v'", char)
	} else {
		lexer.reader.UnreadByte()
	}
}

func TestReadChars(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abbbaacdefg"))

	chars, err := lexer.readChars('a', 'b')

	assertSuccess(err, t)
	assertAt(lexer, 'c', t)
	assertEqual(chars, "abbbaa", t)
}

func TestReadName(t *testing.T) {
	scenarios := []struct{
		data string
		name string
	}{
		{"abcd efgh", "abcd"},
	}

	for _, s := range scenarios {
		lexer, _ := newLexer(strings.NewReader(s.data))
		name, _ := lexer.readName()
		assertEqual(name, s.name, t)
	}
}

func TestAccept(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abbbaacdefg"))

	err := lexer.accept('a', 'b')

	assertSuccess(err, t)
	assertAt(lexer, 'c', t)
}

func TestAcceptWhitespace(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("  \t  \t  \n  abc"))

	err := lexer.acceptWhitespace()

	assertSuccess(err, t)
	assertAt(lexer, 'a', t)
}

func TestRequireChar(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abcde"))

	err := lexer.requireChar('a')

	assertSuccess(err, t)
	assertAt(lexer, 'b', t)

	err = lexer.requireChar('c')

	assertFailure(err, t)
}
