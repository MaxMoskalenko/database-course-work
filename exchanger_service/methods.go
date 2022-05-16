package exchanger_service

import (
	"database-course-work/auth_service"
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
)

func AddCommodity(
	db *sql_service.Database,
	commodity *h.Commodity,
	companyJWT string,
) {
	company, err := auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !db.CheckIsRecordExist("companies", "tag", company.Tag) {
		fmt.Printf("⛔️ No such company %s\n", company.Tag)
		return
	}

	if !db.CheckIsRecordExist("commodity_types", "label", commodity.Label) {
		fmt.Printf("⛔️ No such commodity type %s\n", commodity.Label)
		return
	}

	if !db.CheckIsRecordExist("users", "email", commodity.Owner.Email) {
		fmt.Printf("⛔️ No such user %s \n", commodity.Owner.Email)
		return
	}

	if commodity.Volume <= 0 {
		panic(fmt.Errorf("⛔️ Company cannot take commodities from user"))
	}

	commodity.Source = &h.CommoditySource{
		Type:      "company",
		CompanyId: db.GetId("companies", "tag", company.Tag),
	}

	commodity.Owner.Id = db.GetId("users", "email", commodity.Owner.Email)
	commodity.Id = db.GetId("commodity_types", "label", commodity.Label)

	db.AddCommodity(commodity)

}

func CheckCommodities(
	db *sql_service.Database,
	userJWT string,
) [](*h.Commodity) {
	user := auth_service.GetUser(db, userJWT)

	return db.GetUserCommodities(user.Email)
}

func AddOrder(
	db *sql_service.Database,
	order *h.Order,
	userJWT string,
) {
	order.Owner = auth_service.GetUser(db, userJWT)

	if order.PrefBroker.Email != "" {
		order.PrefBroker = db.GetUserData(order.PrefBroker.Email)
		if !order.PrefBroker.IsBroker {
			fmt.Println("⛔️ Preferable broker is not a broker")
			order.PrefBroker = &h.User{}
		}
	}

	if !db.CheckIsRecordExist("commodity_types", "label", order.Commodity.Label) {
		panic(fmt.Errorf("⛔️ No such commodity type %s", order.Commodity.Label))
	}

	db.AddOrder(order)
}

func ReadUserOrders(
	db *sql_service.Database,
	userJWT string,
) [](*h.Order) {
	user := auth_service.GetUser(db, userJWT)

	return db.ReadUserOrders(user)
}

func ReadAllOrders(
	db *sql_service.Database,
	brokerJWT string,
) [](*h.Order) {
	broker := auth_service.GetUser(db, brokerJWT)

	if !broker.IsBroker {
		panic(fmt.Errorf("⛔️ Broker is not a broker"))
	}

	return db.ReadAllOrders(broker.Id)
}

func UpdateOrder(
	db *sql_service.Database,
	orderId int,
	newOrder *h.Order,
	userJWT string,
) {
	user := auth_service.GetUser(db, userJWT)

	if db.GetOrderOwnerId(orderId) != user.Id {
		panic(fmt.Errorf("⛔️ Order is not owned by user"))
	}

	if newOrder.PrefBroker.Email != "" {
		newOrder.PrefBroker = db.GetUserData(newOrder.PrefBroker.Email)
		if !newOrder.PrefBroker.IsBroker {
			fmt.Println("⛔️ Preferable broker is not a broker")
			newOrder.PrefBroker = &h.User{}
		}
	}

	if !db.CheckIsRecordExist("commodity_types", "label", newOrder.Commodity.Label) {
		panic(fmt.Errorf("⛔️ No such commodity type %s", newOrder.Commodity.Label))
	}

	db.UpdateOrder(orderId, newOrder)
}

func DeleteOrder(
	db *sql_service.Database,
	orderId int,
	userJWT string,
) {
	user := auth_service.GetUser(db, userJWT)

	if db.GetOrderOwnerId(orderId) != user.Id {
		panic(fmt.Errorf("⛔️ Order is not owned by user"))
	}

	db.DeleteOrder(orderId)
}

func ExecuteOrder(
	db *sql_service.Database,
	firstOrderId int,
	secondOrderId int,
	volume float64,
	brokerJWT string,
) {
	broker := auth_service.GetUser(db, brokerJWT)

	if !broker.IsBroker {
		panic(fmt.Errorf("⛔️ Broker is not a broker"))
	}

	firstOrder := db.GetOrderById(firstOrderId, broker.Id)
	secondOrder := db.GetOrderById(secondOrderId, broker.Id)

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
		db.PerformExchange(
			firstOrder,
			secondOrder,
			volume,
		)
	}

	if secondOrder.Side == "sell" {
		db.PerformExchange(
			secondOrder,
			firstOrder,
			volume,
		)
	}
}
