package tgff

import (
	"io"
)

type Result struct {
	graphCount uint
	tableCount uint
}

func Parse(reader io.Reader) *Result {
	return &Result{}
}
