package sql_server

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	gorm *gorm.DB
}

func Connect() Database {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	USER := os.Getenv("MYSQL_DB_USER")
	PASS := os.Getenv("MYSQL_DB_PASSWORD")
	HOST := os.Getenv("MYSQL_DB_HOST")
	PORT := os.Getenv("MYSQL_DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", USER, PASS, HOST, PORT)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.Exec("CREATE DATABASE IF NOT EXISTS commodity_market;")

	db.Exec("USE commodity_market;")

	return Database{db}
}
