package auth_service

import (
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
)

func SignUp(
	db *sql_service.Database,
	user *h.User,
) string {
	if !h.ValidPassword(user.Password) {
		fmt.Printf("⛔️ Your password is incorrect\n")
		return ""
	}

	if !h.ValidEmail(user.Email) {
		fmt.Printf("⛔️ Your email is incorrect\n")
		return ""
	}

	if !h.ValidBankAccount(user.BankAccount) {
		fmt.Printf("⛔️ Your email is incorrect\n")
		return ""
	}

	if db.CheckIsRecordExist("users", "email", user.Email) {
		fmt.Printf("⛔️ Your email is already in use\n")
		return ""
	}

	user.Password = h.Hash(user.Password)

	db.SignUp(user)

	if len(user.License.Code) != 0 {
		AssignBroker(db, user)
	}

	jwt, err := generateJWT(user)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return jwt
}

func AssignBroker(
	db *sql_service.Database,
	user *h.User,
) {

	if !db.CheckIsRecordExist("licenses", "license_code", user.License.Code) {
		panic(fmt.Errorf("⛔️ Invalid license code"))
	}

	user.Id = db.GetId("users", "email", user.Email)
	license := db.GetLicense(user.License.Code)

	if license.IsTaken {
		panic(fmt.Errorf("⛔️ This license code is taken"))
	}

	db.AssignBroker(user.Id, license.Id)
}

func SignUpCompany(
	db *sql_service.Database,
	company *h.Company,
) string {
	if !h.ValidPassword(company.Password) {
		fmt.Printf("⛔️ Your password is incorrect\n")
		return ""
	}

	if db.CheckIsRecordExist("companies", "tag", company.Tag) {
		fmt.Printf("⛔️ This tag is already in use\n")
		return ""
	}

	company.Password = h.Hash(company.Password)

	db.SignUpCompany(company)

	jwt, err := generateCompanyJWT(company)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return jwt
}

func SignIn(
	db *sql_service.Database,
	login *h.User,
) string {
	login.Password = h.Hash(login.Password)

	user := db.GetUserOnLogin(login)

	if user.Email == "" {
		fmt.Printf("⛔️ Wrong credentials\n")
		return ""
	}

	jwt, err := generateJWT(user)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return jwt
}

func SignInCompany(
	db *sql_service.Database,
	login *h.Company,
) string {
	login.Password = h.Hash(login.Password)

	company := db.GetCompanyOnLogin(login)

	if company.Tag == "" {
		fmt.Printf("⛔️ Wrong credentials\n")
		return ""
	}

	jwt, err := generateCompanyJWT(company)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return jwt
}

func GetUser(
	db *sql_service.Database,
	jwt string,
) *h.User {
	user, err := ReadJWT(jwt)
	if err != nil {
		panic(err)
	}

	userData := db.GetUserData(user.Email)

	if userData == nil {
		panic(fmt.Errorf("⛔️ User does not exist"))
	}

	user.Id = userData.Id
	user.Name = userData.Name
	user.Surname = userData.Surname
	user.IsBroker = userData.IsBroker

	return user
}
