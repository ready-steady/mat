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

	err := lexer.readChars("ab")

	assertSuccess(err, t)
	assertAt(lexer, 'c', t)
	assertEqual(lexer.flush(), "abbbaa", t)
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
		_ = lexer.readName()
		assertEqual(lexer.flush(), s.name, t)
	}
}

func TestSkip(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abbbaacdefg"))

	err := lexer.skip("ab")

	assertSuccess(err, t)
	assertAt(lexer, 'c', t)
}

func TestSkipWhitespace(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("  \t  \t  \n  abc"))

	err := lexer.skipWhitespace()

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
