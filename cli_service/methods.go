package cli_service

import (
	"bufio"
	"database-course-work/auth_service"
	ex_service "database-course-work/exchanger_service"
	h "database-course-work/helpers"
	"database-course-work/sql_service"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var user h.CliUser

func readFromConsole(title string, validator string) string {
	var value string
	for {
		fmt.Print(title)
		reader := bufio.NewReader(os.Stdin)
		value, _ = reader.ReadString('\n')
		value = strings.TrimSuffix(value, "\n")

		matched, _ := regexp.MatchString(validator, value)
		if len(validator) == 0 || matched {
			break
		}
		fmt.Println("Try again")
	}

	return value
}

func Launch(db *sql_service.Database) {
	value := ""
	for value != "0" {
		fmt.Print("-- MENU --\nSelect action:\n1) Sign in User\n2) Sign in Company\n3) Sign up User\n4) Sign up Company\n0) Exit\nChoice: ")
		fmt.Scanln(&value)

		if value == "1" {
			signInUser(db)
		} else if value == "2" {
			signInCompany(db)
		} else if value == "3" {
			signUpUser(db)
		} else if value == "4" {
			signUpCompany(db)
		} else if value == "0" {
			return
		} else {
			fmt.Println("Invalid input, try again or enter 0 to exit")
		}

		if len(user.Jwt) != 0 && user.UserType == "user" {
			isExit := launchUserMenu(db)
			if isExit {
				return
			}
		}

		if len(user.Jwt) != 0 && user.UserType == "company" {
			isExit := launchCompanyMenu(db)
			if isExit {
				return
			}
		}
	}
}

func signInUser(db *sql_service.Database) {
	fmt.Println("-- USER SIGN IN --")
	email := readFromConsole("Enter email: ", "^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]+$")
	password := readFromConsole("Enter password: ", ".")

	jwt, err := auth_service.SignIn(
		db,
		&h.User{
			Email:    email,
			Password: password,
		},
	)

	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
	}

	if len(jwt) != 0 {
		user.Jwt = jwt
		user.SelfUser, err = auth_service.GetUser(db, jwt)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		user.UserType = "user"
	}
}

func signInCompany(db *sql_service.Database) {
	fmt.Println("-- COMPANY SIGN IN --")
	tag := readFromConsole("Enter tag: ", ".")
	password := readFromConsole("Enter password: ", ".")

	fmt.Printf("%s %s", tag, password)

	jwt, err := auth_service.SignInCompany(
		db,
		&h.Company{
			Tag:      tag,
			Password: password,
		},
	)

	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
	}

	if len(jwt) != 0 {
		user.Jwt = jwt
		user.SelfCompany, err = auth_service.GetCompany(db, jwt)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		user.UserType = "company"
	}
}

func signUpUser(db *sql_service.Database) {
	fmt.Println("-- USER SIGN UP --")

	name := readFromConsole("Enter user name: ", ".")
	surname := readFromConsole("Enter user surname: ", ".")
	email := readFromConsole("Enter user email: ", "^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]+$")
	bankAccount := readFromConsole("Enter bank account number (16 digits): ", "^[0-9]{16}$")
	password := readFromConsole("Enter user password: ", "^[a-zA-Z0-9]{8,}$")
	license := readFromConsole("Enter license or just skip it: ", "")

	jwt, err := auth_service.SignUp(
		db,
		&h.User{
			Name:        name,
			Surname:     surname,
			Email:       email,
			BankAccount: bankAccount,
			Password:    password,
			License: &h.License{
				Code: license,
			},
		},
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}

	if len(jwt) != 0 {
		user.Jwt = jwt
		user.SelfUser, err = auth_service.GetUser(db, jwt)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		user.UserType = "user"
	}
}

func signUpCompany(db *sql_service.Database) {
	fmt.Println("-- COMPANY SIGN UP --")

	tag := readFromConsole("Enter company tag: ", ".")
	title := readFromConsole("Enter comapny title: ", ".")
	email := readFromConsole("Enter company contact email (or skip): ", "^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]+|$")
	phoneNumber := readFromConsole("Enter company contact phone number (+1 (111) 111-1111) or skip: ", "^\\+[0-9]{1,4} \\([0-9]{3}\\) [0-9]{3}-[0-9]{4}|$")
	password := readFromConsole("Enter user password: ", "^[a-zA-Z0-9]{8,}$")

	jwt, err := auth_service.SignUpCompany(
		db,
		&h.Company{
			Tag:         tag,
			Title:       title,
			Email:       email,
			PhoneNumber: phoneNumber,
			Password:    password,
		},
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}

	if len(jwt) != 0 {
		user.Jwt = jwt
		user.SelfCompany, err = auth_service.GetCompany(db, jwt)
		if err != nil {
			fmt.Printf("⛔️ %s\n", err.Error())
			return
		}
		user.UserType = "company"
	}
}

func launchUserMenu(db *sql_service.Database) bool {
	fmt.Printf("Welcome %s %s\n", user.SelfUser.Name, user.SelfUser.Surname)
	value := ""

	for value != "0" {
		fmt.Print("-- USER MENU --\nSelect action:\n1) Check my commodities\n2) List commodities types\n3) Add order\n4) List orders\n5) Cancel order\n6) Broker section\n7) Sign out\n0) Exit\nChoice: ")
		fmt.Scanln(&value)

		if value == "1" {
			checkCommodities(db)
		} else if value == "2" {
			listCommodities(db)
		} else if value == "3" {
			addOrder(db)
		} else if value == "4" {
			listOrders(db)
		} else if value == "5" {
			cancelOrder(db)
		} else if value == "6" {
			isExit := brokerSection(db)
			if isExit {
				return isExit
			}
		} else if value == "7" {
			user.Jwt = ""
			user.SelfUser = nil
			user.UserType = ""
			return false
		} else if value == "0" {
			return true
		} else {
			fmt.Println("Invalid input, try again or enter 0 to exit")
		}
	}

	return true
}

func launchCompanyMenu(db *sql_service.Database) bool {
	fmt.Printf("Company %s\n", user.SelfCompany.Title)
	value := ""

	for value != "0" {
		fmt.Print("-- COMPANY MENU --\nSelect action:\n1) Add commodity\n2) List commodities types\n3) Sign out\n0) Exit\nChoice: ")
		fmt.Scanln(&value)

		if value == "1" {
			addCommodity(db)
		} else if value == "2" {
			listCommodities(db)
		} else if value == "3" {
			user.Jwt = ""
			user.SelfCompany = nil
			user.UserType = ""
			return false
		} else if value == "0" {
			return true
		} else {
			fmt.Println("Invalid input, try again or enter 0 to exit")
		}
	}

	return true
}

func addCommodity(db *sql_service.Database) {
	label := readFromConsole("Enter label of resource to give: ", ".")
	volume, err := strconv.ParseFloat(readFromConsole("Enter volume of commodity: ", "[0-9.]+"), 64)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}

	email := readFromConsole("Enter email of recipient user: ", "^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]+$")

	err = ex_service.AddCommodity(
		db,
		&h.Commodity{
			Label:  label,
			Volume: volume,
			Owner: &h.User{
				Email: email,
			},
		},
		user.Jwt,
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
}

func checkCommodities(db *sql_service.Database) {
	fmt.Println("-- USER COMMODITIES LIST --")

	commodities, err := ex_service.CheckCommodities(
		db,
		user.Jwt,
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
	h.PrintCommodities(commodities)
}

func listCommodities(db *sql_service.Database) {
	fmt.Println("-- COMMODITIES TYPES --")

	commodities, err := db.GetAvailableCommodities()
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
	h.PrintCommoditiesList(commodities)
}

func addOrder(db *sql_service.Database) {
	fmt.Println("-- USER ADDING ORDER --")
	side := readFromConsole("Enter side (buy or sell): ", "^(buy|sell)$")
	commodityLabel := readFromConsole("Enter commodity label (you can check labels from user menu): ", ".")
	strVolume := readFromConsole("Enter volume of order: ", "[0-9.]+")
	brokerEmail := readFromConsole("Enter email of preferable broker (or skip): ", "^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\\.[a-zA-Z0-9]+|$")

	volume, err := strconv.ParseFloat(strVolume, 64)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
	err = ex_service.AddOrder(
		db,
		&h.Order{
			Side: side,
			Commodity: &h.Commodity{
				Label:  commodityLabel,
				Volume: volume,
			},
			PrefBroker: &h.User{
				Email: brokerEmail,
			},
		},
		user.Jwt,
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
}

func listOrders(db *sql_service.Database) {
	fmt.Println("-- LIST OF ORDERS --")
	orders, err := ex_service.ReadUserOrders(
		db,
		user.Jwt,
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
	h.PrintPersonalOrders(orders)
}

func cancelOrder(db *sql_service.Database) {
	fmt.Println("-- CANCEL OF ORDER --")
	orderIdStr := readFromConsole("Enter id of order: ", "[0-9]+")

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
	err = ex_service.CancelOrder(
		db,
		orderId,
		user.Jwt,
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
}

func brokerSection(db *sql_service.Database) bool {
	value := ""

	for value != "0" && !user.SelfUser.IsBroker {
		fmt.Print("-- BROKER MENU --\nSelect action:\n1) Add license\n2) Back\n0) Exit\nChoice: ")
		fmt.Scanln(&value)

		if value == "1" {
			addLicense(db)
		} else if value == "2" {
			return false
		} else if value == "0" {
			return true
		} else {
			fmt.Println("Invalid input, try again or enter 0 to exit")
		}
	}

	for value != "0" {
		fmt.Print("-- BROKER MENU --\nSelect action:\n1) List all orders\n2) Execute order\n3) Back\n0) Exit\nChoice: ")
		fmt.Scanln(&value)

		if value == "1" {
			listAllOrders(db)
		} else if value == "2" {
			executeOrder(db)
		} else if value == "3" {
			return false
		} else if value == "0" {
			return true
		} else {
			fmt.Println("Invalid input, try again or enter 0 to exit")
		}
	}
	return true
}

func addLicense(db *sql_service.Database) {
	fmt.Println("-- ADD OF LICENSE --")

	license := readFromConsole("Enter license code: ", ".")

	err := auth_service.AssignBroker(
		db,
		&h.User{
			Email: user.SelfUser.Email,
			License: &h.License{
				Code: license,
			},
		},
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}

	user.SelfUser, err = auth_service.GetUser(db, user.Jwt)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
}

func listAllOrders(db *sql_service.Database) {
	fmt.Println("-- LIST OF ALL ORDERS --")

	orders, err := ex_service.ReadAllOrders(
		db,
		user.Jwt,
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
	h.PrintNativeOrders(orders)
}

func executeOrder(db *sql_service.Database) {
	fmt.Println("-- EXECUTE OF ORDER --")

	firstOrderId, err := strconv.Atoi(readFromConsole("Enter first order id: ", "[0-9]+"))
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}

	secondOrderId, err := strconv.Atoi(readFromConsole("Enter second order id: ", "[0-9]+"))
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}

	volume, err := strconv.ParseFloat(readFromConsole("Enter volume of order execution: ", "[0-9.]+"), 32)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
	err = ex_service.ExecuteOrder(
		db,
		firstOrderId,
		secondOrderId,
		volume,
		user.Jwt,
	)
	if err != nil {
		fmt.Printf("⛔️ %s\n", err.Error())
		return
	}
}
