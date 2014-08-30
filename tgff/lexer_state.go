package tgff

import (
	"errors"
	"io"
)

type lexState func(*lexer) lexState

const (
	controlMarker   = '@'
	digitChars      = "0123456789"
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
	if err := l.skipChars(whitespaceChars); err != nil {
		return lexErrorState(err)
	}

	switch c, err := l.peek(); {
	case err == io.EOF:
		return nil
	case err != nil:
		return lexErrorState(err)
	case c == '@':
		return lexControlState
	case isMember("-+0123456789", c):
		return lexNumberState
	default:
		return lexErrorState(errors.New("unknown token"))
	}
}

func lexControlState(l *lexer) lexState {
	if err := l.skipChar(controlMarker); err != nil {
		return lexErrorState(err)
	}

	if err := l.readName(); err != nil {
		return lexErrorState(err)
	}

	l.emit(controlToken)

	return lexUncertainState
}

func lexNumberState(l *lexer) lexState {
	if err := l.readChars("+-", digitChars, ".", digitChars); err != nil {
		return lexErrorState(err)
	}

	if l.length() == 0 {
		return lexErrorState(errors.New("expected a number"))
	}

	l.emit(numberToken)

	return lexUncertainState
}
