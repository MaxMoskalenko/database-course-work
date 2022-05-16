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
		INSERT INTO commodities_account (owner_user_id, commodity_id, volume, source)
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
	} else if commodity.Source.Type == "trade" {
		sqlStatement = `
			INSERT INTO source_commodities_trade 
				(transaction_id, source_owner_id, source_order_id, source_broker_id)
			VALUES (?, ?, ?, ?);
		`

		err = tx.QueryRow(
			sqlStatement,
			lastId,
			commodity.Source.UserId,
			commodity.Source.OrderId,
			commodity.Source.BrokerId,
		).Err()

		if err != nil {
			tx.Rollback()
			panic(err)
		}
	} else {
		tx.Rollback()
		panic(fmt.Errorf("⛔️ Invalid source type"))
	}

	err = tx.Commit()

	if err != nil {
		panic(err)
	}

}

func (db *Database) GetUserCommodities(userEmail string) [](*h.Commodity) {
	// userId := db.GetId("users", "email", userEmail)

	// sqlStatement := fmt.Sprintf(`
	// 	SELECT CT.label, C.volume, CU.unit
	// 	FROM %s.commodities AS C
	// 	INNER JOIN (
	// 		SELECT label, unit_id, id
	// 		FROM commodity_market.commodity_types
	// 	) AS CT
	// 		ON CT.id=C.commodity_id
	// 	INNER JOIN (
	// 		SELECT unit, id
	// 		FROM commodity_market.units
	// 	) AS CU
	// 		ON CU.id = CT.unit_id
	// 	WHERE user_id=?;
	// `, database)

	// rows, err := db.sql.Query(
	// 	sqlStatement,
	// 	userId,
	// )
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// defer rows.Close()

	var commodities []*h.Commodity

	// for rows.Next() {
	// 	var commodity h.Commodity
	// 	if err := rows.Scan(&commodity.Label, &commodity.Volume, &commodity.Unit); err != nil {
	// 		panic(err)
	// 	}
	// 	commodity.Owner = db.GetUserData(userEmail)
	// 	commodities = append(commodities, &commodity)
	// }

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

func (db *Database) AddOrder(order *h.Order) {
	commodityId := db.GetId("commodity_types", "label", order.Commodity.Label)

	sqlStatement := `
		INSERT INTO orders (owner_id, side, state, commodity_id, volume, pref_broker_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	err := db.sql.QueryRow(
		sqlStatement,
		order.Owner.Id,
		order.Side,
		"active",
		commodityId,
		order.Commodity.Volume,
		h.ConvertZeroToNil(order.PrefBroker.Id),
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) ReadUserOrders(user *h.User) [](*h.Order) {
	sqlStatement := `
		SELECT O.id, O.side, O.state, CT.label, CU.unit, O.volume, PB.name, PB.surname 
		FROM orders AS O
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
			FROM users
		) AS PB
		ON O.pref_broker_id IS NOT NULL AND PB.id = O.pref_broker_id;
	`

	rows, err := db.sql.Query(sqlStatement, user.Id)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var orders [](*h.Order)

	for rows.Next() {
		var order h.Order
		var commodity h.Commodity
		var broker h.User
		var brokerName, brokerSurname sql.NullString

		if err := rows.Scan(
			&order.Id,
			&order.Side,
			&order.State,
			&commodity.Label,
			&commodity.Unit,
			&commodity.Volume,
			&brokerName,
			&brokerSurname,
		); err != nil {
			panic(err)
		}
		broker.Name = brokerName.String
		broker.Surname = brokerSurname.String
		order.Owner = user
		order.Commodity = &commodity
		order.PrefBroker = &broker
		orders = append(orders, &order)
	}

	return orders
}

func (db *Database) ReadAllOrders(brokerId int) [](*h.Order) {
	sqlStatement := `
		SELECT O.id, O.side, O.state, CT.label, CU.unit, O.volume, PB.name, PB.surname, U.name, U.surname, U.email
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
		WHERE O.pref_broker_id=? OR O.pref_broker_id IS NULL;
	`

	rows, err := db.sql.Query(sqlStatement, brokerId)

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

func (db *Database) UpdateOrder(oldId int, newOrder *h.Order) {
	commodityId := db.GetId("commodity_types", "label", newOrder.Commodity.Label)

	sqlStatement := `
		UPDATE orders 
		SET side=?, commodity_id=?, volume=?, pref_broker_id=?
		WHERE id=?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		newOrder.Side,
		commodityId,
		newOrder.Commodity.Volume,
		h.ConvertZeroToNil(newOrder.PrefBroker.Id),
		oldId,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) DeleteOrder(orderId int) {
	sqlStatement := `
		DELETE FROM orders 
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
		SELECT id, side, commodity_id, volume, owner_id
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
) {
	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	updateTransactionOrder(tx, sellOrder)
	updateTransactionOrder(tx, buyOrder)
	upsertTransactionCommodities(tx, sellOrder.Owner.Id, sellOrder.Commodity.Id, -volumeChange)
	upsertTransactionCommodities(tx, buyOrder.Owner.Id, buyOrder.Commodity.Id, volumeChange)

	err = tx.Commit()

	if err != nil {
		panic(err)
	}
}
