package repo

type UserRole int

const (
	DefaultUserRole UserRole = iota + 0
	ModeratorRole
	AdminRole
)

var UserRoles = map[UserRole]string{
	DefaultUserRole: "user",
	ModeratorRole:   "moderator",
	AdminRole:       "admin",
}
var RoleToValue = map[string]UserRole{
	"user":      DefaultUserRole,
	"moderator": ModeratorRole,
	"admin":     AdminRole,
}

type User struct {
	OId      string   `sql:"o_id"`
	Username string   `sql:"username"`
	Password string   `sql:"pswd"`
	Role     UserRole `sql:"role"`
}

type UserRepository interface {
	Create(user User) error
	Get(username string) (User, error)
	GetByOId(oid string) (User, error)
}
