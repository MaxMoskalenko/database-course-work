package sql_service

import (
	h "database-course-work/helpers"
	"database/sql"
)

func updateTransactionOrder(
	tx *sql.Tx,
	order *h.Order,
) error {
	sqlStatement := `
		UPDATE commodity_market.orders
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
		return err
	}

	return nil
}

func upsertTransactionCommodities(
	tx *sql.Tx,
	userId int,
	commodityId int,
	volume float64,
	source *h.CommoditySource,
) error {
	sqlStatement := `
		INSERT INTO
			commodity_market.commodities_account (owner_id, commodity_id, volume, source)
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
		return err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStatement = `
		INSERT INTO 
			commodity_market.source_commodities_trade (transaction_id, source_order_id, dest_order_id, broker_id)
		VALUES (?, ?, ?, ?)
	`

	err = tx.QueryRow(
		sqlStatement,
		lastId,
		source.SourceOrderId,
		source.DestOrderId,
		source.BrokerId,
	).Err()

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
