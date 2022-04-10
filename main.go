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
	}

	// _ init_exchange ${database_name} ${exchanger_name} ${tag}
	if request == "init_exchange" {
		if mode == "api" {
			db.InitExchange(os.Args[2], os.Args[3], os.Args[4])
		}
	}

	// _ signup_user ${exchanger_tag} ${name} ${surname} ${email} ${password} ${is_broker}
	if request == "signup_user" {
		if mode == "api" {
			auth_service.SignUp(
				&db,
				os.Args[2],
				&h.User{
					os.Args[3],
					os.Args[4],
					os.Args[5],
					os.Args[6],
					h.BoolToNumber(os.Args[7] == "true"),
					"",
				},
			)
		}
	}

	// _ signup_broker ${exchanger_tag} ${name} ${surname} ${email} ${password} ${is_broker} ${license}
	if request == "signup_broker" {
		if mode == "api" {
			auth_service.SignUp(
				&db,
				os.Args[2],
				&h.User{
					os.Args[3],
					os.Args[4],
					os.Args[5],
					os.Args[6],
					h.BoolToNumber(os.Args[7] == "true"),
					os.Args[8],
				},
			)
		}
	}

	// _ signup_broker ${exchanger_tag} ${email} ${password}
	if request == "signin_user" {
		if mode == "api" {
			jwt := auth_service.SignIn(
				&db,
				os.Args[2],
				&h.User{
					"",
					"",
					os.Args[3],
					os.Args[4],
					0,
					"",
				},
			)
			fmt.Println(jwt)
		}
	}

	if request == "test" {
		fmt.Println("🛠 " + os.Args[2])
	}
}
