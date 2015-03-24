package tgff

import (
	"strings"
	"testing"

	"github.com/ready-steady/assert"
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
	lexer := fakeLexer("abbbaacdefg")

	err := lexer.readAny("ab")

	assert.Success(err, t)
	assertLexerAt(lexer, 'c', t)
	assert.Equal(lexer.value(), "abbbaa", t)
}

func TestLexerReadName(t *testing.T) {
	scenarios := []struct {
		data string
		name string
	}{
		{"abcd efgh", "abcd"},
	}

	for _, s := range scenarios {
		lexer := fakeLexer(s.data)
		_ = lexer.readName()
		assert.Equal(lexer.value(), s.name, t)
	}
}

func TestLexerSkipAny(t *testing.T) {
	lexer := fakeLexer("abbbaacdefg")

	err := lexer.skipAny("ab")

	assert.Success(err, t)
	assertLexerAt(lexer, 'c', t)
}

func TestLexerSkipSequence(t *testing.T) {
	lexer := fakeLexer("abcde")

	err := lexer.skipSequence("ab")

	assert.Success(err, t)
	assertLexerAt(lexer, 'c', t)

	err = lexer.skipSequence("d")

	assert.Failure(err, t)
}

func TestLexerRunControl(t *testing.T) {
	assertTokens(fakeLexerRun("  \t @ABCD\n   @AB_CD_42"), []token{
		token{controlToken, "ABCD"},
		token{controlToken, "AB_CD_42"},
	}, t)
}

func TestLexerRunComment(t *testing.T) {
	tokens := fakeLexerRun("  \t \n   # one two\n #--- \n # three ")

	assertTokens(tokens, []token{
		{titleToken, "one"},
		{titleToken, "two"},
		{titleToken, "three"},
	}, t)
}

func TestLexerRunEmpty(t *testing.T) {
	assertTokens(fakeLexerRun(""), []token{}, t)
	assertTokens(fakeLexerRun("#"), []token{}, t)
	assertTokens(fakeLexerRun(" \n #\r\n   #"), []token{}, t)
}

func TestLexerRunIdent(t *testing.T) {
	assertTokens(fakeLexerRun(" \t ABCD\t \n\n   AB_CD_42 \t\r"), []token{
		token{identToken, "ABCD"},
		token{identToken, "AB_CD_42"},
	}, t)
}

func TestLexerRunName(t *testing.T) {
	assertTokens(fakeLexerRun("\t\t  abcd\t \n  \r ab_cd_42 \t"), []token{
		token{nameToken, "abcd"},
		token{nameToken, "ab_cd_42"},
	}, t)
}

func TestLexerRunNumber(t *testing.T) {
	assertTokens(fakeLexerRun("\t\t  0.42\t \n -4.2 \r +42.0 \t"), []token{
		token{numberToken, "0.42"},
		token{numberToken, "-4.2"},
		token{numberToken, "+42.0"},
	}, t)
}

func fakeLexerRun(data string) []token {
	lexer, stream := newLexer(strings.NewReader(data), make(chan bool))

	go lexer.run()

	tokens := []token{}
	for token := range stream {
		tokens = append(tokens, *token)
	}

	return tokens
}

func fakeLexer(data string) *lexer {
	lexer, _ := newLexer(strings.NewReader(data), make(chan bool))
	return lexer
}
