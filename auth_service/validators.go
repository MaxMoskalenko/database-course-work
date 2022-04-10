package auth_service

import "regexp"

func validPassword(password string) bool {
	passwordRegExp := "^[a-zA-Z0-9]{8,}$"
	matched, _ := regexp.MatchString(passwordRegExp, password)

	return matched
}

func validEmail(email string) bool {
	emailRegExp := "[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]"
	matched, _ := regexp.MatchString(emailRegExp, email)

	return matched
}
