package auth_service

import (
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
)

func SignUp(
	db *sql_service.Database,
	user *h.User,
) (string, error) {
	if !h.ValidPassword(user.Password) {
		return "", fmt.Errorf("your password is incorrect")
	}

	if !h.ValidEmail(user.Email) {
		return "", fmt.Errorf("your email is incorrect")
	}

	if !h.ValidBankAccount(user.BankAccount) {
		return "", fmt.Errorf("your email is incorrect")
	}

	if db.CheckIsRecordExist("users", "email", user.Email) {
		return "", fmt.Errorf("your email is already in use")
	}

	user.Password = h.Hash(user.Password)

	err := db.SignUp(user)

	if err != nil {
		return "", err
	}

	var noBrokerError error
	if len(user.License.Code) != 0 {
		noBrokerError = AssignBroker(db, user)
	}

	jwt, err := generateJWT(user)

	if err != nil {
		return "", err
	}

	return jwt, noBrokerError
}

func AssignBroker(
	db *sql_service.Database,
	user *h.User,
) error {
	if !db.CheckIsRecordExist("licenses", "license_code", user.License.Code) {
		return fmt.Errorf("invalid license code")
	}

	var err error

	user.Id, err = db.GetId("users", "email", user.Email)
	if err != nil {
		return err
	}

	license, err := db.GetLicense(user.License.Code)
	if err != nil {
		return err
	}

	if license.IsTaken {
		return fmt.Errorf("this license code is taken")
	}

	err = db.AssignBroker(user.Id, license.Id)
	if err != nil {
		return err
	}

	return nil
}

func SignUpCompany(
	db *sql_service.Database,
	company *h.Company,
) (string, error) {
	if !h.ValidPassword(company.Password) {
		return "", fmt.Errorf("your password is incorrect")
	}

	if db.CheckIsRecordExist("companies", "tag", company.Tag) {
		return "", fmt.Errorf("this tag is already in use")
	}

	if len(company.Email) != 0 && !h.ValidEmail(company.Email) {
		return "", fmt.Errorf("company email is invalid")
	}

	if len(company.PhoneNumber) != 0 && !h.ValidPhone(company.PhoneNumber) {
		return "", fmt.Errorf("company phone number is invalid")
	}

	company.Password = h.Hash(company.Password)

	err := db.SignUpCompany(company)

	if err != nil {
		return "", err
	}

	jwt, err := generateCompanyJWT(company)

	if err != nil {
		return "", err
	}

	return jwt, nil
}

func SignIn(
	db *sql_service.Database,
	login *h.User,
) (string, error) {
	login.Password = h.Hash(login.Password)

	user, err := db.GetUserOnLogin(login)

	if err != nil {
		return "", err
	}

	if user.Email == "" {
		return "", fmt.Errorf("wrong credentials")
	}

	jwt, err := generateJWT(user)

	if err != nil {
		return "", err
	}

	return jwt, nil
}

func SignInCompany(
	db *sql_service.Database,
	login *h.Company,
) (string, error) {
	login.Password = h.Hash(login.Password)

	company, err := db.GetCompanyOnLogin(login)

	if err != nil {
		return "", err
	}

	if company.Tag == "" {
		return "", fmt.Errorf("wrong credentials")
	}

	jwt, err := generateCompanyJWT(company)

	if err != nil {
		return "", err
	}

	return jwt, nil
}

func GetUser(
	db *sql_service.Database,
	jwt string,
) (*h.User, error) {
	user, err := ReadJWT(jwt)
	if err != nil {
		return nil, err
	}

	userData, err := db.GetUserData(user.Email)
	if err != nil {
		return nil, err
	}

	if userData == nil {
		return nil, fmt.Errorf("user does not exist")
	}

	user.Id = userData.Id
	user.Name = userData.Name
	user.Surname = userData.Surname
	user.IsBroker = userData.IsBroker

	return user, nil
}

func GetCompany(
	db *sql_service.Database,
	jwt string,
) (*h.Company, error) {
	company, err := ReadCompanyJWT(jwt)
	if err != nil {
		return nil, err
	}

	companyData, err := db.GetCompanyData(company.Tag)
	if err != nil {
		return nil, err
	}

	if companyData == nil {
		return nil, fmt.Errorf("user does not exist")
	}

	company.Id = companyData.Id
	company.Title = companyData.Title

	return company, nil
}
