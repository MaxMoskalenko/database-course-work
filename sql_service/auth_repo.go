package sql_service

import (
	h "database-course-work/helpers"
)

func (db *Database) SignUp(user *h.User) error {
	sqlStatement := `
		INSERT INTO commodity_market.users (name, surname, email, bank_account, password) 
		VALUES (?, ?, ?, ?, ?); 
	`

	return db.sql.QueryRow(
		sqlStatement,
		user.Name,
		user.Surname,
		user.Email,
		user.BankAccount,
		user.Password,
	).Err()
}

func (db *Database) AssignBroker(userId int, licenseId int) error {
	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStatement := `
		INSERT INTO commodity_market.brokers (user_id, license_id)
		VALUES (?, ?);
	`

	err = tx.QueryRow(
		sqlStatement,
		userId,
		licenseId,
	).Err()

	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStatement = `
		UPDATE commodity_market.licenses
		SET is_taken=true
		WHERE id=?;
	`

	err = tx.QueryRow(
		sqlStatement,
		licenseId,
	).Err()

	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStatement = `
		UPDATE commodity_market.users
		SET is_broker=true
		WHERE id=?
	`

	err = tx.QueryRow(
		sqlStatement,
		userId,
	).Err()

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *Database) SignUpCompany(company *h.Company) error {
	sqlStatement := `
		INSERT INTO commodity_market.companies (title, tag, password, email, phone_number) 
		VALUES (?, ?, ?, ?, ?); 
	`

	return db.sql.QueryRow(
		sqlStatement,
		company.Title,
		company.Tag,
		company.Password,
		company.Email,
		company.PhoneNumber,
	).Err()
}

func (db *Database) GetUserOnLogin(login *h.User) (*h.User, error) {
	var user h.User

	sqlStatement := `
		SELECT email 
		FROM commodity_market.users 
		WHERE email=? AND password=?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		login.Email,
		login.Password,
	).Scan(&user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetCompanyOnLogin(login *h.Company) (*h.Company, error) {
	var company h.Company

	sqlStatement := `
		SELECT tag 
		FROM commodity_market.companies 
		WHERE tag=? AND password=?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		login.Tag,
		login.Password,
	).Scan(&company.Tag)

	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (db *Database) GetUserData(email string) (*h.User, error) {
	var user h.User

	sqlStatement := `
		SELECT id, email, name, surname, is_broker
		FROM commodity_market.users
		WHERE email=?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		email,
	).Scan(&user.Id, &user.Email, &user.Name, &user.Surname, &user.IsBroker)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetCompanyData(tag string) (*h.Company, error) {
	var company h.Company

	sqlStatement := `
		SELECT id, title
		FROM commodity_market.companies
		WHERE tag=?;
	`

	err := db.sql.QueryRow(
		sqlStatement,
		tag,
	).Scan(&company.Id, &company.Title)

	if err != nil {
		return nil, err
	}

	return &company, nil
}
