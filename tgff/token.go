package tgff

type tokenKind uint

type token struct {
	kind tokenKind
	value string
}

const (
	errorToken tokenKind = iota
	controlToken
)
