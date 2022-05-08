package sql_service

import (
	h "database-course-work/helpers"
	"database/sql"
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
		commodity.Owner = db.GetUserData(database, userEmail)
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
			ON U.id = C.user_id;
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

		commodity.Owner = db.GetUserData(database, email)
		commodities = append(commodities, &commodity)
	}

	return commodities
}

func (db *Database) GetAvailableCommodities() [](*h.Commodity) {
	sqlStatement := `
		SELECT CT.label, CU.unit
		FROM commodity_market.commodity_types AS CT
		INNER JOIN (
			SELECT unit, id
			FROM commodity_market.units
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

func (db *Database) GetCommodityVolume(database string, userId int, commodityLabel string) int {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	commodityId := db.getId("commodity_market", "commodity_types", "label", commodityLabel)

	sqlStatement := fmt.Sprintf(`
		SELECT volume 
		FROM %s.commodities
		WHERE user_id=? AND commodity_id=?;
	`, database)

	var volume int

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

func (db *Database) AddOrder(database string, order *h.Order) {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	commodityId := db.getId("commodity_market", "commodity_types", "label", order.Commodity.Label)

	sqlStatement := fmt.Sprintf(`
		INSERT INTO %s.orders (owner_id, side, state, commodity_id, volume, pref_broker_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, database)

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

func (db *Database) ReadOrders(database string, user *h.User, isOpen bool) [](*h.Order) {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	var filter string

	if isOpen {
		filter = "WHERE O.owner_id=? AND O.state='active'"
	} else {
		filter = "WHERE O.owner_id=?"
	}

	sqlStatement := fmt.Sprintf(`
		SELECT O.id, O.side, O.state, CT.label, CU.unit, O.volume, PB.name, PB.surname 
		FROM %s.orders AS O
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
			FROM %s.users
		) AS PB
		ON O.pref_broker_id IS NOT NULL AND PB.id = O.pref_broker_id
		%s;
	`, database, database, filter)

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

func (db *Database) ReadOrdersNative(database string, brokerId int) [](*h.Order) {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	sqlStatement := fmt.Sprintf(`
		SELECT O.id, O.side, O.state, CT.label, CU.unit, O.volume, PB.name, PB.surname, U.name, U.surname, U.email
		FROM %s.orders AS O
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
			FROM %s.users
		) AS PB
		ON O.pref_broker_id IS NOT NULL AND PB.id = O.pref_broker_id
		JOIN (
			SELECT id, name, surname, email
			FROM %s.users
		) AS U
		ON U.id = O.owner_id 
		WHERE O.pref_broker_id=? OR O.pref_broker_id IS NULL;
	`, database, database, database)

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

func (db *Database) ReadOrdersForeign() [](*h.Order) {
	sqlStatement := `
		CALL commodity_market.GetForeignOrders();
	`

	rows, err := db.sql.Query(sqlStatement)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var orders [](*h.Order)

	for rows.Next() {
		order := h.Order{
			Owner:     &h.User{},
			Commodity: &h.Commodity{},
		}

		if err := rows.Scan(
			&order.Id,
			&order.Side,
			&order.State,
			&order.Commodity.Label,
			&order.Commodity.Unit,
			&order.Commodity.Volume,
			&order.Owner.Name,
			&order.Owner.Surname,
			&order.Owner.Email,
			&order.Owner.ExchangerTag,
		); err != nil {
			panic(err)
		}
		orders = append(orders, &order)
	}

	return orders
}

func (db *Database) GetOrderOwnerId(database string, orderId int) int {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	var ownerId int

	sqlStatement := fmt.Sprintf(`
		SELECT owner_id
		FROM %s.orders
		WHERE id=?
	`, database)

	err := db.sql.QueryRow(
		sqlStatement,
		orderId,
	).Scan(&ownerId)

	if err != nil {
		panic(err)
	}
	return ownerId
}

func (db *Database) UpdateOrder(database string, oldId int, newOrder *h.Order) {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	commodityId := db.getId("commodity_market", "commodity_types", "label", newOrder.Commodity.Label)

	sqlStatement := fmt.Sprintf(`
		UPDATE %s.orders 
		SET side=?, commodity_id=?, volume=?, pref_broker_id=?
		WHERE id=?;
	`, database)

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

func (db *Database) DeleteOrder(database string, orderId int) {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	sqlStatement := fmt.Sprintf(`
		DELETE FROM %s.orders 
		WHERE id=?;
	`, database)

	err := db.sql.QueryRow(
		sqlStatement,
		orderId,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) GetOrderById(database string, orderId int, prefBrokerId int) *h.Order {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	sqlStatement := fmt.Sprintf(`
		SELECT id, side, commodity_id, volume, owner_id
		FROM %s.orders
		WHERE 
			(pref_broker_id=? OR pref_broker_id IS NULL) AND 
			state='active' AND 
			id=?;
	`, database)

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

func (db *Database) PerformNativeExchange(
	database string,
	sellOrder *h.Order,
	buyOrder *h.Order,
	volumeChange float64,
) {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	updateTransactionOrder(tx, database, sellOrder)
	updateTransactionOrder(tx, database, buyOrder)
	upsertTransactionCommodities(tx, database, sellOrder.Owner.Id, sellOrder.Commodity.Id, -volumeChange)
	upsertTransactionCommodities(tx, database, buyOrder.Owner.Id, buyOrder.Commodity.Id, volumeChange)

	err = tx.Commit()

	if err != nil {
		panic(err)
	}
}

func (db *Database) PerformForeignExchange(
	sellDatabase string,
	buyDatabase string,
	sellOrder *h.Order,
	buyOrder *h.Order,
	raceId int,
	volumeChange float64,
) {
	if !h.ValidDatabase(sellDatabase) || !h.ValidDatabase(buyDatabase) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	updateTransactionOrder(tx, sellDatabase, sellOrder)
	updateTransactionOrder(tx, buyDatabase, buyOrder)
	upsertTransactionCommodities(tx, sellDatabase, sellOrder.Owner.Id, sellOrder.Commodity.Id, -volumeChange)
	upsertTransactionCargo(tx, buyDatabase, raceId, buyOrder.Commodity.Id, buyOrder.Owner.Id, volumeChange)

	err = tx.Commit()

	if err != nil {
		panic(err)
	}
}
