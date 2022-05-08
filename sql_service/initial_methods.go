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
			tag VARCHAR (255) NOT NULL,
			password VARCHAR (255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (tag)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS shipment_companies (
			id INT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR (255) NOT NULL,
			tag VARCHAR (255) NOT NULL,
			password VARCHAR (255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (tag)
		);`,
	)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS exchangers (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR (255) NOT NULL,
			tag VARCHAR (255) NOT NULL,
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

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS commodity_market.races (
			id INT PRIMARY KEY AUTO_INCREMENT,
			from_id INT NOT NULL,
			to_id INT NOT NULL,
			race_date TIMESTAMP NOT NULL,
			company_id INT NOT NULL,
			status ENUM ('preparing', 'arrive', 'permanent') DEFAULT 'preparing',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (from_id) REFERENCES exchangers(id),
			FOREIGN KEY (to_id) REFERENCES exchangers(id)
		);
	`)

	db.sql.Exec(`
		CREATE TRIGGER checkCorrectRaceRoute 
		BEFORE INSERT ON commodity_market.races
		FOR EACH ROW
		BEGIN
			IF NEW.from_id = NEW.to_id
			THEN
				SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = "Start and finish points of race are the same";
			END IF;
		END;
	`)

	db.sql.Exec(`
		CREATE PROCEDURE GetForeignOrders()
		BEGIN
		DECLARE isDone INT;
		DECLARE databaseName VARCHAR(255); 
		DECLARE exchTag VARCHAR(255);
		
		DECLARE ExchangerCursor CURSOR FOR SELECT database_name, tag FROM commodity_market.exchangers;
		DECLARE CONTINUE HANDLER FOR NOT FOUND SET isDone = 1;
		OPEN ExchangerCursor;
		
		SET @globalQuery := 'SELECT * FROM (';
		SET isDone = 0;
		REPEAT
			FETCH ExchangerCursor INTO databaseName, exchTag;
		
			SET @globalQuery := CONCAT(
				@globalQuery,
				' SELECT  O.id, O.side, O.state, CT.label, CU.unit, O.volume, U.name as user_name, U.surname as user_surname, U.email, \'',
				exchTag,
				'\' AS tag FROM ',
				databaseName,
				'.orders AS O ',
				'JOIN (SELECT id, label, unit_id FROM commodity_market.commodity_types) AS CT ON CT.id = O.commodity_id ',
				'JOIN (SELECT id, unit FROM commodity_market.units) AS CU ON CU.id = CT.unit_id ',
				'JOIN (SELECT id, name, surname, email FROM kyiv_central_ex.users) AS U ON U.id = O.owner_id ',
				'WHERE O.pref_broker_id IS NULL UNION'
			);
		UNTIL isDone END REPEAT;
		
		CLOSE ExchangerCursor;
		
		SET @globalQuery := CONCAT(@globalQuery, ' SELECT NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL ) as G WHERE state = \'active\';');
		PREPARE stat FROM @globalQuery;
		EXECUTE stat;
		DEALLOCATE PREPARE stat;
		END
	`)
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
			volume FLOAT NOT NULL,
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
			volume FLOAT NOT NULL,
			executed_volume FLOAT NOT NULL DEFAULT 0,
			pref_broker_id INT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (owner_id) REFERENCES users(id),
			FOREIGN KEY (pref_broker_id) REFERENCES brokers(user_id),
			FOREIGN KEY (commodity_id) REFERENCES commodity_market.commodity_types(id)
		)
	`)

	db.sql.Exec(`
		CREATE TABLE IF NOT EXISTS expected_cargo (
			race_id INT NOT NULL,
			user_id INT NOT NULL,
			commodity_id INT NOT NULL,
			volume FLOAT NOT NULL,
			FOREIGN KEY (race_id) REFERENCES commodity_market.races(id),
			FOREIGN KEY (commodity_id) REFERENCES commodity_market.commodity_types(id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			PRIMARY KEY (race_id, user_id, commodity_id)
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
