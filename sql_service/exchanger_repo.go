package sql_service

import (
	h "database-course-work/helpers"
	"database/sql"
	"fmt"
)

func (db *Database) AddCommodity(commodity *h.Commodity) {
	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStatement := `
		INSERT INTO commodities_account (owner_id, commodity_id, volume, source)
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
		panic(err)
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if commodity.Source.Type == "company" {
		sqlStatement = `
			INSERT INTO source_commodities_company (transaction_id, source_company_id)
			VALUES (?, ?);
		`

		err = tx.QueryRow(
			sqlStatement,
			lastId,
			commodity.Source.CompanyId,
		).Err()

		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()

	if err != nil {
		panic(err)
	}

}

func (db *Database) GetUserCommodities(userId int) [](*h.Commodity) {
	sqlStatement := `
		SELECT CT.label, SUM(C.volume), CU.unit, U.name, U.surname, U.email 
			FROM commodities_account AS C
		JOIN commodity_types AS CT
			ON CT.id=C.commodity_id
		JOIN units AS CU
			ON CU.id = CT.unit_id 
		JOIN users AS U 
			ON U.id = C.owner_id
		WHERE C.owner_id = ?
		GROUP BY C.owner_id, C.commodity_id;
	`

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
			panic(err)
		}
		commodities = append(commodities, &commodity)
	}

	return commodities
}

func (db *Database) GetAvailableCommodities() [](*h.Commodity) {
	sqlStatement := `
		SELECT CT.label, CU.unit
		FROM commodity_types AS CT
		INNER JOIN (
			SELECT unit, id
			FROM units
		) AS CU
			ON CU.id = CT.unit_id;
	`

	rows, err := db.sql.Query(sqlStatement)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var commodities []*h.Commodity

	for rows.Next() {
		var commodity h.Commodity
		if err := rows.Scan(&commodity.Label, &commodity.Unit); err != nil {
			panic(err)
		}

		commodities = append(commodities, &commodity)
	}

	return commodities
}

func (db *Database) GetUnlockedVolume(userId int, commodityId int) float64 {
	sqlStatement := `
		CALL GetUnlockedVolume(?, ?);
	`

	volume := 0.0

	err := db.sql.QueryRow(
		sqlStatement,
		userId,
		commodityId,
	).Scan(&volume)

	if err != nil {
		panic(err)
	}

	return volume
}

func (db *Database) AddOrder(order *h.Order) {
	sqlStatement := `
		INSERT INTO orders (owner_id, side, commodity_id, volume, pref_broker_id)
		VALUES (?, ?, ?, ?, ?)
	`

	err := db.sql.QueryRow(
		sqlStatement,
		order.Owner.Id,
		order.Side,
		order.Commodity.Id,
		order.Commodity.Volume,
		h.ConvertZeroToNil(order.PrefBroker.Id),
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) ReadOrders(brokerId int, userId int) [](*h.Order) {
	sqlStatement := `
		SELECT O.id, O.side, O.state, CT.label, CU.unit, O.volume, O.executed_volume, PB.name, PB.surname, U.name, U.surname, U.email
		FROM orders AS O
		JOIN (
			SELECT id, label, unit_id
			FROM commodity_types
		) AS CT
		ON CT.id = O.commodity_id
		JOIN (
			SELECT id, unit
			FROM units
		) AS CU
		ON CU.id = CT.unit_id
		LEFT JOIN (
			SELECT id, name, surname
			FROM users
		) AS PB
		ON O.pref_broker_id IS NOT NULL AND PB.id = O.pref_broker_id
		JOIN (
			SELECT id, name, surname, email
			FROM users
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
		fmt.Println(err.Error())
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
			panic(err)
		}
		broker.Name = brokerName.String
		broker.Surname = brokerSurname.String
		order.Owner = &user
		order.Commodity = &commodity
		order.PrefBroker = &broker
		orders = append(orders, &order)
	}

	return orders
}

func (db *Database) GetOrderOwnerId(orderId int) int {
	var ownerId int

	sqlStatement := `
		SELECT owner_id
		FROM orders
		WHERE id=?
	`

	err := db.sql.QueryRow(
		sqlStatement,
		orderId,
	).Scan(&ownerId)

	if err != nil {
		panic(err)
	}
	return ownerId
}

func (db *Database) CancelOrder(orderId int) {
	sqlStatement := `
		UPDATE orders
		SET state = 'canceled'
		WHERE id=?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		orderId,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) GetOrderById(orderId int, prefBrokerId int) *h.Order {
	sqlStatement := `
		SELECT id, side, commodity_id, volume, executed_volume, owner_id
		FROM orders
		WHERE 
			(pref_broker_id=? OR pref_broker_id IS NULL) AND 
			state='active' AND 
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
	)

	if err != nil {
		panic(err)
	}

	return &order
}

func (db *Database) PerformExchange(
	sellOrder *h.Order,
	buyOrder *h.Order,
	volumeChange float64,
	brokerId int,
) {
	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	updateTransactionOrder(tx, sellOrder)
	updateTransactionOrder(tx, buyOrder)
	upsertTransactionCommodities(
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
	)
	upsertTransactionCommodities(
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
	)

	err = tx.Commit()

	if err != nil {
		panic(err)
	}
}
