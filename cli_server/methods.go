package cli_server

import (
	"fmt"
	"regexp"
	"strings"
)

func readFromConsole(title string, validator string) string {
	var value string
	for {
		fmt.Println(title)
		fmt.Scanln(&value)
		matched, _ := regexp.MatchString(validator, value)
		if len(validator) == 0 || matched {
			break
		}
	}
	return value
}

func GetUserOnSignUp() SignUpResponse {
	userType := readFromConsole("🤨 Please enter type of user\nE - Exchange User\nB - Broker\nC - Company\nS - Shipment Company", "^[eEbBcCsS]$")

	if strings.ToLower(userType) == "e" || strings.ToLower(userType) == "b" {
		name := readFromConsole("🐵 Please enter your name", ".")
		surname := readFromConsole("🙉 Please enter your surname", ".")
		email := readFromConsole("🙊 Please enter your email", "[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]")
		password := readFromConsole("🙈 Please enter your password", ".")
		return SignUpResponse{
			userType,
			UserSignUp{
				name,
				surname,
				email,
				password,
			},
			CompanySignUp{},
		}
	}
	if strings.ToLower(userType) == "c" || strings.ToLower(userType) == "s" {
		title := readFromConsole("🎩 Please enter title of your company", ".")
		email := readFromConsole("🤓 Please enter email of your company", ".")
		password := readFromConsole("🧐 Please enter password of your company", ".")
		return SignUpResponse{
			userType,
			UserSignUp{},
			CompanySignUp{
				title,
				email,
				password,
			},
		}

	}
	return SignUpResponse{}
}
