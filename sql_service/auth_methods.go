package sql_service

import (
	h "database-course-work/helpers"
	"fmt"
)

func (db *Database) SignUp(database string, user *h.User) {
	if !h.ValidDatabase(database) {
		fmt.Println("⛔️ Invalid database name")
		return
	}
	db.sql.Exec(fmt.Sprintf("USE %s;", database))

	sqlStatement := `
		INSERT INTO users (name, surname, email, password, is_broker) 
		VALUES (?, ?, ?, ?, ?); 
	`

	err := db.sql.QueryRow(
		sqlStatement,
		user.Name,
		user.Surname,
		user.Email,
		user.Password,
		user.IsBroker,
	).Err()

	if err != nil {
		panic(err)
	}

	if user.IsBroker == 1 && len(user.License) != 0 {
		id := db.getId(database, "users", "email", user.Email)

		sqlStatement = fmt.Sprintf("INSERT INTO %s.brokers (user_id, license) VALUES (?, ?);", database)
		err = db.sql.QueryRow(
			sqlStatement,
			id,
			user.License,
		).Err()

		if err != nil {
			panic(err)
		}
	}
}

func (db *Database) SignUpCompany(table string, company *h.Company) {
	db.sql.Exec("USE commodity_market;")

	sqlStatement := fmt.Sprintf(`
		INSERT INTO %s (title, tag, password) 
		VALUES (?, ?, ?); 
	`, table)

	err := db.sql.QueryRow(
		sqlStatement,
		company.Title,
		company.Tag,
		company.Password,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) GetUserOnLogin(database string, login *h.User) *h.User {
	if !h.ValidDatabase(database) {
		fmt.Println("🛠 Invalid database name")
		return nil
	}

	var user h.User

	db.sql.Exec(fmt.Sprintf("USE %s;", database))
	sqlStatement := "SELECT email FROM users WHERE email=? AND password=?;"

	db.sql.QueryRow(
		sqlStatement,
		login.Email,
		login.Password,
	).Scan(&user.Email)

	return &user
}

func (db *Database) GetCompanyOnLogin(table string, login *h.Company) *h.Company {
	if !h.ValidTable(table) {
		fmt.Println("🛠 Invalid table name")
		return nil
	}

	db.sql.Exec("USE commodity_market;")

	var company h.Company

	sqlStatement := fmt.Sprintf("SELECT tag FROM %s WHERE tag=? AND password=?;", table)

	db.sql.QueryRow(
		sqlStatement,
		login.Tag,
		login.Password,
	).Scan(&company.Tag)

	return &company
}

func (db *Database) GetUserData(database string, email string) *h.User {
	if !h.ValidDatabase(database) {
		fmt.Println("🛠 Invalid table name")
		return nil
	}

	var user h.User

	sqlStatement := fmt.Sprintf(`
		SELECT email, name, surname, is_broker
		FROM %s.users
		WHERE email=?;
	`, database)

	db.sql.QueryRow(
		sqlStatement,
		email,
	).Scan(&user.Email, &user.Name, &user.Surname, &user.IsBroker)

	user.ExchangerTag = db.GetTagByDatabase(database)

	return &user
}
