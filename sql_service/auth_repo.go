package sql_service

import (
	h "database-course-work/helpers"
)

func (db *Database) SignUp(user *h.User) {
	sqlStatement := `
		INSERT INTO users (name, surname, email, bank_account, password) 
		VALUES (?, ?, ?, ?, ?); 
	`

	err := db.sql.QueryRow(
		sqlStatement,
		user.Name,
		user.Surname,
		user.Email,
		user.BankAccount,
		user.Password,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) AssignBroker(userId int, licenseId int) {
	tx, err := db.sql.Begin()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStatement := `
		INSERT INTO brokers (user_id, license_id)
		VALUES (?, ?);
	`

	err = tx.QueryRow(
		sqlStatement,
		userId,
		licenseId,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStatement = `
		UPDATE licenses
		SET is_taken=true
		WHERE id=?;
	`

	err = tx.QueryRow(
		sqlStatement,
		licenseId,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStatement = `
		UPDATE users
		SET is_broker=true
		WHERE id=?
	`

	err = tx.QueryRow(
		sqlStatement,
		userId,
	).Err()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()

	if err != nil {
		panic(err)
	}

}

func (db *Database) SignUpCompany(company *h.Company) {
	sqlStatement := `
		INSERT INTO companies (title, tag, password, email, phone_number) 
		VALUES (?, ?, ?, ?, ?); 
	`

	err := db.sql.QueryRow(
		sqlStatement,
		company.Title,
		company.Tag,
		company.Password,
		company.Email,
		company.PhoneNumber,
	).Err()

	if err != nil {
		panic(err)
	}
}

func (db *Database) GetUserOnLogin(login *h.User) *h.User {
	var user h.User

	sqlStatement := `
		SELECT email 
		FROM users 
		WHERE email=? AND password=?;
	`

	db.sql.QueryRow(
		sqlStatement,
		login.Email,
		login.Password,
	).Scan(&user.Email)

	return &user
}

func (db *Database) GetCompanyOnLogin(login *h.Company) *h.Company {
	var company h.Company

	sqlStatement := `
		SELECT tag 
		FROM companies 
		WHERE tag=? AND password=?;
	`

	db.sql.QueryRow(
		sqlStatement,
		login.Tag,
		login.Password,
	).Scan(&company.Tag)

	return &company
}

func (db *Database) GetUserData(email string) *h.User {
	var user h.User

	sqlStatement := `
		SELECT id, email, name, surname, is_broker
		FROM users
		WHERE email=?;
	`

	db.sql.QueryRow(
		sqlStatement,
		email,
	).Scan(&user.Id, &user.Email, &user.Name, &user.Surname, &user.IsBroker)

	return &user
}
