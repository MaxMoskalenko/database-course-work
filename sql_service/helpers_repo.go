package sql_service

import (
	h "database-course-work/helpers"
	"fmt"
)

func (db *Database) CheckIsRecordExist(table string, column string, value string) bool {
	if !h.ValidTable(table) || !h.ValidTable(column) {
		fmt.Println("🛠 Invalid table/column name")
		return false
	}

	var result string
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?;", column, table, column)

	db.sql.QueryRow(
		sqlStatement,
		value,
	).Scan(&result)

	return len(result) != 0
}

func (db *Database) GetId(table string, column string, value string) int {
	if !h.ValidTable(table) || !h.ValidTable(column) {
		panic(fmt.Errorf("🛠 Invalid table or column name"))
	}

	id := 0
	sqlStatement := fmt.Sprintf("SELECT id FROM %s WHERE %s=?;", table, column)
	err := db.sql.QueryRow(
		sqlStatement,
		value,
	).Scan(&id)

	if err != nil {
		panic(err)
	}

	return id

}

func (db *Database) GetLicense(code string) *h.License {
	sqlStatement := `
		SELECT id, license_code, is_taken
		FROM licenses
		WHERE license_code=?
	`

	var license h.License

	err := db.sql.QueryRow(
		sqlStatement,
		code,
	).Scan(&license.Id, &license.Code, &license.IsTaken)

	if err != nil {
		panic(err)
	}

	return &license
}
