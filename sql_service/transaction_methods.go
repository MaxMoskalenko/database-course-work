package sql_service

import (
	h "database-course-work/helpers"
	"database/sql"
	"fmt"
)

func updateTransactionOrder(
	tx *sql.Tx,
	database string,
	order *h.Order,
) {
	sqlStatement := fmt.Sprintf(`
		UPDATE %s.orders
		SET volume = ?, state = ?
		WHERE id = ?;
	`, database)

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

func upsertTransactionCommodities(
	tx *sql.Tx,
	database string,
	userId int,
	commodityId int,
	volume float64,
) {
	sqlStatement := fmt.Sprintf(`
		INSERT INTO
			%s.commodities (user_id, commodity_id, volume)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE volume = volume + ?;
	`, database)

	err := tx.QueryRow(
		sqlStatement,
		userId,
		commodityId,
		volume,
		volume,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

func upsertTransactionCargo(
	tx *sql.Tx,
	database string,
	raceId int,
	commodityId int,
	userId int,
	volume float64,
) {
	sqlStatement := fmt.Sprintf(`
		INSERT
			%s.expected_cargo (race_id, user_id, commodity_id, volume)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE volume = volume + ?;
	`, database)

	err := tx.QueryRow(
		sqlStatement,
		raceId,
		userId,
		commodityId,
		volume,
		volume,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}
}
