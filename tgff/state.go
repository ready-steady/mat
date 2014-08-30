package tgff

type state func(*lexer) state

const (
	controlMarker = '@'
)

func errorState(err error) state {
	return func(l *lexer) state {
		l.emit(errorToken{err})
		return nil
	}
}

func stripState(l *lexer) state {
	if err := l.acceptWhitespace(); err != nil {
		return errorState(err)
	}

	return controlState
}

func controlState(l *lexer) state {
	if err := l.requireChar(controlMarker); err != nil {
		return errorState(err)
	}

	if name, err := l.readName(); err != nil {
		return errorState(err)
	} else {
		l.emit(controlToken{name})

		switch name {
		case "HYPERPERIOD":
		default:
		}
	}

	return nil
}
