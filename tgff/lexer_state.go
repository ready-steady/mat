package tgff

import (
	"errors"
	"io"
)

type lexState func(*lexer) lexState

const (
	blockCloseChar  = '}'
	blockOpenChar   = '{'
	controlChar     = '@'
	digitChars      = "0123456789"
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
	case isMember(signChars, c) || isMember(digitChars, c):
		return lexNumberState
	case c == blockOpenChar:
		return lexBlockOpenState
	case c == blockCloseChar:
		return lexBlockCloseState
	default:
		return lexErrorState(errors.New("unknown token"))
	}
}

func lexControlState(l *lexer) lexState {
	if err := l.skipSequence(string(controlChar)); err != nil {
		return lexErrorState(err)
	}

	if err := l.readName(); err != nil {
		return lexErrorState(err)
	}

	l.emit(controlToken)

	return lexUncertainState
}

func lexNumberState(l *lexer) lexState {
	if err := l.readAny(signChars, digitChars, ".", digitChars); err != nil {
		return lexErrorState(err)
	}

	l.emit(numberToken)

	return lexUncertainState
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
