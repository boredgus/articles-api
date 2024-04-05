package models

import (
	"a-article/internal/auth"
	"a-article/internal/domain"
	"a-article/internal/mailing"
	"a-article/internal/models/repo"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserModel interface {
	ConfirmSignup(username, passcode string) error
	RequestSignup(user domain.User) error
	Authorize(username, password string) (string, error)
	Delete(issuerRole, userToDeleteOId string) error
	UpdateRole(issuerRole, userToUpdateOId, roleToSet string) error
}

var InvalidAuthParameterErr = errors.New("username or password is invalid")
var UsernameDuplicationErr = errors.New("user with such username already exists")
var UserNotFoundErr = errors.New("user not found")
var ExpiredPasscodeErr = errors.New("passcode is expired")

func NewUserModel(repo repo.UserRepository) UserModel {
	return &user{
		repo:    repo,
		token:   auth.NewJWT(),
		crptr:   auth.NewCryptor(),
		mailman: mailing.NewMailman(),
	}
}

type user struct {
	repo    repo.UserRepository
	token   auth.Token[auth.JWTPayload]
	crptr   auth.Cryptor
	mailman mailing.Mailman
}

func (u *user) RequestSignup(user domain.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("%w: %w", InvalidDataErr, err)
	}
	_, err := u.repo.Get(user.Username)
	if err == nil {
		return UsernameDuplicationErr
	}
	if !errors.Is(err, NotFoundErr) {
		return err
	}
	encryptedPswd, err := u.crptr.Encrypt(user.Password)
	if err != nil {
		return err
	}
	passcode, err := auth.GeneratePasscode()
	if err != nil {
		return err
	}
	encryptedPasscode, err := u.crptr.Encrypt(passcode)
	if err != nil {
		return err
	}
	if err := u.repo.RegisterSignupRequest(repo.SignupRequest{
		Email:    user.Username,
		Password: encryptedPswd,
		Passcode: encryptedPasscode,
	}); err != nil {
		return err
	}
	u.mailman.ConfirmSignupEmail(user.Username, passcode)
	return nil
}

func (u *user) ConfirmSignup(email, passcode string) error {
	_, err := u.repo.Get(email)
	if err == nil {
		return UsernameDuplicationErr
	}
	if !errors.Is(err, NotFoundErr) {
		return err
	}
	reqData, err := u.repo.GetSignupRequest(email)
	if err != nil {
		return err
	}
	if reqData.AttemptedAt.Add(auth.PasscodeExpiresAfter).Before(time.Now().UTC()) {
		return ExpiredPasscodeErr
	}
	if !u.crptr.Compare(reqData.Passcode, passcode) {
		return fmt.Errorf("%w: passcode does not match", InvalidDataErr)
	}

	if err := u.repo.Create(repo.User{
		OId:      uuid.New().String(),
		Username: reqData.Email,
		Password: reqData.Password,
		Role:     domain.DefaultUserRole,
	}); err != nil {
		return err
	}
	u.mailman.WelcomeEmail(reqData.Email)
	return nil
}

func (u *user) Authorize(username, password string) (token string, err error) {
	userFromDB, err := u.repo.Get(username)
	if errors.Is(err, NotFoundErr) || !u.crptr.Compare(userFromDB.Password, password) {
		return "", InvalidAuthParameterErr
	}
	if err != nil {
		return "", err
	}
	return u.token.Generate(auth.JWTPayload{
		Username: userFromDB.Username,
		UserOId:  userFromDB.OId,
		Role:     string(userFromDB.Role),
	})
}

func (u *user) Delete(issuerRole, userToDeleteOId string) error {
	if issuerRole != string(domain.AdminRole) {
		return fmt.Errorf("%w: you have to be admin", NotEnoughRightsErr)
	}
	userFromDB, err := u.repo.GetByOId(userToDeleteOId)
	if err != nil {
		return err
	}
	if userFromDB.Role == domain.AdminRole {
		return fmt.Errorf("%w: you cannot delete admin", NotEnoughRightsErr)
	}
	return u.repo.Delete(userToDeleteOId)
}

func (u *user) UpdateRole(issuerRole, userToUpdateOId, roleToSet string) error {
	if issuerRole != string(domain.AdminRole) {
		return fmt.Errorf("%w: you have to be admin", NotEnoughRightsErr)
	}
	if !domain.UserRole(roleToSet).IsValid() {
		return fmt.Errorf("%w: unknown role '%v'", InvalidDataErr, roleToSet)
	}
	userFromDB, err := u.repo.GetByOId(userToUpdateOId)
	if err != nil {
		return err
	}
	if userFromDB.Role == domain.AdminRole {
		return fmt.Errorf("%w: you cannot update admin", NotEnoughRightsErr)
	}
	return u.repo.UpdateRole(userToUpdateOId, domain.UserRole(roleToSet))
}
