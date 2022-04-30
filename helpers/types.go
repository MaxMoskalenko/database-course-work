package helpers

type User struct {
	Id           int
	Name         string
	Surname      string
	Email        string
	Password     string
	IsBroker     uint8
	License      string
	ExchangerTag string
}

type Company struct {
	Tag      string
	Title    string
	Password string
	Type     string
}

type Exchanger struct {
	DatabaseName string
	Name         string
	Tag          string
}

type Commodity struct {
	Label  string
	Volume int
	Unit   string
	Id     int
	Owner  *User
}

type Order struct {
	Id         int
	Owner      *User
	Side       string
	State      string
	Commodity  *Commodity
	PrefBroker *User
}
