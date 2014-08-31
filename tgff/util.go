package tgff

import (
	"strings"
)

func isMember(chars string, c byte) bool {
	return strings.IndexByte(chars, c) >= 0
}

func isLowerLetter(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isUpperLetter(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isIdently(c byte) bool {
	return isUpperLetter(c) || c == '_'
}

func isNamely(c byte) bool {
	return isLowerLetter(c) || isDigit(c) || c == '_'
}
