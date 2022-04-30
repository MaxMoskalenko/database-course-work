package helpers

func ConvertZeroToNil(number int) interface{} {
	if number == 0 {
		return nil
	}
	return number
}
