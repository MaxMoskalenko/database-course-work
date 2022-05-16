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
	sqlStatement := fmt.Sprintf("SELECT %s FROM commodity_market.%s WHERE %s=?;", column, table, column)

	db.sql.QueryRow(
		sqlStatement,
		value,
	).Scan(&result)

	return len(result) != 0
}

func (db *Database) GetId(table string, column string, value string) (int, error) {
	if !h.ValidTable(table) || !h.ValidTable(column) {
		fmt.Println("🛠 Invalid table or column name")
		return 0, nil
	}

	id := 0
	sqlStatement := fmt.Sprintf("SELECT id FROM commodity_market.%s WHERE %s=?;", table, column)
	err := db.sql.QueryRow(
		sqlStatement,
		value,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (db *Database) GetLicense(code string) (*h.License, error) {
	sqlStatement := `
		SELECT id, license_code, is_taken
		FROM commodity_market.licenses
		WHERE license_code=?
	`

	var license h.License

	err := db.sql.QueryRow(
		sqlStatement,
		code,
	).Scan(&license.Id, &license.Code, &license.IsTaken)

	if err != nil {
		return nil, err
	}

	return &license, err
}
