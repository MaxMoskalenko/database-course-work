package helpers

type User struct {
	Name     string `orm:"name"`
	Surname  string `orm:"surname"`
	Email    string `orm:"email"`
	Password string `orm:"password"`
	IsBroker uint8  `orm:"is_broker"`
	License  string `orm:"license"`
}
