package helpers

type User struct {
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
}
