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

func lexerRun(data string) []token {
	lexer, stream := newLexer(strings.NewReader(data))

	go lexer.run()

	tokens := []token{}
	for token := range stream {
		tokens = append(tokens, token)
	}

	return tokens
}

func TestLexerRunControl(t *testing.T) {
	assertTokens(lexerRun("  \t @ABCD\n   @AB_CD_42"), []token{
		token{controlToken, "ABCD", nil},
		token{controlToken, "AB_CD_42", nil},
	}, t)
}

func TestLexerRunComment(t *testing.T) {
	tokens := lexerRun("  \t \n   # one two\n #--- \n # three ")

	assertTokens(tokens, []token{
		{titleToken, "one", nil},
		{titleToken, "two", nil},
		{titleToken, "three", nil},
	}, t)
}

func TestLexerRunEmpty(t *testing.T) {
	assertTokens(lexerRun(""), []token{}, t)
	assertTokens(lexerRun("#"), []token{}, t)
	assertTokens(lexerRun(" \n #\r\n   #"), []token{}, t)
}

func TestLexerRunIdent(t *testing.T) {
	assertTokens(lexerRun(" \t ABCD\t \n\n   AB_CD_42 \t\r"), []token{
		token{identToken, "ABCD", nil},
		token{identToken, "AB_CD_42", nil},
	}, t)
}

func TestLexerRunName(t *testing.T) {
	assertTokens(lexerRun("\t\t  abcd\t \n  \r ab_cd_42 \t"), []token{
		token{nameToken, "abcd", nil},
		token{nameToken, "ab_cd_42", nil},
	}, t)
}

func TestLexerRunNumber(t *testing.T) {
	assertTokens(lexerRun("\t\t  0.42\t \n -4.2 \r +42.0 \t"), []token{
		token{numberToken, "0.42", nil},
		token{numberToken, "-4.2", nil},
		token{numberToken, "+42.0", nil},
	}, t)
}
