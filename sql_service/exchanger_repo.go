package sql_service

import (
	h "database-course-work/helpers"
	"database/sql"
)

func (db *Database) AddCommodity(commodity *h.Commodity) error {
	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStatement := `
		INSERT INTO commodity_market.commodities_account (owner_id, commodity_id, volume, source)
		VALUES (?, ?, ?, ?);
	`

	res, err := tx.Exec(
		sqlStatement,
		commodity.Owner.Id,
		commodity.Id,
		commodity.Volume,
		commodity.Source.Type,
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

	if commodity.Source.Type == "company" {
		sqlStatement = `
			INSERT INTO commodity_market.source_commodities_company (transaction_id, source_company_id)
			VALUES (?, ?);
		`

		err = tx.QueryRow(
			sqlStatement,
			lastId,
			commodity.Source.CompanyId,
		).Err()

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *Database) GetUserCommodities(userId int) ([](*h.Commodity), error) {
	sqlStatement := `
		SELECT CT.label, SUM(C.volume), CU.unit, U.name, U.surname, U.email 
			FROM commodity_market.commodities_account AS C
		JOIN commodity_market.commodity_types AS CT
			ON CT.id=C.commodity_id
		JOIN commodity_market.units AS CU
			ON CU.id = CT.unit_id 
		JOIN commodity_market.users AS U 
			ON U.id = C.owner_id
		WHERE C.owner_id = ?
		GROUP BY C.owner_id, C.commodity_id;
	`

	rows, err := db.sql.Query(
		sqlStatement,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commodities []*h.Commodity

	for rows.Next() {
		commodity := h.Commodity{
			Owner: &h.User{},
		}
		if err := rows.Scan(
			&commodity.Label,
			&commodity.Volume,
			&commodity.Unit,
			&commodity.Owner.Name,
			&commodity.Owner.Surname,
			&commodity.Owner.Email,
		); err != nil {
			return nil, err
		}
		commodities = append(commodities, &commodity)
	}

	return commodities, nil
}

func (db *Database) GetAvailableCommodities() ([](*h.Commodity), error) {
	sqlStatement := `
		SELECT CT.label, CU.unit
		FROM commodity_market.commodity_types AS CT
		INNER JOIN (
			SELECT unit, id
			FROM units
		) AS CU
			ON CU.id = CT.unit_id;
	`

	rows, err := db.sql.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commodities []*h.Commodity

	for rows.Next() {
		var commodity h.Commodity
		if err := rows.Scan(&commodity.Label, &commodity.Unit); err != nil {
			return nil, err
		}

		commodities = append(commodities, &commodity)
	}

	return commodities, nil
}

func (db *Database) GetUnlockedVolume(userId int, commodityId int) (float64, error) {
	sqlStatement := `
		CALL commodity_market.GetUnlockedVolume(?, ?);
	`

	volume := 0.0

	err := db.sql.QueryRow(
		sqlStatement,
		userId,
		commodityId,
	).Scan(&volume)

	if err != nil {
		return 0, err
	}

	return volume, nil
}

func (db *Database) AddOrder(order *h.Order) error {
	sqlStatement := `
		INSERT INTO commodity_market.orders (owner_id, side, commodity_id, volume, pref_broker_id)
		VALUES (?, ?, ?, ?, ?)
	`

	return db.sql.QueryRow(
		sqlStatement,
		order.Owner.Id,
		order.Side,
		order.Commodity.Id,
		order.Commodity.Volume,
		h.ConvertZeroToNil(order.PrefBroker.Id),
	).Err()
}

func (db *Database) ReadOrders(brokerId int, userId int) ([](*h.Order), error) {
	sqlStatement := `
		SELECT O.id, O.side, O.state, CT.label, CU.unit, O.volume, O.executed_volume, PB.name, PB.surname, U.name, U.surname, U.email
		FROM commodity_market.orders AS O
		JOIN (
			SELECT id, label, unit_id
			FROM commodity_market.commodity_types
		) AS CT
		ON CT.id = O.commodity_id
		JOIN (
			SELECT id, unit
			FROM commodity_market.units
		) AS CU
		ON CU.id = CT.unit_id
		LEFT JOIN (
			SELECT id, name, surname
			FROM commodity_market.users
		) AS PB
		ON O.pref_broker_id IS NOT NULL AND PB.id = O.pref_broker_id
		JOIN (
			SELECT id, name, surname, email
			FROM commodity_market.users
		) AS U
		ON U.id = O.owner_id 
		WHERE 
			(?=0 OR O.owner_id=?) AND
			(?=0 OR O.pref_broker_id=? OR O.pref_broker_id IS NULL) AND
			O.state = 'active';
	`

	rows, err := db.sql.Query(
		sqlStatement,
		userId,
		userId,
		brokerId,
		brokerId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders [](*h.Order)

	for rows.Next() {
		var order h.Order
		var commodity h.Commodity
		var broker, user h.User
		var brokerName, brokerSurname sql.NullString

		if err := rows.Scan(
			&order.Id,
			&order.Side,
			&order.State,
			&commodity.Label,
			&commodity.Unit,
			&commodity.Volume,
			&order.ExecutedVolume,
			&brokerName,
			&brokerSurname,
			&user.Name,
			&user.Surname,
			&user.Email,
		); err != nil {
			return nil, err
		}
		broker.Name = brokerName.String
		broker.Surname = brokerSurname.String
		order.Owner = &user
		order.Commodity = &commodity
		order.PrefBroker = &broker
		orders = append(orders, &order)
	}

	return orders, nil
}

func (db *Database) GetOrderOwnerId(orderId int) (int, error) {
	var ownerId int

	sqlStatement := `
		SELECT owner_id
		FROM commodity_market.orders
		WHERE id=?
	`

	err := db.sql.QueryRow(
		sqlStatement,
		orderId,
	).Scan(&ownerId)

	if err != nil {
		return 0, err
	}
	return ownerId, err
}

func (db *Database) CancelOrder(orderId int) error {
	sqlStatement := `
		UPDATE commodity_market.orders
		SET state = 'canceled'
		WHERE id=?;
	`

	return db.sql.QueryRow(
		sqlStatement,
		orderId,
	).Err()
}

func (db *Database) GetOrderById(orderId int, prefBrokerId int) (*h.Order, error) {
	sqlStatement := `
		SELECT id, side, commodity_id, volume, executed_volume, owner_id, state
		FROM commodity_market.orders
		WHERE 
			(pref_broker_id=? OR pref_broker_id IS NULL) AND  
			id=?;
	`

	order := h.Order{
		Owner:      &h.User{},
		Commodity:  &h.Commodity{},
		PrefBroker: &h.User{},
	}

	err := db.sql.QueryRow(
		sqlStatement,
		prefBrokerId,
		orderId,
	).Scan(
		&order.Id,
		&order.Side,
		&order.Commodity.Id,
		&order.Commodity.Volume,
		&order.ExecutedVolume,
		&order.Owner.Id,
		&order.State,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (db *Database) PerformExchange(
	sellOrder *h.Order,
	buyOrder *h.Order,
	volumeChange float64,
	brokerId int,
) error {
	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = updateTransactionOrder(tx, sellOrder); err != nil {
		return err
	}

	if err = updateTransactionOrder(tx, buyOrder); err != nil {
		return err
	}

	if err = upsertTransactionCommodities(
		tx,
		sellOrder.Owner.Id,
		sellOrder.Commodity.Id,
		-volumeChange,
		&h.CommoditySource{
			SourceUserId:  buyOrder.Owner.Id,
			SourceOrderId: buyOrder.Id,
			DestOrderId:   sellOrder.Id,
			BrokerId:      brokerId,
		},
	); err != nil {
		return err
	}

	if err = upsertTransactionCommodities(
		tx,
		buyOrder.Owner.Id,
		buyOrder.Commodity.Id,
		volumeChange,
		&h.CommoditySource{
			SourceUserId:  sellOrder.Owner.Id,
			SourceOrderId: sellOrder.Id,
			DestOrderId:   buyOrder.Id,
			BrokerId:      brokerId,
		},
	); err != nil {
		return err
	}

	return tx.Commit()
}
