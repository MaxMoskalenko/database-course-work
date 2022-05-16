package sql_service

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func (db *Database) InitDatabase() {
	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS companies (
			id INT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR (255) NOT NULL,
			tag VARCHAR (255) NOT NULL,
			password VARCHAR (255) NOT NULL,
			email VARCHAR (255) NOT NULL,
			phone_number VARCHAR (255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (tag)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS units (
			id INT PRIMARY KEY AUTO_INCREMENT,
			unit VARCHAR (255) NOT NULL,
			UNIQUE (unit)
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS commodity_types (
			id INT PRIMARY KEY AUTO_INCREMENT,
			label VARCHAR (255) NOT NULL,
			unit_id INT NOT NULL,
			UNIQUE (label)
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS licenses (
			id INT PRIMARY KEY AUTO_INCREMENT,
			license_code VARCHAR(255) NOT NULL,
			is_taken BOOLEAN DEFAULT FALSE
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR (255) NOT NULL,
			surname VARCHAR (255) NOT NULL,
			email VARCHAR (255) NOT NULL,
			bank_account VARCHAR (255) NOT NULL,
			password VARCHAR (255) NOT NULL,
			is_broker BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (email)
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS brokers (
			user_id INT PRIMARY KEY NOT NULL,
			license_id INT NOT NULL
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS commodities_account (
			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
			owner_id INT NOT NULL,
			commodity_id INT NOT NULL,
			volume FLOAT NOT NULL,
			source ENUM('company', 'trade') NOT NULL,
			transaction_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS source_commodities_company (
			transaction_id INT NOT NULL PRIMARY KEY,
			source_company_id INT NOT NULL
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS source_commodities_trade (
			transaction_id INT NOT NULL PRIMARY KEY,
			source_user_id INT NOT NULL,
			source_order_id INT NOT NULL,
			dest_order_id INT NOT NULL,
			broker_id INT NOT NULL
		);
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id INT PRIMARY KEY AUTO_INCREMENT,
			owner_id INT NOT NULL,
			side ENUM ('buy', 'sell'),
			state ENUM ('active', 'executed', 'canceled') DEFAULT 'active',
			commodity_id INT NOT NULL,
			volume FLOAT NOT NULL,
			executed_volume FLOAT NOT NULL DEFAULT 0,
			pref_broker_id INT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`)

	db.sql.Exec(`
		ALTER TABLE orders
		ADD FOREIGN KEY (owner_id) REFERENCES users(id),
		ADD FOREIGN KEY (pref_broker_id) REFERENCES brokers(user_id),
		ADD FOREIGN KEY (commodity_id) REFERENCES commodity_types(id);
	`)

	db.sql.Exec(`
		ALTER TABLE source_commodities_trade
		ADD FOREIGN KEY (transaction_id) REFERENCES commodities_account(id),
		ADD FOREIGN KEY (source_user_id) REFERENCES users(id),
		ADD FOREIGN KEY (source_order_id) REFERENCES orders(id),
		ADD FOREIGN KEY (dest_order_id) REFERENCES orders(id),
		ADD FOREIGN KEY (broker_id) REFERENCES users(id)
	`)

	db.sql.Exec(`
		ALTER TABLE source_commodities_company
		ADD FOREIGN KEY (transaction_id) REFERENCES commodities_account(id),
		ADD FOREIGN KEY (source_company_id) REFERENCES companies(id);
	`)

	db.sql.Exec(`
		ALTER TABLE commodities_account
		ADD FOREIGN KEY (owner_id) REFERENCES users(id),
		ADD FOREIGN KEY (commodity_id) REFERENCES commodity_types(id);
	`)

	db.sql.Exec(`
		ALTER TABLE brokers
		ADD FOREIGN KEY (license_id) REFERENCES licenses(id),
		ADD FOREIGN KEY (user_id) REFERENCES users(id);
	`)

	db.sql.Exec(`
		ALTER TABLE commodity_types
		ADD FOREIGN KEY (unit_id) REFERENCES units(id);
	`)

	db.sql.Exec(`
		CREATE TRIGGER turnNameAndSurnameToUppercase 
		BEFORE INSERT ON users
		FOR EACH ROW
		BEGIN
			SET @first_name_letter = UPPER(SUBSTRING(NEW.name, 1, 1));
			SET @other_name_letters = SUBSTR(NEW.name, 2);

			SET @first_surname_letter = UPPER(SUBSTRING(NEW.surname, 1, 1));
			SET @other_surname_letters = SUBSTR(NEW.surname, 2);

			SET NEW.name = CONCAT(@first_name_letter, @other_name_letters);
			SET NEW.surname = CONCAT(@first_surname_letter, @other_surname_letters);
		END;
	`)

	db.sql.Exec(`
		CREATE PROCEDURE GetUnlockedVolume(IN i_user_id INT, IN i_commodity_id INT)
		BEGIN
			SET @total_volume := (
				SELECT COALESCE(SUM(volume), 0) FROM commodities_account
				WHERE commodity_id = i_commodity_id  AND owner_id = i_user_id
			);
		
			SET @locked_volume := (
				SELECT COALESCE((SUM(volume)-SUM(executed_volume)), 0)
				FROM orders
				WHERE commodity_id = i_commodity_id AND owner_id = i_user_id AND side = 'sell' AND state = 'active'
			);
		
			SELECT @total_volume - @locked_volume AS unlocked_volume;
		END;
	`)

	db.fillFromFile(
		"units",
		"./data/units",
		"INSERT IGNORE INTO units (unit) VALUES (?);",
	)
	db.fillFromFile(
		"commodity_types",
		"./data/commodity_types",
		"INSERT IGNORE INTO commodity_types (label, unit_id) VALUES (?, ?);",
	)
	db.fillFromFile(
		"licenses",
		"./data/licenses",
		"INSERT IGNORE INTO licenses (license_code) VALUES (?);",
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
