package login

type LoginSvc interface {
	Register(username, password string) error
	Authorize(username, password string) string
}
