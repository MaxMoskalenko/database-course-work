package sql_service

import (
	"bufio"
	h "database-course-work/helpers"
	"fmt"
	"log"
	"os"
	"strings"
)

func (db *Database) InitCommodityMarket() {
	db.sql.Exec("USE commodity_market;")

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS commodity_companies (
			id INT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR (255) NOT NULL,
			tag VARCHAR (16) NOT NULL,
			password VARCHAR (255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (tag)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS shipment_companies (
			id INT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR (255) NOT NULL,
			tag VARCHAR (16) NOT NULL,
			password VARCHAR (255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (tag)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS exchangers (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR (255) NOT NULL,
			tag VARCHAR (16) NOT NULL,
			database_name VARCHAR (255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (tag, database_name)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS units (
			id INT PRIMARY KEY AUTO_INCREMENT,
			unit VARCHAR (255) NOT NULL,
			UNIQUE (unit)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS commodity_types (
			id INT PRIMARY KEY AUTO_INCREMENT,
			label VARCHAR (255) NOT NULL,
			unit_id INT NOT NULL,
			UNIQUE (label),
			FOREIGN KEY (unit_id) REFERENCES units(id)
		);`,
	)

	db.fillFromFile(
		"units",
		"./data/units",
		"INSERT IGNORE INTO commodity_market.units (unit) VALUES (?);",
	)
	db.fillFromFile(
		"commodity_types",
		"./data/commodity_types",
		"INSERT IGNORE INTO commodity_market.commodity_types (label, unit_id) VALUES (?, ?);",
	)
}

func (db *Database) fillFromFile(table string, path string, sqlStatement string) {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		strArgs := strings.Split(scanner.Text(), " ")
		intrArgs := make([]interface{}, len(strArgs))
		for i, v := range strArgs {
			intrArgs[i] = v
		}

		_, err := db.sql.Exec(
			sqlStatement,
			intrArgs...,
		)
		if err != nil {
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS users (
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
	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS brokers (
			user_id INT PRIMARY KEY,
			license VARCHAR(25) NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS commodities (
			user_id INT NOT NULL,
			commodity_id INT NOT NULL,
			volume INT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (commodity_id) REFERENCES commodity_market.commodity_types(id),
			PRIMARY KEY (user_id, commodity_id)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id INT PRIMARY KEY AUTO_INCREMENT,
			owner_id INT NOT NULL,
			side ENUM ('buy', 'sell'),
			state ENUM ('active', 'executed'),
			commodity_id INT NOT NULL,
			volume INT NOT NULL,
			pref_broker_id INT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (owner_id) REFERENCES users(id),
			FOREIGN KEY (pref_broker_id) REFERENCES brokers(user_id),
			FOREIGN KEY (commodity_id) REFERENCES commodity_market.commodity_types(id)
		)
	`)

	sqlStatement := "INSERT INTO commodity_market.exchangers (name, tag, database_name) VALUES (?, ?, ?);"

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
