package tgff

import (
	"fmt"
	"strconv"
)

type tokenKind uint

type token struct {
	kind  tokenKind
	value string
}

const (
	errorToken tokenKind = iota
	blockCloseToken
	blockOpenToken
	controlToken
	identToken
	nameToken
	numberToken
	titleToken
)

func (t token) String() string {
	return fmt.Sprintf("%v (%v)", t.kind, t.value)
}

func (k tokenKind) String() string {
	switch k {
	case errorToken:
		return "Error"
	case blockCloseToken:
		return "Block Close"
	case blockOpenToken:
		return "Block Open"
	case controlToken:
		return "Control"
	case identToken:
		return "Ident"
	case nameToken:
		return "Name"
	case numberToken:
		return "Number"
	case titleToken:
		return "Title"
	default:
		return "Unknown"
	}
}

func (t token) Uint32() uint32 {
	value, _ := strconv.ParseUint(t.value, 10, 32)
	return uint32(value)
}

func (t token) Float64() float64 {
	value, _ := strconv.ParseFloat(t.value, 64)
	return value
}
