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
