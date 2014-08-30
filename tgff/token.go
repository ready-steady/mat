package tgff

type tokenKind uint

type token struct {
	kind  tokenKind
	value string
	more  []interface{}
}

const (
	errorToken tokenKind = iota
	controlToken
	numberToken
)
