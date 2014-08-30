package tgff

func isMember(item byte, set []byte) bool {
	for i := 0; i < len(set); i++ {
		if item == set[i] {
			return true
		}
	}
	return false
}
