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
		SET volume = ?, state = ?, executed_volume = ?
		WHERE id = ?;
	`

	err := tx.QueryRow(
		sqlStatement,
		order.Commodity.Volume,
		order.State,
		order.ExecutedVolume,
		order.Id,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

func upsertTransactionCommodities(
	tx *sql.Tx,
	userId int,
	commodityId int,
	volume float64,
	source *h.CommoditySource,
) {
	sqlStatement := `
		INSERT INTO
			commodities_account (owner_id, commodity_id, volume, source)
		VALUES (?, ?, ?, 'trade');
	`

	res, err := tx.Exec(
		sqlStatement,
		userId,
		commodityId,
		volume,
	)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStatement = `
		INSERT INTO 
			source_commodities_trade (transaction_id, source_user_id, source_order_id, dest_order_id, broker_id)
		VALUES (?, ?, ?, ?, ?)
	`

	err = tx.QueryRow(
		sqlStatement,
		lastId,
		source.SourceUserId,
		source.SourceOrderId,
		source.DestOrderId,
		source.BrokerId,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}
}
