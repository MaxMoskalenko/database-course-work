package helpers

func BoolToNumber(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
