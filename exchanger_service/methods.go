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
) error {
	company, err := auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		return err
	}

	if !db.CheckIsRecordExist("companies", "tag", company.Tag) {
		return fmt.Errorf("no such company %s", company.Tag)
	}

	if !db.CheckIsRecordExist("commodity_types", "label", commodity.Label) {
		return fmt.Errorf("no such commodity type %s", commodity.Label)
	}

	if !db.CheckIsRecordExist("users", "email", commodity.Owner.Email) {
		return fmt.Errorf("no such user %s", commodity.Owner.Email)
	}

	if commodity.Volume <= 0 {
		return fmt.Errorf("company cannot take commodities from user")
	}

	companyId, err := db.GetId("companies", "tag", company.Tag)

	if err != nil {
		return err
	}

	commodity.Source = &h.CommoditySource{
		Type:      "company",
		CompanyId: companyId,
	}

	commodity.Owner.Id, err = db.GetId("users", "email", commodity.Owner.Email)
	if err != nil {
		return err
	}

	commodity.Id, err = db.GetId("commodity_types", "label", commodity.Label)
	if err != nil {
		return err
	}

	return db.AddCommodity(commodity)
}

func CheckCommodities(
	db *sql_service.Database,
	userJWT string,
) ([](*h.Commodity), error) {
	user, err := auth_service.GetUser(db, userJWT)
	if err != nil {
		return nil, err
	}

	return db.GetUserCommodities(user.Id)
}

func AddOrder(
	db *sql_service.Database,
	order *h.Order,
	userJWT string,
) error {
	var err error
	order.Owner, err = auth_service.GetUser(db, userJWT)

	if err != nil {
		return err
	}

	if order.PrefBroker.Email != "" {
		order.PrefBroker, err = db.GetUserData(order.PrefBroker.Email)
		if err != nil {
			return err
		}
		if !order.PrefBroker.IsBroker {
			fmt.Println("⛔️ preferable broker is not a broker")
			order.PrefBroker = &h.User{}
		}
	}

	if !db.CheckIsRecordExist("commodity_types", "label", order.Commodity.Label) {
		return fmt.Errorf("no such commodity type %s", order.Commodity.Label)
	}

	order.Commodity.Id, err = db.GetId("commodity_types", "label", order.Commodity.Label)
	if err != nil {
		return err
	}

	unlockedVolume, err := db.GetUnlockedVolume(order.Owner.Id, order.Commodity.Id)
	if err != nil {
		return err
	}

	if order.Side == "sell" && unlockedVolume < order.Commodity.Volume {
		return fmt.Errorf("you have insufficient amount to sell")
	}

	return db.AddOrder(order)
}

func ReadUserOrders(
	db *sql_service.Database,
	userJWT string,
) ([](*h.Order), error) {
	user, err := auth_service.GetUser(db, userJWT)

	if err != nil {
		return nil, err
	}

	return db.ReadOrders(0, user.Id)
}

func ReadAllOrders(
	db *sql_service.Database,
	brokerJWT string,
) ([](*h.Order), error) {
	broker, err := auth_service.GetUser(db, brokerJWT)

	if err != nil {
		return nil, err
	}

	if !broker.IsBroker {
		return nil, fmt.Errorf("broker is not a broker")
	}

	return db.ReadOrders(broker.Id, 0)
}

func CancelOrder(
	db *sql_service.Database,
	orderId int,
	userJWT string,
) error {
	user, err := auth_service.GetUser(db, userJWT)
	if err != nil {
		return err
	}

	orderOwnerId, err := db.GetOrderOwnerId(orderId)
	if err != nil {
		return err
	}

	if orderOwnerId != user.Id {
		return fmt.Errorf("order is not owned by user")
	}

	return db.CancelOrder(orderId)
}

func ExecuteOrder(
	db *sql_service.Database,
	firstOrderId int,
	secondOrderId int,
	volume float64,
	brokerJWT string,
) error {
	broker, err := auth_service.GetUser(db, brokerJWT)
	if err != nil {
		return err
	}

	if !broker.IsBroker {
		return fmt.Errorf("broker is not a broker")
	}

	firstOrder, err := db.GetOrderById(firstOrderId, broker.Id)
	if err != nil {
		return err
	}

	secondOrder, err := db.GetOrderById(secondOrderId, broker.Id)
	if err != nil {
		return err
	}

	if firstOrder.State != "active" || secondOrder.State != "active" {
		return fmt.Errorf("one of the orders is not active")
	}

	if firstOrder.Owner.Id == secondOrder.Owner.Id {
		return fmt.Errorf("order owners are the same")
	}

	if firstOrder.Side == secondOrder.Side {
		return fmt.Errorf("orders have the same side")
	}

	if firstOrder.Commodity.Volume-firstOrder.ExecutedVolume < volume ||
		secondOrder.Commodity.Volume-secondOrder.ExecutedVolume < volume {
		return fmt.Errorf("executable volume is bigger than one of the order`s volume")
	}
	firstOrder.ExecutedVolume += volume
	secondOrder.ExecutedVolume += volume

	if firstOrder.Commodity.Volume == firstOrder.ExecutedVolume {
		firstOrder.State = "executed"
	} else {
		firstOrder.State = "active"
	}

	if secondOrder.Commodity.Volume == secondOrder.ExecutedVolume {
		secondOrder.State = "executed"
	} else {
		secondOrder.State = "active"
	}

	if firstOrder.Side == "sell" {
		err = db.PerformExchange(
			firstOrder,
			secondOrder,
			volume,
			broker.Id,
		)
		if err != nil {
			return err
		}
	}

	if secondOrder.Side == "sell" {
		err = db.PerformExchange(
			secondOrder,
			firstOrder,
			volume,
			broker.Id,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
