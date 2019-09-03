package util

func IsAlpha(b rune) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z'
}

func IsDigit(b rune) bool {
	return b >= '0' && b <= '9'
}
