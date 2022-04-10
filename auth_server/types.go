package auth_server

type User struct {
	id        int
	user_type string
}

type ExchangeUser struct {
	User
	name    string
	surname string
	email   string
}

type CompanyUser struct {
	User
	title string
}

type ShipmentUser struct {
	User
	title string
}

type Broker struct {
	ExchangeUser
	license string
}
