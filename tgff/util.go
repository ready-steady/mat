package tgff

func isMember(item byte, set []byte) bool {
	for i := 0; i < len(set); i++ {
		if item == set[i] {
			return true
		}
	}
	return false
}

func isAlpha(char byte) bool {
	return char >= 'A' && char <= 'z'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isNamely(char byte) bool {
	return isAlpha(char) || isDigit(char) || char == '_'
}
