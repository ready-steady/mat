package tgff

import (
	"strings"
	"testing"
)

func assertLexerAt(lexer *lexer, char byte, t *testing.T) {
	if c, _ := lexer.reader.ReadByte(); c != char {
		t.Fatalf("at '%v' instead of '%v'", char)
	} else {
		lexer.reader.UnreadByte()
	}
}

func assertTokens(actual, expected []token, t *testing.T) {
	if len(actual) != len(expected) {
		goto error
	}

	for i := range actual {
		if actual[i].kind != expected[i].kind || actual[i].value != expected[i].value {
			goto error
		}
	}

	return

error:
	t.Fatalf("got %v instead of %v", actual, expected)
}

func TestLexerReadAny(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abbbaacdefg"))

	err := lexer.readAny("ab")

	assertSuccess(err, t)
	assertLexerAt(lexer, 'c', t)
	assertEqual(lexer.value(), "abbbaa", t)
}

func TestLexerReadName(t *testing.T) {
	scenarios := []struct {
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

func TestLexerSkipAny(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abbbaacdefg"))

	err := lexer.skipAny("ab")

	assertSuccess(err, t)
	assertLexerAt(lexer, 'c', t)
}

func TestLexerSkipSequence(t *testing.T) {
	lexer, _ := newLexer(strings.NewReader("abcde"))

	err := lexer.skipSequence("ab")

	assertSuccess(err, t)
	assertLexerAt(lexer, 'c', t)

	err = lexer.skipSequence("d")

	assertFailure(err, t)
}

func TestLexerRun(t *testing.T) {
	lexer, stream := newLexer(strings.NewReader("   \n\n @abcd   42"))

	go lexer.run()

	tokens := []token{}
	for token := range stream {
		tokens = append(tokens, token)
	}

	assertTokens(tokens, []token{
		token{controlToken, "abcd", nil},
		token{numberToken, "42", nil},
	}, t)
}
