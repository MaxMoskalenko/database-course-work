package main

import (
	"database-course-work/auth_service"
	"database-course-work/cli_service"
	ex_service "database-course-work/exchanger_service"
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	request := os.Args[1]
	db, err := sql_service.Connect()

	if err != nil {
		fmt.Printf("⛔️ Something went wrong during connection : %s", err.Error())
	}

	// _ init {su_password}
	if request == "init" {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}

		if os.Args[2] != os.Getenv("SUPERUSER_PASSWORD") {
			fmt.Println("⛔️ Superuser password is not correct")
			return
		}
		db.InitDatabase()
		return
	}

	// _ add_licence ${license_code} ${su_password}
	if request == "add_license" {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}

		if os.Args[3] != os.Getenv("SUPERUSER_PASSWORD") {
			fmt.Println("⛔️ Superuser password is not correct")
			return
		}
		err = db.AddLicense(os.Args[2])

		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
		}
		return
	}

	// _ signup_user ${name} ${surname} ${email} ${bank_account} ${password} ${license}
	if request == "signup_user" {
		if len(os.Args) < 8 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 8)
			return
		}
		jwt, err := auth_service.SignUp(
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
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		fmt.Println(jwt)
		return
	}

	// _ signin_user ${email} ${password}
	if request == "signin_user" {
		if len(os.Args) < 4 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 4)
			return
		}
		jwt, err := auth_service.SignIn(
			&db,
			&h.User{
				Email:    os.Args[2],
				Password: os.Args[3],
			},
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		fmt.Println(jwt)
		return
	}

	// _ assign_broker ${email} ${license}
	if request == "assign_broker" {
		if len(os.Args) < 4 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 4)
			return
		}
		err := auth_service.AssignBroker(
			&db,
			&h.User{
				Email: os.Args[2],
				License: &h.License{
					Code: os.Args[3],
				},
			},
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		return
	}

	// _ signup_company ${tag} ${title} ${email} ${phone_number} ${password}
	if request == "signup_company" {
		if len(os.Args) < 7 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 7)
			return
		}
		jwt, err := auth_service.SignUpCompany(
			&db,
			&h.Company{
				Tag:         os.Args[2],
				Title:       os.Args[3],
				Email:       os.Args[4],
				PhoneNumber: os.Args[5],
				Password:    os.Args[6],
			},
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		fmt.Println(jwt)
		return
	}

	// _ signin_company ${tag} ${password}
	if request == "signin_company" {
		if len(os.Args) < 4 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 4)
			return
		}
		jwt, err := auth_service.SignInCompany(
			&db,
			&h.Company{
				Tag:      os.Args[2],
				Password: os.Args[3],
			},
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		fmt.Println(jwt)
		return
	}

	// _ add_commodity ${user_email} ${commodity_label} ${volume} ${company_jwt}
	if request == "add_commodity" {
		if len(os.Args) < 6 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 6)
			return
		}
		volume, err := strconv.ParseFloat(os.Args[4], 64)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		err = ex_service.AddCommodity(
			&db,
			&h.Commodity{
				Label:  os.Args[3],
				Volume: volume,
				Owner: &h.User{
					Email: os.Args[2],
				},
			},
			os.Args[5],
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		return
	}

	// _ check_commodity ${user_jwt}
	if request == "check_commodity" {
		if len(os.Args) < 2 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 2)
			return
		}
		commodities, err := ex_service.CheckCommodities(
			&db,
			os.Args[2],
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		h.PrintCommodities(commodities)
		return
	}

	// _ list_commodities
	if request == "list_commodities" {
		commodities, err := db.GetAvailableCommodities()
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		h.PrintCommoditiesList(commodities)
		return
	}

	// _ add_order ${side} ${commodity_label} ${volume} ${preferable_broker_email} ${user_jwt}
	if request == "add_order" {
		if len(os.Args) < 7 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 7)
			return
		}
		volume, err := strconv.ParseFloat(os.Args[4], 64)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		err = ex_service.AddOrder(
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
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		return
	}

	// _ list_orders_my ${user_jwt}
	if request == "list_orders_my" {
		if len(os.Args) < 3 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 3)
			return
		}
		orders, err := ex_service.ReadUserOrders(
			&db,
			os.Args[2],
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		h.PrintPersonalOrders(orders)
		return
	}

	// _ list_orders_all ${broker_jwt}
	if request == "list_orders_all" {
		if len(os.Args) < 3 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 3)
			return
		}
		orders, err := ex_service.ReadAllOrders(
			&db,
			os.Args[2],
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		h.PrintNativeOrders(orders)
		return
	}

	// _ cancel_order ${order_id} ${user_jwt}
	if request == "cancel_order" {
		if len(os.Args) < 4 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 4)
			return
		}
		orderId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		err = ex_service.CancelOrder(
			&db,
			orderId,
			os.Args[3],
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		return
	}

	// _ execute_order ${first_order_id} ${second_order_id} ${volume} ${broker_jwt}
	if request == "execute_order" {
		if len(os.Args) < 6 {
			fmt.Printf("⛔️ Not enough arguments for %s. Should be %d\n", request, 6)
			return
		}
		firstOrderId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		secondOrderId, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		volume, err := strconv.ParseFloat(os.Args[4], 32)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		err = ex_service.ExecuteOrder(
			&db,
			firstOrderId,
			secondOrderId,
			volume,
			os.Args[5],
		)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		return
	}

	if request == "cli" {
		cli_service.Launch(&db)
		return
	}

	fmt.Println("⛔️ Unknown command")
}
