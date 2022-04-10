package sql_server

import "fmt"

func (db *Database) InitCommodityMarket() {
	db.gorm.Exec("USE commodity_market;")

	db.gorm.Exec(
		`CREATE TABLE IF NOT EXISTS commodity_companies (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR (255) NOT NULL,
		tag VARCHAR (255) NOT NULL,
		password VARCHAR (255) NOT NULL,
		created_at DATETIME
		);`,
	)
	db.gorm.Exec(
		`CREATE TABLE IF NOT EXISTS shipment_companies (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR (255) NOT NULL,
		tag VARCHAR (255) NOT NULL,
		password VARCHAR (255) NOT NULL,
		created_at DATETIME
		);`,
	)
	db.gorm.Exec(`CREATE TABLE IF NOT EXISTS exchangers (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR (255) NOT NULL,
		tag VARCHAR (255) NOT NULL,
		database_name VARCHAR (255) NOT NULL,
		created_at DATETIME
		);`,
	)
}

func (db *Database) InitExchange(databaseName string) {
	tx := db.gorm.Exec("CREATE DATABASE IF NOT EXISTS " + databaseName + ";")
	fmt.Println(tx)
	db.gorm.Exec("USE " + databaseName + ";")
	db.gorm.Exec(
		`CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR (255) NOT NULL,
		surname VARCHAR (255) NOT NULL,
		email VARCHAR (255) NOT NULL,
		password VARCHAR (255) NOT NULL,
		is_broker BIT,
		created_at DATETIME
		);`,
	)
	db.gorm.Exec(
		`CREATE TABLE IF NOT EXISTS brokers (
		user_id INT PRIMARY KEY AUTO_INCREMENT,
		license VARCHAR(25) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
	)
}
