package domain

type User struct {
	Username string `json:"username"`
	Password string `json:"pswd"`
}

func NewUser(name, pswd string) User {
	return User{Username: name, Password: pswd}
}
