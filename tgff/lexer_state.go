package tgff

import (
	"errors"
	"fmt"
	"io"
)

type lexState func(*lexer) lexState

const (
	blockCloseChar  = '}'
	blockOpenChar   = '{'
	controlChar     = '@'
	commentChar     = '#'
	digitChars      = "0123456789"
	dotChar         = '.'
	signChars       = "-+"
	whitespaceChars = " \t\n\r"
)

func lexErrorState(err error) lexState {
	return func(l *lexer) lexState {
		l.set(err.Error())
		l.emit(errorToken, err)

		return nil
	}
}

func lexUncertainState(l *lexer) lexState {
	if err := l.skipAny(whitespaceChars); err != nil {
		return lexErrorState(err)
	}

	switch c, err := l.peek(); {
	case err == io.EOF:
		return nil
	case err != nil:
		return lexErrorState(err)
	case c == controlChar:
		return lexControlState
	case c == commentChar:
		return lexCommentState
	case c == blockOpenChar:
		return lexBlockOpenState
	case c == blockCloseChar:
		return lexBlockCloseState
	case isMember(signChars, c) || isMember(digitChars, c):
		return lexNumberState
	case isIdently(c):
		return lexIdentifierState
	case isNamely(c):
		return lexNameState
	default:
		return lexErrorState(errors.New(fmt.Sprintf("unknown token starting from '%c'", c)))
	}
}

func lexBlockOpenState(l *lexer) lexState {
	if err := l.readOne(blockOpenChar); err != nil {
		return lexErrorState(err)
	}

	l.emit(blockOpenToken)

	return lexUncertainState
}

func lexBlockCloseState(l *lexer) lexState {
	if err := l.readOne(blockCloseChar); err != nil {
		return lexErrorState(err)
	}

	l.emit(blockCloseToken)

	return lexUncertainState
}

func lexControlState(l *lexer) lexState {
	if err := l.skipOne(controlChar); err != nil {
		return lexErrorState(err)
	}

	if err := l.readIdent(); err != nil {
		return lexErrorState(err)
	}

	l.emit(controlToken)

	return lexUncertainState
}

func lexCommentState(l *lexer) lexState {
	if err := l.skipLine(); err != nil {
		return lexErrorState(err)
	}

	return lexUncertainState
}

func lexIdentifierState(l *lexer) lexState {
	if err := l.readIdent(); err != nil {
		return lexErrorState(err)
	}

	l.emit(identifierToken)

	return lexUncertainState
}

func lexNameState(l *lexer) lexState {
	if err := l.readName(); err != nil {
		return lexErrorState(err)
	}

	l.emit(nameToken)

	return lexUncertainState
}

func lexNumberState(l *lexer) lexState {
	if err := l.readAny(signChars, digitChars, string(dotChar), digitChars); err != nil {
		return lexErrorState(err)
	}

	l.emit(numberToken)

	return lexUncertainState
}
