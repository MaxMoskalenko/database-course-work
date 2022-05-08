package sql_service

import (
	h "database-course-work/helpers"
	"fmt"
)

func (db *Database) InsertRace(race *h.Race) {
	sqlStatement := `
		SELECT from_id, id AS to_id, company_id
		FROM commodity_market.exchangers AS e1
		LEFT JOIN (
			SELECT id AS from_id, tag
			FROM commodity_market.exchangers 
		) AS e2
		ON e2.tag = ?
		LEFT JOIN (
			SELECT id AS company_id, tag
			FROM commodity_market.shipment_companies 
		) AS sh
		ON sh.tag = ?
		WHERE e1.tag = ?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		race.FromExch.Tag,
		race.Company.Tag,
		race.ToExch.Tag,
	).Scan(&race.FromExch.Id, &race.ToExch.Id, &race.Company.Id)

	if err != nil {
		panic(err)
	}

	sqlStatement = `
		INSERT INTO commodity_market.races (from_id, to_id, race_date, company_id)
		VALUES (?, ?, ?, ?);
	`

	err = db.sql.QueryRow(
		sqlStatement,
		race.FromExch.Id,
		race.ToExch.Id,
		race.DateStamp,
		race.Company.Id,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) ReadRaces() [](*h.Race) {
	sqlStatement := `	
		SELECT R.id, FE.tag AS from_tag, TE.tag AS to_tag, R.race_date, SHC.tag AS company_tag
		FROM commodity_market.races AS R
		JOIN (
			SELECT id, tag
			FROM commodity_market.exchangers 
		) AS FE
		ON FE.id = R.from_id 
		JOIN (
			SELECT id, tag
			FROM commodity_market.exchangers
		) AS TE
		ON TE.id = R.to_id 
		JOIN (
			SELECT id, tag
			FROM commodity_market.shipment_companies
		) AS SHC
		ON SHC.id = R.company_id;
	`

	rows, err := db.sql.Query(sqlStatement)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var races [](*h.Race)

	for rows.Next() {
		race := h.Race{
			FromExch: &h.Exchanger{},
			ToExch:   &h.Exchanger{},
			Company:  &h.Company{},
		}
		if err := rows.Scan(
			&race.Id,
			&race.FromExch.Tag,
			&race.ToExch.Tag,
			&race.DateStamp,
			&race.Company.Tag,
		); err != nil {
			panic(err)
		}
		races = append(races, &race)
	}
	return races
}

func (db *Database) GetCompanyTag(raceId int) string {
	sqlStatement := `
		SELECT SHC.tag
		FROM commodity_market.races AS R
		JOIN (
			SELECT id, tag
			FROM commodity_market.shipment_companies
		) AS SHC
		ON R.id = ?;
	`

	tag := ""
	err := db.sql.QueryRow(
		sqlStatement,
		raceId,
	).Scan(&tag)

	if err != nil {
		panic(err)
	}

	return tag
}

func (db *Database) UpdateRace(id int, race *h.Race) {
	sqlStatement := `
		SELECT from_id, id AS to_id, company_id
		FROM commodity_market.exchangers AS e1
		LEFT JOIN (
			SELECT id AS from_id, tag
			FROM commodity_market.exchangers 
		) AS e2
		ON e2.tag = ?
		LEFT JOIN (
			SELECT id AS company_id, tag
			FROM commodity_market.shipment_companies 
		) AS sh
		ON sh.tag = ?
		WHERE e1.tag = ?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		race.FromExch.Tag,
		race.Company.Tag,
		race.ToExch.Tag,
	).Scan(&race.FromExch.Id, &race.ToExch.Id, &race.Company.Id)

	if err != nil {
		panic(err)
	}

	sqlStatement = `
		UPDATE commodity_market.races
		SET from_id = ?, to_id = ?, race_date = ?, company_id = ?
		WHERE id = ?;
	`

	err = db.sql.QueryRow(
		sqlStatement,
		race.FromExch.Id,
		race.ToExch.Id,
		race.DateStamp,
		race.Company.Id,
		id,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) DeleteRace(id int) {
	sqlStatement := `
		DELETE FROM commodity_market.races 
		WHERE id=?;
	`
	err := db.sql.QueryRow(
		sqlStatement,
		id,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) GetRaceById(id int) *h.Race {
	sqlStatement := `
		SELECT R.id AS race_id, E1.tag AS from_tag, E2.tag AS to_tag, R.status
		FROM commodity_market.races R 
		JOIN commodity_market.exchangers E1
			ON R.from_id = E1.id
		JOIN commodity_market.exchangers E2
			ON R.to_id = E2.id
		WHERE R.id = ?;
	`

	race := &h.Race{
		FromExch: &h.Exchanger{},
		ToExch:   &h.Exchanger{},
	}

	err := db.sql.QueryRow(
		sqlStatement,
		id,
	).Scan(
		&race.Id,
		&race.FromExch.Tag,
		&race.ToExch.Tag,
		&race.Status,
	)

	if err != nil {
		panic(err)
	}

	return race
}

func (db *Database) FinishRace(database string, raceId int) {
	if !h.ValidDatabase(database) {
		panic(fmt.Errorf("🛠 Invalid database name"))
	}

	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStatement := fmt.Sprintf(`
		SELECT user_id, commodity_id, volume
		FROM %s.expected_cargo
		WHERE race_id = ?;
	`, database)

	rows, err := tx.Query(
		sqlStatement,
		raceId,
	)

	if err != nil {
		panic(err)
	}

	var commodities [](*h.Commodity)

	for rows.Next() {
		commodity := h.Commodity{
			Owner: &h.User{},
		}
		if err := rows.Scan(
			&commodity.Owner.Id,
			&commodity.Id,
			&commodity.Volume,
		); err != nil {
			tx.Rollback()
			panic(err)
		}

		commodities = append(commodities, &commodity)
	}
	rows.Close()

	for _, commodity := range commodities {
		upsertTransactionCommodities(tx, database, commodity.Owner.Id, commodity.Id, commodity.Volume)
		upsertTransactionCargo(tx, database, raceId, commodity.Id, commodity.Owner.Id, -commodity.Volume)
	}

	sqlStatement = `
		UPDATE commodity_market.races
		SET status = IF(status = 'preparing', 'arrive', status)
		WHERE id = ?
	`

	err = tx.QueryRow(sqlStatement, raceId).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()

	if err != nil {
		panic(err)
	}
}
