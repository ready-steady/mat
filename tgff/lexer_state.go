package tgff

import (
	"errors"
	"fmt"
	"io"
)

type lexState func(*lexer) lexState

const (
	blockCloser = '}'
	blockOpener = '{'
	controlMark = '@'
	commentMark = '#'
	whitespaces = " \t\n\r"
)

func lexErrorState(err error) lexState {
	return func(l *lexer) lexState {
		l.set(err.Error())
		l.emit(errorToken, err)

		return nil
	}
}

func lexUncertainState(l *lexer) lexState {
	if err := l.skipAny(whitespaces); err != nil {
		return lexErrorState(err)
	}

	switch c, err := l.peek(); {
	case err == io.EOF:
		return nil
	case err != nil:
		return lexErrorState(err)
	case c == controlMark:
		return lexControlState
	case c == commentMark:
		return lexCommentState
	case c == blockOpener:
		return lexBlockOpenState
	case c == blockCloser:
		return lexBlockCloseState
	case isNumberly(c):
		return lexNumberState
	case isIdently(c):
		return lexIdentState
	case isNamely(c):
		return lexNameState
	default:
		return lexErrorState(errors.New(fmt.Sprintf("unknown token starting from '%c'", c)))
	}
}

func lexBlockOpenState(l *lexer) lexState {
	if err := l.readOne(blockOpener); err != nil {
		return lexErrorState(err)
	}

	l.emit(blockOpenToken)

	return lexUncertainState
}

func lexBlockCloseState(l *lexer) lexState {
	if err := l.readOne(blockCloser); err != nil {
		return lexErrorState(err)
	}

	l.emit(blockCloseToken)

	return lexUncertainState
}

func lexControlState(l *lexer) lexState {
	if err := l.skipOne(controlMark); err != nil {
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

func lexIdentState(l *lexer) lexState {
	if err := l.readIdent(); err != nil {
		return lexErrorState(err)
	}

	l.emit(identToken)

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
	if err := l.readAny(signs, digits, string(point), digits); err != nil {
		return lexErrorState(err)
	}

	l.emit(numberToken)

	return lexUncertainState
}
