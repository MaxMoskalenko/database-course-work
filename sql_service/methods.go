package sql_service

import (
	h "database-course-work/helpers"
	"fmt"
)

func (db *Database) InitCommodityMarket() {
	db.gorm.Exec("USE commodity_market;")

	db.gorm.Exec(
		`CREATE TABLE IF NOT EXISTS commodity_companies (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR (255) NOT NULL,
		tag VARCHAR (16) NOT NULL,
		password VARCHAR (255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tag)
		);`,
	)
	db.gorm.Exec(
		`CREATE TABLE IF NOT EXISTS shipment_companies (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR (255) NOT NULL,
		tag VARCHAR (16) NOT NULL,
		password VARCHAR (255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tag)
		);`,
	)
	db.gorm.Exec(`CREATE TABLE IF NOT EXISTS exchangers (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR (255) NOT NULL,
		tag VARCHAR (16) NOT NULL,
		database_name VARCHAR (255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tag)
		);`,
	)
}

func (db *Database) InitExchange(databaseName string, exchangerName string, tag string) {
	db.gorm.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", databaseName))
	db.gorm.Exec(fmt.Sprintf("USE %s;", databaseName))

	db.gorm.Exec(
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
	db.gorm.Exec(
		`CREATE TABLE IF NOT EXISTS brokers (
		user_id INT PRIMARY KEY AUTO_INCREMENT,
		license VARCHAR(25) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
	)

	db.gorm.Exec("USE commodity_market;")

	db.gorm.Exec(
		`INSERT INTO exchangers (name, tag, database_name) VALUES (?, ?, ?);`,
		exchangerName,
		tag,
		databaseName,
	)
}

func (db *Database) GetDatabaseByTag(tag string) string {
	db.gorm.Exec("USE commodity_market;")

	var dbName string
	db.gorm.Raw(
		`SELECT database_name FROM exchangers WHERE tag=?;`,
		tag,
	).Scan(&dbName)
	return dbName
}

func (db *Database) CheckIsRecordExist(database string, table string, column string, value string) bool {
	db.gorm.Exec(fmt.Sprintf("USE %s;", database))

	var result string

	db.gorm.Raw(
		fmt.Sprintf(
			`SELECT %s FROM %s WHERE %s='%s';`,
			column,
			table,
			column,
			value,
		),
	).Scan(&result)

	return len(result) != 0
}

func (db *Database) SignUp(database string, user *h.User) {
	db.gorm.Exec(fmt.Sprintf("USE %s;", database))

	db.gorm.Exec(
		`INSERT INTO users (name, surname, email, password, is_broker) VALUES (?, ?, ?, ?, ?);`,
		user.Name,
		user.Surname,
		user.Email,
		user.Password,
		user.IsBroker,
	)
	var id int
	db.gorm.Raw(
		`SELECT id FROM users WHERE email=?;`,
		user.Email,
	).Scan(&id)

	if user.IsBroker == 1 && len(user.License) != 0 {
		db.gorm.Exec(
			`INSERT INTO brokers (user_id, license) VALUES (?, ?);`,
			id,
			user.License,
		)
	}
}

func (db *Database) GetUserOnLogin(database string, login *h.User) *h.User {
	db.gorm.Exec(fmt.Sprintf("USE %s;", database))

	var user h.User

	db.gorm.Raw(
		`SELECT name, email, surname, is_broker
		FROM users 
		WHERE email=? AND password=?;`,
		login.Email,
		login.Password,
	).Scan(&user)

	return &user
}
