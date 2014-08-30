package tgff

func isAlpha(c byte) bool {
	return c >= 'A' && c <= 'z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isNamely(c byte) bool {
	return isAlpha(c) || isDigit(c) || c == '_'
}
