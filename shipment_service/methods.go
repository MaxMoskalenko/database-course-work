package shipment_service

import (
	"database-course-work/auth_service"
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
	"time"
)

func CreateRace(
	db *sql_service.Database,
	race *h.Race,
	companyJWT string,
) {
	var err error
	race.Company, err = auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		panic(err)
	}

	if race.Company.Type != "s" || !db.CheckIsRecordExist("commodity_market", "shipment_companies", "tag", race.Company.Tag) {
		panic(fmt.Errorf("⛔️ This shipment company doesn`t exist"))
	}

	if !db.CheckIsRecordExist("commodity_market", "exchangers", "tag", race.FromExch.Tag) ||
		!db.CheckIsRecordExist("commodity_market", "exchangers", "tag", race.FromExch.Tag) {
		panic(fmt.Errorf("⛔️ Exchangers with these tags do not exist"))
	}

	race.DateStamp, err = time.Parse("2006-01-02 15:04", race.DateValue)
	if err != nil {
		fmt.Println(err)
		return
	}

	db.InsertRace(race)
}

func ReadRaces(
	db *sql_service.Database,
) [](*h.Race) {
	return db.ReadRaces()
}

func UpdateRace(
	db *sql_service.Database,
	id int,
	race *h.Race,
	companyJWT string,
) {
	var err error
	race.Company, err = auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		panic(err)
	}

	if race.Company.Type != "s" || !db.CheckIsRecordExist("commodity_market", "shipment_companies", "tag", race.Company.Tag) {
		panic(fmt.Errorf("⛔️ This shipment company doesn`t exist"))
	}

	if db.GetCompanyTag(id) != race.Company.Tag {
		panic(fmt.Errorf("⛔️ Race is not owned by company"))
	}

	if !db.CheckIsRecordExist("commodity_market", "exchangers", "tag", race.FromExch.Tag) ||
		!db.CheckIsRecordExist("commodity_market", "exchangers", "tag", race.FromExch.Tag) {
		panic(fmt.Errorf("⛔️ Exchangers with these tags do not exist"))
	}

	race.DateStamp, err = time.Parse("2006-01-02 15:04", race.DateValue)
	if err != nil {
		fmt.Println(err)
		return
	}

	db.UpdateRace(id, race)
}

func DeleteRace(
	db *sql_service.Database,
	id int,
	companyJWT string,
) {
	company, err := auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		panic(err)
	}

	if company.Type != "s" || !db.CheckIsRecordExist("commodity_market", "shipment_companies", "tag", company.Tag) {
		panic(fmt.Errorf("⛔️ This shipment company doesn`t exist"))
	}

	if db.GetCompanyTag(id) != company.Tag {
		panic(fmt.Errorf("⛔️ Race is not owned by company"))
	}

	db.DeleteRace(id)
}

func FinishRace(
	db *sql_service.Database,
	raceId int,
	companyJWT string,
) {
	company, err := auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		panic(err)
	}

	if company.Type != "s" || !db.CheckIsRecordExist("commodity_market", "shipment_companies", "tag", company.Tag) {
		panic(fmt.Errorf("⛔️ This shipment company doesn`t exist"))
	}

	if db.GetCompanyTag(raceId) != company.Tag {
		panic(fmt.Errorf("⛔️ Race is not owned by company"))
	}

	race := db.GetRaceById(raceId)
	database := db.GetDatabaseByTag(race.ToExch.Tag)

	if len(database) == 0 {
		panic(fmt.Errorf("⛔️ Database with a %s tag does not exist", race.ToExch.Tag))
	}

	if race.Status == "arrive" {
		panic(fmt.Errorf("⛔️ This race has already arrived"))
	}

	db.FinishRace(database, raceId)
}
