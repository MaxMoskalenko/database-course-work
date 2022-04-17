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
	// _ check_commodity_broker ${exchanger_tag} ${user_email} ${broker_jwt}

	if request == "test" {
		fmt.Println("🛠 " + os.Args[2])
	}

	fmt.Println("⛔️ Unknown command")
}
