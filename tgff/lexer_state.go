package tgff

import (
	"errors"
	"io"
)

type lexState func(*lexer) lexState

const (
	blockCloser = '}'
	blockOpener = '{'
	commentLine = '-'
	commentMark = '#'
	controlMark = '@'
	lineSpace   = " \t"
	newLine     = "\n\r"
	whitespace  = " \t\n\r"
)

func lexErrorState(err error) lexState {
	return func(l *lexer) lexState {
		l.set(err.Error())
		l.send(errorToken)

		return nil
	}
}

func lexEndOrErrorState(err error) lexState {
	if err == io.EOF {
		return nil
	} else {
		return lexErrorState(err)
	}
}

func lexUncertainState(l *lexer) lexState {
	if err := l.skipAny(whitespace); err != nil {
		return lexErrorState(err)
	}

	switch c, err := l.peek(); {
	case err != nil:
		return lexEndOrErrorState(err)
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
		return lexErrorState(errors.New("encountered an unknown token"))
	}
}

func lexBlockOpenState(l *lexer) lexState {
	if err := l.readOne(blockOpener); err != nil {
		return lexErrorState(err)
	}

	if !l.send(blockOpenToken) {
		return nil
	}

	return lexUncertainState
}

func lexBlockCloseState(l *lexer) lexState {
	if err := l.readOne(blockCloser); err != nil {
		return lexErrorState(err)
	}

	if !l.send(blockCloseToken) {
		return nil
	}

	return lexUncertainState
}

func lexControlState(l *lexer) lexState {
	if err := l.skipOne(controlMark); err != nil {
		return lexErrorState(err)
	}

	if err := l.readIdent(); err != nil {
		return lexErrorState(err)
	}

	if !l.send(controlToken) {
		return nil
	}

	return lexUncertainState
}

func lexCommentState(l *lexer) lexState {
	if err := l.skipOne(commentMark); err != nil {
		return lexErrorState(err)
	}

	c, err := l.peek()

	if err != nil {
		return lexEndOrErrorState(err)
	}

	if c != commentLine {
		return lexHeaderState
	}

	if err = l.skipLine(); err != nil {
		return lexErrorState(err)
	}

	return lexUncertainState
}

func lexHeaderState(l *lexer) lexState {
	for {
		if err := l.skipAny(lineSpace); err != nil {
			return lexEndOrErrorState(err)
		}

		if c, err := l.peek(); err != nil {
			return lexEndOrErrorState(err)
		} else if isMember(newLine, c) {
			return lexUncertainState
		}

		if err := l.readName(); err != nil {
			return lexEndOrErrorState(err)
		}

		if !l.send(titleToken) {
			return nil
		}
	}
}

func lexIdentState(l *lexer) lexState {
	if err := l.readIdent(); err != nil {
		return lexErrorState(err)
	}

	if !l.send(identToken) {
		return nil
	}

	return lexUncertainState
}

func lexNameState(l *lexer) lexState {
	if err := l.readName(); err != nil {
		return lexErrorState(err)
	}

	if !l.send(nameToken) {
		return nil
	}

	return lexUncertainState
}

func lexNumberState(l *lexer) lexState {
	if err := l.readAnyOneOf(signs); err != nil {
		return lexErrorState(err)
	}

	if err := l.readOneOf(digits); err != nil {
		return lexErrorState(err)
	}

	if err := l.readAny(digits, string(point), digits); err != nil {
		return lexErrorState(err)
	}

	if !l.send(numberToken) {
		return nil
	}

	return lexUncertainState
}
