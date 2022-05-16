package helpers

import (
	"regexp"
)

func ValidPassword(password string) bool {
	passwordRegExp := "^[a-zA-Z0-9]{8,}$"
	matched, _ := regexp.MatchString(passwordRegExp, password)

	return matched
}

func ValidBankAccount(bankAccount string) bool {
	bankAccountRegExp := "^[0-9]{16}$"
	matched, _ := regexp.MatchString(bankAccountRegExp, bankAccount)

	return matched
}

func ValidEmail(email string) bool {
	emailRegExp := "^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]+$"
	matched, _ := regexp.MatchString(emailRegExp, email)

	return matched
}

func ValidPhone(phone string) bool {
	phoneRegExp := "^\\+[0-9]{1,4} \\([0-9]{3}\\) [0-9]{3}-[0-9]{4}$"
	matched, _ := regexp.MatchString(phoneRegExp, phone)

	return matched
}

func ValidDatabase(exchangerName string) bool {
	databaseRegExp := "^[a-z_]{1,127}$"
	matched, _ := regexp.MatchString(databaseRegExp, exchangerName)

	return matched
}

func ValidTable(str string) bool {
	tableRegExp := "^[a-zA-Z0-9_]{1,127}$"
	matched, _ := regexp.MatchString(tableRegExp, str)

	return matched
}

func ValidExchangerTag(tag string) bool {
	tagRegExp := "^[a-zA-Z0-9]{1,15}$"
	matched, _ := regexp.MatchString(tagRegExp, tag)

	return matched
}
