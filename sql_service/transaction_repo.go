package sql_service

import (
	h "database-course-work/helpers"
	"database/sql"
)

func updateTransactionOrder(
	tx *sql.Tx,
	order *h.Order,
) {
	sqlStatement := `
		UPDATE orders
		SET volume = ?, state = ?
		WHERE id = ?;
	`

	err := tx.QueryRow(
		sqlStatement,
		order.Commodity.Volume,
		order.State,
		order.Id,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

// TODO commodities things
func upsertTransactionCommodities(
	tx *sql.Tx,
	userId int,
	commodityId int,
	volume float64,
) {
	// sqlStatement := `
	// 	INSERT INTO
	// 		commodities (user_id, commodity_id, volume)
	// 	VALUES (?, ?, ?)
	// 	ON DUPLICATE KEY UPDATE volume = volume + ?;
	// `

	// err := tx.QueryRow(
	// 	sqlStatement,
	// 	userId,
	// 	commodityId,
	// 	volume,
	// 	volume,
	// ).Err()

	// if err != nil {
	// 	tx.Rollback()
	// 	panic(err)
	// }
}
