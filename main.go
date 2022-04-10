package main

import (
	sql_server "database-course-work/sql_server"
	"os"
)

func main() {
	mode := os.Getenv("INPUT_MODE")
	request := os.Args[1]
	db := sql_server.Connect()

	// _ init
	if request == "init" {
		if mode == "api" {
			db.InitCommodityMarket()
		}
	}

	// _ init_exchange ${exchange_name}
	if request == "init_exchange" {
		if mode == "api" {
			db.InitExchange(os.Args[2])
		}
	}
}
