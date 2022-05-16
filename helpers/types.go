package helpers

import "time"

type User struct {
	Id          int
	Name        string
	Surname     string
	Email       string
	BankAccount string
	Password    string
	IsBroker    bool
	License     *License
}

type License struct {
	Id      int
	Code    string
	IsTaken bool
}

type Company struct {
	Id          int
	Tag         string
	Title       string
	Email       string
	PhoneNumber string
	Password    string
}

type Exchanger struct {
	Id           int
	DatabaseName string
	Name         string
	Tag          string
}

type Commodity struct {
	Label  string
	Volume float64
	Unit   string
	Id     int
	Owner  *User
	Source *CommoditySource
}

type CommoditySource struct {
	Type          string
	CompanyId     int
	SourceUserId  int
	SourceOrderId int
	DestOrderId   int
	BrokerId      int
}

type Order struct {
	Id             int
	Owner          *User
	Side           string
	State          string
	Commodity      *Commodity
	ExecutedVolume float64
	PrefBroker     *User
	Exchnager      *Exchanger
}

type Race struct {
	Id          int
	FromExch    *Exchanger
	ToExch      *Exchanger
	DateStamp   time.Time
	DateValue   string
	Commodities [](*Commodity)
	Company     *Company
	Status      string
}
