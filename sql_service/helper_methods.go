package sql_service

import (
	h "database-course-work/helpers"
	"fmt"
)

func (db *Database) GetDatabaseByTag(tag string) string {
	var dbName string
	db.sql.Exec("USE commodity_market;")

	sqlStatement := "SELECT database_name FROM exchangers WHERE tag=?;"

	err := db.sql.QueryRow(
		sqlStatement,
		tag,
	).Scan(&dbName)

	if err != nil {
		panic(err)
	}

	return dbName
}

func (db *Database) CheckIsRecordExist(database string, table string, column string, value string) bool {
	if !h.ValidDatabase(database) {
		fmt.Println("🛠 Invalid database name")
		return false
	}

	if !h.ValidTable(table) || !h.ValidTable(column) {
		fmt.Println("🛠 Invalid table/column name")
		return false
	}

	db.sql.Exec(fmt.Sprintf("USE %s;", database))

	var result string
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?;", column, table, column)

	db.sql.QueryRow(
		sqlStatement,
		value,
	).Scan(&result)

	return len(result) != 0
}

func (db *Database) getId(database string, table string, column string, value string) int {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	if !h.ValidTable(table) || !h.ValidTable(column) {
		panic(fmt.Errorf("🛠 Invalid table or column name"))
	}

	id := 0
	sqlStatement := fmt.Sprintf("SELECT id FROM %s.%s WHERE %s=?;", database, table, column)
	err := db.sql.QueryRow(
		sqlStatement,
		value,
	).Scan(&id)

	if err != nil {
		panic(err)
	}

	return id

}
