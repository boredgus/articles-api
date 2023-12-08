package repo

type User struct {
	OId      string `sql:"o_id"`
	Username string `sql:"username"`
	Password string `sql:"pswd"`
}

type UserRepository interface {
	Create(user User) error
	Get(username string) (User, error)
	GetByOId(oid string) (User, error)
}
