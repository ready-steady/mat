package tgff

import (
	"errors"
)

type lexState func(*lexer) lexState

const (
	controlMarker = '@'
)

func lexErrorState(err error) lexState {
	return func(l *lexer) lexState {
		l.set(err.Error())
		l.emit(errorToken, err)

		return nil
	}
}

func lexControlState(l *lexer) lexState {
	if err := l.skipWhitespace(); err != nil {
		return lexErrorState(err)
	}

	if err := l.skipChar(controlMarker); err != nil {
		return lexErrorState(err)
	}

	if err := l.readName(); err != nil {
		return lexErrorState(err)
	}

	l.emit(controlToken)

	return lexNumberState
}

func lexNumberState(l *lexer) lexState {
	if err := l.skipWhitespace(); err != nil {
		return lexErrorState(err)
	}

	if err := l.readChars("+-.0123456789eE"); err != nil {
		return lexErrorState(err)
	}

	if l.length() == 0 {
		return lexErrorState(errors.New("expected a number"))
	}

	l.emit(numberToken)

	return nil
}
