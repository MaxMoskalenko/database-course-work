package sql_service

import (
	h "database-course-work/helpers"
	"fmt"
)

func (db *Database) AddCommodity(database string, userEmail string, commodity *h.Commodity) {
	if !h.ValidDatabase(database) {
		fmt.Println("🛠 Invalid database name")
		return
	}

	userId := db.getId(database, "users", "email", userEmail)
	commodityId := db.getId("commodity_market", "commodity_types", "label", commodity.Label)

	sqlStatement := fmt.Sprintf(`
		INSERT INTO %s.commodities (user_id, commodity_id, volume) 
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE volume = volume + ?;
	`, database)

	err := db.sql.QueryRow(
		sqlStatement,
		userId,
		commodityId,
		commodity.Volume,
		commodity.Volume,
	).Err()

	if err != nil {
		fmt.Println(err)
	}
}

func (db *Database) GetUserCommodities(database string, userEmail string) [](*h.Commodity) {
	if !h.ValidDatabase(database) {
		fmt.Println("🛠 Invalid database name")
		return nil
	}

	userId := db.getId(database, "users", "email", userEmail)

	sqlStatement := fmt.Sprintf(`
		SELECT CT.label, C.volume, CU.unit
		FROM %s.commodities AS C
		INNER JOIN (
			SELECT label, unit_id, id
			FROM commodity_market.commodity_types
		) AS CT
			ON CT.id=C.commodity_id 
		INNER JOIN (
			SELECT unit, id
			FROM commodity_market.units
		) AS CU
			ON CU.id = CT.unit_id 
		WHERE user_id=?;
	`, database)

	rows, err := db.sql.Query(
		sqlStatement,
		userId,
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var commodities []*h.Commodity

	for rows.Next() {
		var commodity h.Commodity
		if err := rows.Scan(&commodity.Label, &commodity.Volume, &commodity.Unit); err != nil {
			panic(err)
		}
		commodity.Owner = *db.GetUserData(database, userEmail)
		commodities = append(commodities, &commodity)
	}

	return commodities
}

func (db *Database) GetAllCommodities(database string) [](*h.Commodity) {
	if !h.ValidDatabase(database) {
		fmt.Println("🛠 Invalid database name")
		return nil
	}

	sqlStatement := fmt.Sprintf(`
		SELECT CT.label, C.volume, CU.unit, U.email
		FROM %s.commodities AS C
		INNER JOIN (
			SELECT label, unit_id, id
			FROM commodity_market.commodity_types
		) AS CT
			ON CT.id=C.commodity_id 
		INNER JOIN (
			SELECT unit, id
			FROM commodity_market.units
		) AS CU
			ON CU.id = CT.unit_id
		INNER JOIN (
			SELECT email , id
			FROM %s.users
		) AS U
			ON U.id = C.user_id 
	`, database, database)

	rows, err := db.sql.Query(sqlStatement)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var commodities []*h.Commodity

	for rows.Next() {
		var commodity h.Commodity
		var email string
		if err := rows.Scan(&commodity.Label, &commodity.Volume, &commodity.Unit, &email); err != nil {
			panic(err)
		}

		commodity.Owner = *db.GetUserData(database, email)
		commodities = append(commodities, &commodity)
	}

	return commodities
}
