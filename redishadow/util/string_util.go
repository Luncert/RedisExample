package util

func IsAlpha(b int) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z'
}

func IsDigit(b int) bool {
	return b >= '0' && b <= '9'
}
