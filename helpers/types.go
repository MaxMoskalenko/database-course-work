package helpers

type User struct {
	Name     string `sql:"name"`
	Surname  string `sql:"surname"`
	Email    string `sql:"email"`
	Password string `sql:"password"`
	IsBroker uint8  `sql:"is_broker"`
	License  string `sql:"license"`
}

type Company struct {
	Tag      string `sql:"tag"`
	Title    string `sql:"title"`
	Password string `sql:"password"`
	Type     string
}

type Exchanger struct {
	DatabaseName string `sql:"database_name"`
	Name         string `sql:"name"`
	Tag          string `sql:"tag"`
}
