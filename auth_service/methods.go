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
) {
	database := db.GetDatabaseByTag(exchanger_tag)

	if len(database) == 0 {
		fmt.Printf("⛔️ Database with a %s tag does not exist\n", exchanger_tag)
		return
	}

	if !validPassword(user.Password) {
		fmt.Printf("⛔️ Your password is incorrect\n")
		return
	}

	if !validEmail(user.Email) {
		fmt.Printf("⛔️ Your email is incorrect\n")
		return
	}

	if db.CheckIsRecordExist(database, "users", "email", user.Email) {
		fmt.Printf("⛔️ Your email is already in use\n")
		return
	}

	user.Password = h.Hash(user.Password)

	db.SignUp(database, user)
}

func SignIn(
	db *sql_service.Database,
	exchanger_tag string,
	login *h.User,
) string {
	database := db.GetDatabaseByTag(exchanger_tag)
	login.Password = h.Hash(login.Password)

	if !db.CheckIsRecordExist(database, "users", "email", login.Email) {
		fmt.Printf("⛔️ Wrong credentials\n")
		return ""
	}

	user := db.GetUserOnLogin(database, login)

	if user.Email == "" {
		fmt.Printf("⛔️ Wrong credentials\n")
		return ""
	}

	jwt, err := generateJWT(user)

	if err != nil {
		fmt.Println(err.Error())
	}

	return jwt
}
