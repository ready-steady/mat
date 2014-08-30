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

func assertToken(actual, expected token, t *testing.T) {
	if actual.kind != expected.kind || actual.value != expected.value {
		t.Fatalf("got %v instead of %v", actual, expected)
	}
}

func TestReadChars(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abbbaacdefg"))

	err := lexer.readChars("ab")

	assertSuccess(err, t)
	assertAt(lexer, 'c', t)
	assertEqual(lexer.value(), "abbbaa", t)
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
		assertEqual(lexer.value(), s.name, t)
	}
}

func TestSkipChars(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abbbaacdefg"))

	err := lexer.skipChars("ab")

	assertSuccess(err, t)
	assertAt(lexer, 'c', t)
}

func TestSkipWhitespace(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("  \t  \t  \n  abc"))

	err := lexer.skipWhitespace()

	assertSuccess(err, t)
	assertAt(lexer, 'a', t)
}

func TestSkipChar(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abcde"))

	err := lexer.skipChar('a')

	assertSuccess(err, t)
	assertAt(lexer, 'b', t)

	err = lexer.skipChar('c')

	assertFailure(err, t)
}

func TestRun(t *testing.T) {
	lexer, stream := newLexer(strings.NewReader("   \n\n @abcd   42"))

	go lexer.run()

	tokens := []token{}
	for token := range stream {
		tokens = append(tokens, token)
	}

	assertEqual(len(tokens), 2, t)
	assertToken(tokens[0], token{controlToken, "abcd", nil}, t)
	assertToken(tokens[1], token{numberToken, "42", nil}, t)
}
