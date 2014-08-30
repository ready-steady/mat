package tgff

type state func(*lexer) state

const (
	controlMarker = '@'
)

func errorState(err error) state {
	return func(l *lexer) state {
		l.set(err.Error())
		l.emit(errorToken, err)

		return nil
	}
}

func controlState(l *lexer) state {
	if err := l.skipWhitespace(); err != nil {
		return errorState(err)
	}

	if err := l.skipChar(controlMarker); err != nil {
		return errorState(err)
	}

	if err := l.readName(); err != nil {
		return errorState(err)
	}

	l.emit(controlToken)

	return numberState
}

func numberState(l *lexer) state {
	if err := l.skipWhitespace(); err != nil {
		return errorState(err)
	}
	return nil
}
