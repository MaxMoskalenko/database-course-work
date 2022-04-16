package main

import (
	auth_service "database-course-work/auth_service"
	h "database-course-work/helpers"
	sql_service "database-course-work/sql_service"
	"fmt"
	"os"
)

func main() {
	mode := os.Getenv("INPUT_MODE")
	request := os.Args[1]
	db := sql_service.Connect()

	// _ init
	if request == "init" {
		if mode == "api" {
			db.InitCommodityMarket()
		}
		return
	}

	// _ init_exchange ${database_name} ${exchanger_name} ${tag}
	if request == "init_exchange" {
		if mode == "api" {
			db.InitExchange(
				&h.Exchanger{
					DatabaseName: os.Args[2],
					Name:         os.Args[3],
					Tag:          os.Args[4],
				},
			)
		}
		return
	}

	// _ signup_user ${exchanger_tag} ${name} ${surname} ${email} ${password} ${is_broker}
	if request == "signup_user" {
		if mode == "api" {
			jwt := auth_service.SignUp(
				&db,
				os.Args[2],
				&h.User{
					Name:     os.Args[3],
					Surname:  os.Args[4],
					Email:    os.Args[5],
					Password: os.Args[6],
					IsBroker: h.BoolToNumber(os.Args[7] == "true"),
				},
			)
			fmt.Println(jwt)
		}
		return
	}

	// _ signup_broker ${exchanger_tag} ${name} ${surname} ${email} ${password} ${is_broker} ${license}
	if request == "signup_broker" {
		if mode == "api" {
			jwt := auth_service.SignUp(
				&db,
				os.Args[2],
				&h.User{
					Name:     os.Args[3],
					Surname:  os.Args[4],
					Email:    os.Args[5],
					Password: os.Args[6],
					IsBroker: h.BoolToNumber(os.Args[7] == "true"),
					License:  os.Args[8],
				},
			)
			fmt.Println(jwt)
		}
		return
	}

	// _ signup_company ${tag} ${title} ${password}
	if request == "signup_company" {
		if mode == "api" {
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
		}
		return
	}

	// _ signup_shipcompany ${tag} ${title} ${password}
	if request == "signup_shipcompany" {
		if mode == "api" {
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
		}
		return
	}

	// _ signup_broker ${exchanger_tag} ${email} ${password}
	if request == "signin_user" {
		if mode == "api" {
			jwt := auth_service.SignIn(
				&db,
				os.Args[2],
				&h.User{
					Email:    os.Args[3],
					Password: os.Args[4],
				},
			)
			fmt.Println(jwt)
		}
		return
	}

	// _ signin_company ${tag} ${password}
	if request == "signin_company" {
		if mode == "api" {
			jwt := auth_service.SignInCompany(
				&db,
				&h.Company{
					Tag:      os.Args[2],
					Password: os.Args[3],
					Type:     "c",
				},
			)
			fmt.Println(jwt)
		}
		return
	}

	// _ signin_shipcompany ${tag} ${password}
	if request == "signin_shipcompany" {
		if mode == "api" {
			jwt := auth_service.SignInCompany(
				&db,
				&h.Company{
					Tag:      os.Args[2],
					Password: os.Args[3],
					Type:     "s",
				},
			)
			fmt.Println(jwt)
		}
		return
	}

	if request == "test" {
		fmt.Println("🛠 " + os.Args[2])
	}

	fmt.Println("⛔️ Unknown command")
}
