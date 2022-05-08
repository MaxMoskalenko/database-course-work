package exchanger_service

import (
	"database-course-work/auth_service"
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
)

func AddCommodity(
	db *sql_service.Database,
	user *h.User,
	commodity *h.Commodity,
	companyJWT string,
) {
	company, err := auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	database := db.GetDatabaseByTag(user.ExchangerTag)

	if len(database) == 0 {
		fmt.Printf("⛔️ Database with a %s tag does not exist\n", user.ExchangerTag)
		return
	}

	if !db.CheckIsRecordExist("commodity_market", h.GetTableFromType(company.Type), "tag", company.Tag) {
		fmt.Printf("⛔️ No such company %s\n", company.Tag)
		return
	}

	if !db.CheckIsRecordExist("commodity_market", "commodity_types", "label", commodity.Label) {
		fmt.Printf("⛔️ No such commodity type %s\n", commodity.Label)
		return
	}

	if !db.CheckIsRecordExist(database, "users", "email", user.Email) {
		fmt.Printf("⛔️ No such user %s \n", user.Email)
		return
	}

	db.AddCommodity(database, user.Email, commodity)

}

func CheckCommodities(
	db *sql_service.Database,
	userJWT string,
) [](*h.Commodity) {
	user := auth_service.GetUser(db, userJWT)

	database := db.GetDatabaseByTag(user.ExchangerTag)

	return db.GetUserCommodities(database, user.Email)
}

func CheckAllCommodities(
	db *sql_service.Database,
	exchangerTag string,
	brokerJWT string,
) [](*h.Commodity) {
	broker := auth_service.GetUser(db, brokerJWT)
	database := db.GetDatabaseByTag(exchangerTag)

	if broker.IsBroker != 1 {
		panic(fmt.Errorf("⛔️ User is not a broker"))
	}

	return db.GetAllCommodities(database)
}

func AddOrder(
	db *sql_service.Database,
	order *h.Order,
	userJWT string,
) {
	order.Owner = auth_service.GetUser(db, userJWT)
	database := db.GetDatabaseByTag(order.Owner.ExchangerTag)

	if order.PrefBroker.Email != "" {
		order.PrefBroker = db.GetUserData(database, order.PrefBroker.Email)
		if order.PrefBroker.IsBroker != 1 {
			fmt.Println("⛔️ Preferable broker is not a broker")
			order.PrefBroker = &h.User{}
		}
	}

	if !db.CheckIsRecordExist("commodity_market", "commodity_types", "label", order.Commodity.Label) {
		panic(fmt.Errorf("⛔️ No such commodity type %s", order.Commodity.Label))
	}

	db.AddOrder(database, order)
}

func ReadOrders(
	db *sql_service.Database,
	isOpen bool,
	userJWT string,
) [](*h.Order) {
	user := auth_service.GetUser(db, userJWT)
	database := db.GetDatabaseByTag(user.ExchangerTag)

	return db.ReadOrders(database, user, isOpen)
}

func ReadOrdersNative(
	db *sql_service.Database,
	exchangerTag string,
	brokerJWT string,
) [](*h.Order) {
	broker := auth_service.GetUser(db, brokerJWT)
	database := db.GetDatabaseByTag(exchangerTag)

	if broker.IsBroker != 1 {
		panic(fmt.Errorf("⛔️ Broker is not a broker"))
	}

	return db.ReadOrdersNative(database, broker.Id)
}

func ReadOrdersForeign(
	db *sql_service.Database,
	brokerJWT string,
) [](*h.Order) {
	broker := auth_service.GetUser(db, brokerJWT)

	if broker.IsBroker != 1 {
		panic(fmt.Errorf("⛔️ Broker is not a broker"))
	}

	return db.ReadOrdersForeign()

}

func UpdateOrder(
	db *sql_service.Database,
	orderId int,
	newOrder *h.Order,
	userJWT string,
) {
	user := auth_service.GetUser(db, userJWT)
	database := db.GetDatabaseByTag(user.ExchangerTag)

	if db.GetOrderOwnerId(database, orderId) != user.Id {
		panic(fmt.Errorf("⛔️ Order is not owned by user"))
	}

	if newOrder.PrefBroker.Email != "" {
		newOrder.PrefBroker = db.GetUserData(database, newOrder.PrefBroker.Email)
		if newOrder.PrefBroker.IsBroker != 1 {
			fmt.Println("⛔️ Preferable broker is not a broker")
			newOrder.PrefBroker = &h.User{}
		}
	}

	if !db.CheckIsRecordExist("commodity_market", "commodity_types", "label", newOrder.Commodity.Label) {
		panic(fmt.Errorf("⛔️ No such commodity type %s", newOrder.Commodity.Label))
	}

	db.UpdateOrder(database, orderId, newOrder)
}

func DeleteOrder(
	db *sql_service.Database,
	orderId int,
	userJWT string,
) {
	user := auth_service.GetUser(db, userJWT)
	database := db.GetDatabaseByTag(user.ExchangerTag)

	if db.GetOrderOwnerId(database, orderId) != user.Id {
		panic(fmt.Errorf("⛔️ Order is not owned by user"))
	}

	db.DeleteOrder(database, orderId)
}

func ExecuteNativeOrder(
	db *sql_service.Database,
	firstOrderId int,
	secondOrderId int,
	volume float64,
	brokerJWT string,
) {
	broker := auth_service.GetUser(db, brokerJWT)
	database := db.GetDatabaseByTag(broker.ExchangerTag)

	if broker.IsBroker != 1 {
		panic(fmt.Errorf("⛔️ Broker is not a broker"))
	}

	firstOrder := db.GetOrderById(database, firstOrderId, broker.Id)
	secondOrder := db.GetOrderById(database, secondOrderId, broker.Id)

	if firstOrder.Owner.Id == secondOrder.Owner.Id {
		panic(fmt.Errorf("⛔️ Order owners are the same"))
	}

	if firstOrder.Side == secondOrder.Side {
		panic(fmt.Errorf("⛔️ Orders have the same side"))
	}

	if firstOrder.Commodity.Volume < volume || secondOrder.Commodity.Volume < volume {
		panic(fmt.Errorf("⛔️ Executable volume is bigger than one of the order`s volume"))
	}
	firstOrder.Commodity.Volume -= volume
	secondOrder.Commodity.Volume -= volume

	if firstOrder.Commodity.Volume == 0 {
		firstOrder.State = "executed"
	} else {
		firstOrder.State = "active"
	}

	if secondOrder.Commodity.Volume == 0 {
		secondOrder.State = "executed"
	} else {
		secondOrder.State = "active"
	}

	if firstOrder.Side == "sell" {
		db.PerformNativeExchange(
			database,
			firstOrder,
			secondOrder,
			volume,
		)
	}

	if secondOrder.Side == "sell" {
		db.PerformNativeExchange(
			database,
			secondOrder,
			firstOrder,
			volume,
		)
	}
}

func ExecuteForeignOrder(
	db *sql_service.Database,
	firstOrder *h.Order,
	secondOrder *h.Order,
	raceId int,
	volume float64,
	brokerJWT string,
) {
	broker := auth_service.GetUser(db, brokerJWT)

	if broker.IsBroker != 1 {
		panic(fmt.Errorf("⛔️ Broker is not a broker"))
	}

	if firstOrder.Exchnager.Tag == secondOrder.Exchnager.Tag {
		panic(fmt.Errorf("⛔️ Exchnagers are the same, please use native exchange"))
	}

	race := db.GetRaceById(raceId)

	firstDatabase := db.GetDatabaseByTag(firstOrder.Exchnager.Tag)
	secondDatabase := db.GetDatabaseByTag(secondOrder.Exchnager.Tag)

	firstOrder = db.GetOrderById(firstDatabase, firstOrder.Id, broker.Id)
	secondOrder = db.GetOrderById(secondDatabase, secondOrder.Id, broker.Id)

	if firstOrder.Side == secondOrder.Side {
		panic(fmt.Errorf("⛔️ Orders have the same side"))
	}

	if (firstOrder.Side == "sell" && db.GetTagByDatabase(firstDatabase) != race.FromExch.Tag) ||
		(firstOrder.Side == "buy" && db.GetTagByDatabase(firstDatabase) != race.ToExch.Tag) {
		panic(fmt.Errorf("⛔️ First exchnager is invalid for this race"))
	}

	if (secondOrder.Side == "sell" && db.GetTagByDatabase(secondDatabase) != race.FromExch.Tag) ||
		(secondOrder.Side == "buy" && db.GetTagByDatabase(secondDatabase) != race.ToExch.Tag) {
		panic(fmt.Errorf("⛔️ Second exchnager is invalid for this race"))
	}

	if firstOrder.Commodity.Volume < volume || secondOrder.Commodity.Volume < volume {
		panic(fmt.Errorf("⛔️ Executable volume is bigger than one of the order`s volume"))
	}
	firstOrder.Commodity.Volume -= volume
	secondOrder.Commodity.Volume -= volume

	if firstOrder.Commodity.Volume == 0 {
		firstOrder.State = "executed"
	} else {
		firstOrder.State = "active"
	}

	if secondOrder.Commodity.Volume == 0 {
		secondOrder.State = "executed"
	} else {
		secondOrder.State = "active"
	}

	if firstOrder.Side == "sell" {
		db.PerformForeignExchange(
			firstDatabase,
			secondDatabase,
			firstOrder,
			secondOrder,
			raceId,
			volume,
		)
	}

	if secondOrder.Side == "sell" {
		db.PerformForeignExchange(
			secondDatabase,
			firstDatabase,
			secondOrder,
			firstOrder,
			raceId,
			volume,
		)
	}
}
