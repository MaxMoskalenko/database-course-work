package auth_service

import (
	h "database-course-work/helpers"
	sql_service "database-course-work/sql_service"
	"fmt"
)

func SignUp(
	db *sql_service.Database,
	exchanger_tag string,
	user *h.User,
) string {
	database := db.GetDatabaseByTag(exchanger_tag)

	if len(database) == 0 {
		fmt.Printf("⛔️ Database with a %s tag does not exist\n", exchanger_tag)
		return ""
	}

	if !h.ValidPassword(user.Password) {
		fmt.Printf("⛔️ Your password is incorrect\n")
		return ""
	}

	if !h.ValidEmail(user.Email) {
		fmt.Printf("⛔️ Your email is incorrect\n")
		return ""
	}

	if db.CheckIsRecordExist(database, "users", "email", user.Email) {
		fmt.Printf("⛔️ Your email is already in use\n")
		return ""
	}

	user.Password = h.Hash(user.Password)

	db.SignUp(database, user)

	jwt, err := generateJWT(user)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return jwt
}

func SignUpCompany(
	db *sql_service.Database,
	company *h.Company,
) string {
	table := h.GetTableFromType(company.Type)

	if table == "" {
		fmt.Println("🛠 Comapny Type is incorrect")
		return ""
	}

	if !h.ValidPassword(company.Password) {
		fmt.Printf("⛔️ Your password is incorrect\n")
		return ""
	}

	if db.CheckIsRecordExist("commodity_market", table, "tag", company.Tag) {
		fmt.Printf("⛔️ This tag is already in use\n")
		return ""
	}

	company.Password = h.Hash(company.Password)

	db.SignUpCompany(table, company)

	jwt, err := generateCompanyJWT(company)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return jwt
}

func SignIn(
	db *sql_service.Database,
	exchanger_tag string,
	login *h.User,
) string {
	database := db.GetDatabaseByTag(exchanger_tag)
	login.Password = h.Hash(login.Password)

	user := db.GetUserOnLogin(database, login)

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
	table := h.GetTableFromType(login.Type)

	if table == "" {
		fmt.Println("🛠 Comapny Type is incorrect")
		return ""
	}

	login.Password = h.Hash(login.Password)

	company := db.GetCompanyOnLogin(table, login)

	if company.Tag == "" {
		fmt.Printf("⛔️ Wrong credentials\n")
		return ""
	}

	company.Type = login.Type

	jwt, err := generateCompanyJWT(company)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return jwt
}
