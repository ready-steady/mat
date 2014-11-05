package tgff

import (
	"strings"
)

const (
	digits = "0123456789"
	point  = '.'
	signs  = "-+"
)

func isMember(chars string, c byte) bool {
	return strings.IndexByte(chars, c) >= 0
}

func isLowercase(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isUppercase(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isIdently(c byte) bool {
	return isUppercase(c) || isDigit(c) || c == '_'
}

func isNamely(c byte) bool {
	return isLowercase(c) || isDigit(c) || c == '_'
}

func isNumberly(c byte) bool {
	return isDigit(c) || isMember("-+.", c)
}
