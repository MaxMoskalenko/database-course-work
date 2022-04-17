package exchanger_service

import (
	"database-course-work/auth_service"
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
)

func AddCommodity(
	db *sql_service.Database,
	user *h.User,
	commodity *h.Commodity,
	companyJWT string,
) {
	company, err := auth_service.ReadCompanyJWT(companyJWT)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	database := db.GetDatabaseByTag(user.ExchangerTag)

	if len(database) == 0 {
		fmt.Printf("⛔️ Database with a %s tag does not exist\n", user.ExchangerTag)
		return
	}

	if !db.CheckIsRecordExist("commodity_market", h.GetTableFromType(company.Type), "tag", company.Tag) {
		fmt.Printf("⛔️ No such company %s\n", company.Tag)
		return
	}

	if !db.CheckIsRecordExist("commodity_market", "commodity_types", "label", commodity.Label) {
		fmt.Printf("⛔️ No such commodity type %s\n", commodity.Label)
		return
	}

	if !db.CheckIsRecordExist(database, "users", "email", user.Email) {
		fmt.Printf("⛔️ No such user %s \n", user.Email)
		return
	}

	db.AddCommodity(database, user.Email, commodity)

}
