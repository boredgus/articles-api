package domain

type User struct {
	Username string `json:"username" sql:"username" form:"username"`
	Password string `json:"password" sql:"pswd" form:"password"`
}

func NewUser(name, pswd string) User {
	return User{Username: name, Password: pswd}
}
