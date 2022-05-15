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
		db.InitDatabase()
		return
	}

	// _ signup_user ${name} ${surname} ${email} ${bank_account} ${password} ${license}
	if request == "signup_user" {
		if len(os.Args) < 8 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 8)
		}
		jwt := auth_service.SignUp(
			&db,
			&h.User{
				Name:        os.Args[2],
				Surname:     os.Args[3],
				Email:       os.Args[4],
				BankAccount: os.Args[5],
				Password:    os.Args[6],
				License: &h.License{
					Code: os.Args[7],
				},
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ signup_company ${tag} ${title} ${email} ${phone_number} ${password}
	if request == "signup_company" {
		if len(os.Args) < 7 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 7)
		}
		jwt := auth_service.SignUpCompany(
			&db,
			&h.Company{
				Tag:         os.Args[2],
				Title:       os.Args[3],
				Email:       os.Args[4],
				PhoneNumber: os.Args[5],
				Password:    os.Args[6],
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ assign_broker ${email} ${license}
	if request == "assign_broker" {
		if len(os.Args) < 4 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 4)
		}
		auth_service.AssignBroker(
			&db,
			&h.User{
				Email: os.Args[2],
				License: &h.License{
					Code: os.Args[3],
				},
			},
		)
		return
	}

	// _ signin_company ${tag} ${password}
	if request == "signin_company" {
		if len(os.Args) < 4 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 4)
		}
		jwt := auth_service.SignInCompany(
			&db,
			&h.Company{
				Tag:      os.Args[2],
				Password: os.Args[3],
			},
		)
		fmt.Println(jwt)
		return
	}

	// _ add_commodity ${user_email} ${commodity_label} ${volume} ${company_jwt}
	if request == "add_commodity" {
		if len(os.Args) < 6 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 6)
		}
		volume, _ := strconv.ParseFloat(os.Args[4], 64)
		ex_service.AddCommodity(
			&db,
			&h.User{
				Email: os.Args[2],
			},
			&h.Commodity{
				Label:  os.Args[3],
				Volume: volume,
			},
			os.Args[5],
		)
		return
	}

	// _ check_commodity ${user_jwt}
	if request == "check_commodity" {
		if len(os.Args) < 2 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 2)
		}
		commodities := ex_service.CheckCommodities(
			&db,
			os.Args[2],
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
		if len(os.Args) < 7 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 7)
		}
		volume, _ := strconv.ParseFloat(os.Args[4], 64)
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

	// _ cancel_order ${order_id} ${user_jwt}

	// _ list_orders_my ${user_jwt}
	if request == "list_orders_my" {
		if len(os.Args) < 3 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 3)
		}
		orders := ex_service.ReadUserOrders(
			&db,
			os.Args[2],
		)
		h.PrintPersonalOrders(orders)
		return
	}

	// _ list_orders_all ${broker_jwt}
	if request == "list_orders_all" {
		if len(os.Args) < 3 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 3)
		}
		orders := ex_service.ReadAllOrders(
			&db,
			os.Args[2],
		)
		h.PrintNativeOrders(orders)
		return
	}

	// _ update_order ${order_id} ${side} ${commodity_label} ${volume} ${preferable_broker_email} ${user_jwt}
	if request == "update_order" {
		if len(os.Args) < 8 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 8)
		}
		orderId, _ := strconv.Atoi(os.Args[2])
		volume, _ := strconv.ParseFloat(os.Args[5], 64)
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
		if len(os.Args) < 4 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 4)
		}
		orderId, _ := strconv.Atoi(os.Args[2])
		ex_service.DeleteOrder(
			&db,
			orderId,
			os.Args[3],
		)
		return
	}

	// _ execute_order ${first_order_id} ${second_order_id} ${volume} ${broker_jwt}
	if request == "execute_order" {
		if len(os.Args) < 6 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 6)
		}
		firstOrderId, _ := strconv.Atoi(os.Args[2])
		secondOrderId, _ := strconv.Atoi(os.Args[3])
		volume, _ := strconv.ParseFloat(os.Args[4], 32)
		ex_service.ExecuteOrder(
			&db,
			firstOrderId,
			secondOrderId,
			volume,
			os.Args[5],
		)
		return
	}

	if request == "test" {
		fmt.Println("🛠 " + os.Args[2])
	}

	fmt.Println("⛔️ Unknown command")
}
