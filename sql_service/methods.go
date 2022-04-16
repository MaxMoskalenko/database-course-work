package sql_service

import (
	h "database-course-work/helpers"
	"fmt"
)

func (db *Database) InitCommodityMarket() {
	db.sql.Exec("USE commodity_market;")

	db.sql.Exec(
		`CREATE TABLE IF NOT EXISTS commodity_companies (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR (255) NOT NULL,
		tag VARCHAR (16) NOT NULL,
		password VARCHAR (255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tag)
		);`,
	)
	db.sql.Exec(
		`CREATE TABLE IF NOT EXISTS shipment_companies (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR (255) NOT NULL,
		tag VARCHAR (16) NOT NULL,
		password VARCHAR (255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tag)
		);`,
	)
	db.sql.Exec(`CREATE TABLE IF NOT EXISTS exchangers (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR (255) NOT NULL,
		tag VARCHAR (16) NOT NULL,
		database_name VARCHAR (255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tag, database_name)
		);`,
	)
}

func (db *Database) UseExchangerDB(exName string) bool {
	if !h.ValidDatabase(exName) {
		fmt.Println("⛔️ Invalid database name")
		return false
	}

	err := db.sql.QueryRow(fmt.Sprintf("USE %s;", exName)).Err()
	if err != nil {
		panic(err)
	}
	return true
}

func (db *Database) InitExchange(ex *h.Exchanger) {
	if !h.ValidDatabase(ex.DatabaseName) || !h.ValidExchangerTag(ex.Tag) {
		fmt.Println("⛔️ Invalid exchange credetials")
		return
	}

	if db.CheckIsRecordExist("commodity_market", "exchangers", "database_name", ex.DatabaseName) ||
		db.CheckIsRecordExist("commodity_market", "exchangers", "tag", ex.Tag) {
		fmt.Println("⛔️ Exchanger has already exists")
		return
	}

	db.sql.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", ex.DatabaseName))
	db.sql.Exec(fmt.Sprintf("USE %s;", ex.DatabaseName))

	db.sql.Exec(
		`CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR (255) NOT NULL,
		surname VARCHAR (255) NOT NULL,
		email VARCHAR (255) NOT NULL,
		password VARCHAR (255) NOT NULL,
		is_broker INT8,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (email)
		);`,
	)
	db.sql.Exec(
		`CREATE TABLE IF NOT EXISTS brokers (
		user_id INT PRIMARY KEY AUTO_INCREMENT,
		license VARCHAR(25) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
	)

	db.sql.Exec("USE commodity_market;")
	sqlStatement := "INSERT INTO exchangers (name, tag, database_name) VALUES (?, ?, ?);"

	_, err := db.sql.Exec(
		sqlStatement,
		ex.Name,
		ex.Tag,
		ex.DatabaseName,
	)
	if err != nil {
		panic(err)
	}
}

func (db *Database) GetDatabaseByTag(tag string) string {
	var dbName string
	db.sql.Exec("USE commodity_market;")

	sqlStatement := "SELECT database_name FROM exchangers WHERE tag=?;"

	err := db.sql.QueryRow(
		sqlStatement,
		tag,
	).Scan(&dbName)

	if err != nil {
		panic(err)
	}

	return dbName
}

func (db *Database) CheckIsRecordExist(database string, table string, column string, value string) bool {
	if !h.ValidDatabase(database) {
		fmt.Println("🛠 Invalid database name")
		return false
	}

	if !h.ValidTable(table) || !h.ValidTable(column) {
		fmt.Println("🛠 Invalid table/column name")
		return false
	}

	db.sql.Exec(fmt.Sprintf("USE %s;", database))

	var result string
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?;", column, table, column)

	db.sql.QueryRow(
		sqlStatement,
		value,
	).Scan(&result)

	return len(result) != 0
}

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
		db.sql.Exec(fmt.Sprintf("USE %s;", database))

		id := 0
		sqlStatement = "SELECT id FROM users WHERE email=?;"
		err = db.sql.QueryRow(
			sqlStatement,
			user.Email,
		).Scan(&id)

		if err != nil {
			panic(err)
		}

		sqlStatement = "INSERT INTO brokers (user_id, license) VALUES (?, ?);"
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
