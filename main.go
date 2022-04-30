package main

import (
	"database-course-work/auth_service"
	ex_service "database-course-work/exchanger_service"
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
	"os"
	"strconv"
)

func main() {
	request := os.Args[1]
	db := sql_service.Connect()

	// _ init
	if request == "init" {
		db.InitCommodityMarket()
		return
	}

	// _ init_exchange ${database_name} ${exchanger_name} ${tag}
	if request == "init_exchange" {
		db.InitExchange(
			&h.Exchanger{
				DatabaseName: os.Args[2],
				Name:         os.Args[3],
				Tag:          os.Args[4],
			},
		)
		return
	}

	// _ signup_user ${exchanger_tag} ${name} ${surname} ${email} ${password} ${is_broker}
	if request == "signup_user" {
		jwt := auth_service.SignUp(
			&db,
			&h.User{
				ExchangerTag: os.Args[2],
				Name:         os.Args[3],
				Surname:      os.Args[4],
				Email:        os.Args[5],
				Password:     os.Args[6],
				IsBroker:     h.BoolToNumber(os.Args[7] == "true"),
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ signup_broker ${exchanger_tag} ${name} ${surname} ${email} ${password} ${is_broker} ${license}
	if request == "signup_broker" {
		jwt := auth_service.SignUp(
			&db,
			&h.User{
				ExchangerTag: os.Args[2],
				Name:         os.Args[3],
				Surname:      os.Args[4],
				Email:        os.Args[5],
				Password:     os.Args[6],
				IsBroker:     h.BoolToNumber(os.Args[7] == "true"),
				License:      os.Args[8],
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ signup_company ${tag} ${title} ${password}
	if request == "signup_company" {
		jwt := auth_service.SignUpCompany(
			&db,
			&h.Company{
				Tag:      os.Args[2],
				Title:    os.Args[3],
				Password: os.Args[4],
				Type:     "c",
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ signup_shipcompany ${tag} ${title} ${password}
	if request == "signup_shipcompany" {
		jwt := auth_service.SignUpCompany(
			&db,
			&h.Company{
				Tag:      os.Args[2],
				Title:    os.Args[3],
				Password: os.Args[4],
				Type:     "s",
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ signup_broker ${exchanger_tag} ${email} ${password}
	if request == "signin_user" {
		jwt := auth_service.SignIn(
			&db,
			&h.User{
				ExchangerTag: os.Args[2],
				Email:        os.Args[3],
				Password:     os.Args[4],
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ signin_company ${tag} ${password}
	if request == "signin_company" {
		jwt := auth_service.SignInCompany(
			&db,
			&h.Company{
				Tag:      os.Args[2],
				Password: os.Args[3],
				Type:     "c",
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ signin_shipcompany ${tag} ${password}
	if request == "signin_shipcompany" {
		jwt := auth_service.SignInCompany(
			&db,
			&h.Company{
				Tag:      os.Args[2],
				Password: os.Args[3],
				Type:     "s",
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ add_commodity ${exchanger_tag} ${user_email} ${commodity_label} ${volume} ${company_jwt}
	if request == "add_commodity" {
		volume, _ := strconv.Atoi(os.Args[5])
		ex_service.AddCommodity(
			&db,
			&h.User{
				ExchangerTag: os.Args[2],
				Email:        os.Args[3],
			},
			&h.Commodity{
				Label:  os.Args[4],
				Volume: volume,
			},
			os.Args[6],
		)
		return
	}

	// _ check_commodity ${user_jwt}
	if request == "check_commodity" {
		commodities := ex_service.CheckCommodities(
			&db,
			os.Args[2],
		)
		h.PrintCommodities(commodities)
		return
	}

	// _ check_commodity_broker ${exchanger_tag} ${broker_jwt}
	if request == "check_commodity_broker" {
		commodities := ex_service.CheckAllCommodities(
			&db,
			os.Args[2],
			os.Args[3],
		)
		h.PrintCommodities(commodities)
		return
	}

	// _ list_commodities
	if request == "list_commodities" {
		commodities := db.GetAvailableCommodities()
		h.PrintCommoditiesList(commodities)
		return
	}

	// _ add_order ${side} ${commodity_label} ${volume} ${preferable_broker_email} ${user_jwt}
	if request == "add_order" {
		volume, _ := strconv.Atoi(os.Args[4])
		ex_service.AddOrder(
			&db,
			&h.Order{
				Side: os.Args[2],
				Commodity: &h.Commodity{
					Label:  os.Args[3],
					Volume: volume,
				},
				PrefBroker: &h.User{
					Email: os.Args[5],
				},
			},
			os.Args[6],
		)
		return
	}

	// _ list_orders ${isOpen} ${user_jwt}
	if request == "list_orders" {
		orders := ex_service.ReadOrders(
			&db,
			os.Args[2] == "true",
			os.Args[3],
		)
		h.PrintPersonalOrders(orders)
		return
	}

	// _ list_orders_all ${exchanger_tag} ${broker_jwt}
	if request == "list_orders_all" {
		orders := ex_service.ReadOrdersAll(
			&db,
			os.Args[2],
			os.Args[3],
		)
		h.PrintAllOrders(orders)
		return
	}

	// _ update_order ${order_id} ${side} ${commodity_label} ${volume} ${preferable_broker_email} ${user_jwt}
	if request == "update_order" {
		orderId, _ := strconv.Atoi(os.Args[2])
		volume, _ := strconv.Atoi(os.Args[5])
		ex_service.UpdateOrder(
			&db,
			orderId,
			&h.Order{
				Side: os.Args[3],
				Commodity: &h.Commodity{
					Label:  os.Args[4],
					Volume: volume,
				},
				PrefBroker: &h.User{
					Email: os.Args[6],
				},
			},
			os.Args[7],
		)
		return
	}

	// _ delete_order ${order_id} ${user_jwt}
	if request == "delete_order" {
		orderId, _ := strconv.Atoi(os.Args[2])
		ex_service.DeleteOrder(
			&db,
			orderId,
			os.Args[3],
		)
		return
	}

	if request == "test" {
		fmt.Println("🛠 " + os.Args[2])
	}

	fmt.Println("⛔️ Unknown command")
}
