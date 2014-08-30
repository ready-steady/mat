package tgff

type token interface {
}

type errorToken struct {
	err error
}

type controlToken struct {
	name string
}
