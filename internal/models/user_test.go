package models

import (
	"a-article/internal/auth"
	"a-article/internal/domain"
	authMocks "a-article/internal/mocks/auth"
	mailingMocks "a-article/internal/mocks/mailing"
	repoMocks "a-article/internal/mocks/repo"
	"a-article/internal/models/repo"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_RequestSignup(t *testing.T) {
	type args struct {
		user domain.User
	}
	type fields struct {
		generatePasscode func() (string, error)
	}
	type mockedRes struct {
		getErr     error
		requestErr error
		encryptErr error
	}
	repoMock := repoMocks.NewUserRepository(t)
	cryptorMock := authMocks.NewCryptor(t)
	mailingMock := mailingMocks.NewMailman(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			Get(mock.Anything).Return(repo.User{}, res.getErr).Once()
		pswdCall := cryptorMock.EXPECT().
			Encrypt(mock.Anything).Return("", res.encryptErr).Twice()
		signupRequestCall := repoMock.EXPECT().
			RegisterSignupRequest(mock.Anything).Return(res.requestErr) //.NotBefore(psscdCall) //.Maybe()
		mailMock := mailingMock.EXPECT().
			ConfirmSignupEmail(mock.Anything, mock.Anything) //.NotBefore(signupRequestCall) //.Maybe()
		return func() {
			getCall.Unset()
			pswdCall.Unset()
			signupRequestCall.Unset()
			mailMock.Unset()
		}
	}
	passcode := "123123"
	passcodeGenerator := func() (string, error) { return passcode, nil }
	validUser := domain.User{Username: "username@com.co", Password: "PASsword/123"}
	customErr := fmt.Errorf("custom err")

	tests := []struct {
		name      string
		fields    fields
		mockedRes mockedRes
		args      args
		wantErr   error
	}{
		{
			name:      "invalid user data",
			mockedRes: mockedRes{},
			args:      args{user: domain.NewUser("qw", "er")},
			wantErr:   InvalidDataErr,
		},
		{
			name:      "username is duplicated",
			mockedRes: mockedRes{},
			args:      args{user: validUser},
			wantErr:   UsernameDuplicationErr,
		},
		{
			name:      "failed to get user with such username",
			mockedRes: mockedRes{getErr: customErr},
			args:      args{user: validUser},
			wantErr:   customErr,
		},
		{
			name:      "password encrypting failed",
			mockedRes: mockedRes{getErr: NotFoundErr, encryptErr: customErr},
			args:      args{user: validUser},
			wantErr:   customErr,
		},
		{
			name:      "failed to generate passcode",
			fields:    fields{generatePasscode: func() (string, error) { return "", customErr }},
			mockedRes: mockedRes{getErr: NotFoundErr},
			args:      args{user: validUser},
			wantErr:   customErr,
		},
		{
			name:      "failed to save signup request",
			fields:    fields{generatePasscode: passcodeGenerator},
			mockedRes: mockedRes{getErr: NotFoundErr, requestErr: customErr},
			args:      args{user: validUser},
			wantErr:   customErr,
		},
		{
			name:      "success",
			fields:    fields{generatePasscode: passcodeGenerator},
			mockedRes: mockedRes{getErr: NotFoundErr},
			args:      args{user: validUser},
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := (&user{repo: repoMock, crptr: cryptorMock, generatePasscode: tt.fields.generatePasscode, mailman: mailingMock}).RequestSignup(tt.args.user)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestUserService_Authorize(t *testing.T) {
	type mockedRes struct {
		user        repo.User
		repoErr     error
		isPswdValid bool
		token       string
		tokenErr    error
	}
	repoMock := repoMocks.NewUserRepository(t)
	pswdMock := authMocks.NewCryptor(t)
	tokenMock := authMocks.NewToken[auth.JWTPayload](t)
	setup := func(res mockedRes) func() {
		repoCall := repoMock.EXPECT().
			Get(mock.Anything).Return(res.user, res.repoErr).Once()
		pswdCall := pswdMock.EXPECT().
			Compare(mock.Anything, mock.Anything).Return(res.isPswdValid).Once()
		tokenCall := tokenMock.EXPECT().
			Generate(mock.Anything).Return(res.token, res.tokenErr).Once()
		return func() {
			repoCall.Unset()
			pswdCall.Unset()
			tokenCall.Unset()
		}
	}
	validUser := domain.NewUser("username", "PASsword/123")
	userToken := "dXNlcm5hbWU6UEFTc3dvcmQvMTIz"
	userFromRepo := repo.User{
		OId:      "f7873e08-787b-45e9-b22f-82bdf505cca5",
		Username: "username",
		Password: "$2a$10$YCuxL/v4Rn07gP/ggFZIXeIxj6W9BhTaTj1CBDFH0Qysp4ZpI6Pw6",
	}
	customErr := fmt.Errorf("custom err")
	tokenErr := fmt.Errorf("token error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantToken string
		wantErr   error
	}{
		{
			name:      "no user with such username",
			mockedRes: mockedRes{user: repo.User{}, repoErr: InvalidAuthParameterErr},
			wantToken: "",
			wantErr:   InvalidAuthParameterErr,
		},
		{
			name:      "failed to get user",
			mockedRes: mockedRes{user: repo.User{}, repoErr: customErr},
			wantToken: "",
			wantErr:   customErr,
		},
		{
			name:      "invalid password",
			mockedRes: mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: false},
			wantToken: "",
			wantErr:   InvalidAuthParameterErr,
		},
		{
			name:      "failed to generate token",
			mockedRes: mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: true, token: "", tokenErr: tokenErr},
			wantToken: "",
			wantErr:   tokenErr,
		},
		{
			name:      "success",
			mockedRes: mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: true, token: userToken, tokenErr: nil},
			wantToken: userToken,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			gotToken, err := (&user{
				repo:  repoMock,
				token: tokenMock,
				crptr: pswdMock,
			}).Authorize(validUser.Username, validUser.Password)
			assert.Equal(t, gotToken, tt.wantToken)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	type args struct {
		issuerRole      string
		userToDeleteOId string
	}
	type mockedRes struct {
		userFromDB repo.User
		getErr     error
		deleteErr  error
	}
	repoMock := repoMocks.NewUserRepository(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			GetByOId(mock.Anything).Return(res.userFromDB, res.getErr).Maybe()
		deleteCall := repoMock.EXPECT().
			Delete(mock.Anything).Return(res.deleteErr).NotBefore(getCall).Maybe()
		return func() {
			getCall.Unset()
			deleteCall.Unset()
		}
	}
	someErr := errors.New("some err")
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:    "issuer is not an admin",
			wantErr: NotEnoughRightsErr,
		},
		{
			name:      "failed to get user from db",
			args:      args{issuerRole: "admin"},
			mockedRes: mockedRes{getErr: someErr},
			wantErr:   someErr,
		},
		{
			name:      "user to delete is an admin",
			args:      args{issuerRole: "admin"},
			mockedRes: mockedRes{userFromDB: repo.User{Role: domain.AdminRole}},
			wantErr:   NotEnoughRightsErr,
		},
		{
			name: "failed to delete user",
			args: args{issuerRole: "admin"},
			mockedRes: mockedRes{
				deleteErr:  someErr,
				userFromDB: repo.User{Role: domain.DefaultUserRole}},
			wantErr: someErr,
		},
		{
			name:      "success",
			args:      args{issuerRole: "admin"},
			mockedRes: mockedRes{userFromDB: repo.User{Role: domain.DefaultUserRole}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := (&user{repo: repoMock}).
				Delete(tt.args.issuerRole, tt.args.userToDeleteOId)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestUserService_UpdateRole(t *testing.T) {
	type args struct {
		issuerRole      string
		userToUpdateOId string
		roleToSet       string
	}
	type mockedRes struct {
		userFromDB repo.User
		getErr     error
		updateErr  error
	}
	repoMock := repoMocks.NewUserRepository(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			GetByOId(mock.Anything).Return(res.userFromDB, res.getErr).Maybe()
		deleteCall := repoMock.EXPECT().
			UpdateRole(mock.Anything, mock.Anything).Return(res.updateErr).NotBefore(getCall).Maybe()
		return func() {
			getCall.Unset()
			deleteCall.Unset()
		}
	}
	someErr := errors.New("some err")
	tests := []struct {
		name      string
		args      args
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:    "issuer is not an admin",
			wantErr: NotEnoughRightsErr,
		},
		{
			name:      "unknown role provided",
			args:      args{issuerRole: "admin"},
			mockedRes: mockedRes{userFromDB: repo.User{Role: domain.AdminRole}},
			wantErr:   InvalidDataErr,
		},
		{
			name:      "failed to get user from db",
			args:      args{issuerRole: "admin", roleToSet: "user"},
			mockedRes: mockedRes{getErr: someErr},
			wantErr:   someErr,
		},
		{
			name:      "user to update is an admin",
			args:      args{issuerRole: "admin", roleToSet: "user"},
			mockedRes: mockedRes{userFromDB: repo.User{Role: domain.AdminRole}},
			wantErr:   NotEnoughRightsErr,
		},
		{
			name: "failed to update user",
			args: args{issuerRole: "admin", roleToSet: "user"},
			mockedRes: mockedRes{
				updateErr:  someErr,
				userFromDB: repo.User{Role: domain.DefaultUserRole}},
			wantErr: someErr,
		},
		{
			name:      "success",
			args:      args{issuerRole: "admin", roleToSet: "user"},
			mockedRes: mockedRes{userFromDB: repo.User{Role: domain.DefaultUserRole}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := (&user{repo: repoMock}).
				UpdateRole(tt.args.issuerRole, tt.args.userToUpdateOId, tt.args.roleToSet)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_user_ConfirmSignup(t *testing.T) {
	type args struct {
		email    string
		passcode string
	}
	type mockedRes struct {
		getErr         error
		getRequestData repo.SignupRequest
		getRequestErr  error
		createErr      error
		compareRes     bool
	}
	repoMock := repoMocks.NewUserRepository(t)
	cryptorMock := authMocks.NewCryptor(t)
	mailingMock := mailingMocks.NewMailman(t)
	setup := func(res mockedRes) func() {
		getCall := repoMock.EXPECT().
			Get(mock.Anything).Return(repo.User{}, res.getErr).Once()
		requestCall := repoMock.EXPECT().
			GetSignupRequest(mock.Anything).Return(res.getRequestData, res.getRequestErr).NotBefore(getCall)
		compareCall := cryptorMock.EXPECT().
			Compare(mock.Anything, mock.Anything).Return(res.compareRes).NotBefore(requestCall)
		createCall := repoMock.EXPECT().
			Create(mock.Anything).Return(res.createErr).NotBefore(requestCall)
		mailMock := mailingMock.EXPECT().
			WelcomeEmail(mock.Anything)
		return func() {
			getCall.Unset()
			requestCall.Unset()
			compareCall.Unset()
			createCall.Unset()
			mailMock.Unset()
		}
	}
	passcode := "123123"
	customErr := fmt.Errorf("custom err")
	tests := []struct {
		name      string
		mockedRes mockedRes
		args      args
		wantErr   error
	}{
		{
			name:      "user with such email alredy exists",
			mockedRes: mockedRes{},
			wantErr:   UsernameDuplicationErr,
		},
		{
			name:      "failed to get user data",
			mockedRes: mockedRes{getErr: customErr},
			wantErr:   customErr,
		},
		{
			name:      "failed to get signup request data",
			mockedRes: mockedRes{getErr: NotFoundErr, getRequestErr: customErr},
			wantErr:   customErr,
		},
		{
			name: "passcode is expired",
			mockedRes: mockedRes{
				getErr:         NotFoundErr,
				getRequestData: repo.SignupRequest{AttemptedAt: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)},
			},
			wantErr: ExpiredPasscodeErr,
		},
		{
			name: "passcode mismatches",
			args: args{passcode: passcode},
			mockedRes: mockedRes{
				getErr:         NotFoundErr,
				getRequestData: repo.SignupRequest{AttemptedAt: time.Now().UTC()},
			},
			wantErr: InvalidDataErr,
		},
		{
			name: "failed to create user",
			args: args{passcode: passcode},
			mockedRes: mockedRes{
				getErr:         NotFoundErr,
				getRequestData: repo.SignupRequest{AttemptedAt: time.Now().UTC(), Passcode: passcode},
				compareRes:     true,
				createErr:      customErr,
			},
			wantErr: customErr,
		},
		{
			name: "success",
			args: args{passcode: passcode},
			mockedRes: mockedRes{
				getErr:         NotFoundErr,
				getRequestData: repo.SignupRequest{AttemptedAt: time.Now().UTC(), Passcode: passcode},
				compareRes:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := (&user{
				repo:    repoMock,
				crptr:   cryptorMock,
				mailman: mailingMock,
			}).ConfirmSignup(tt.args.email, tt.args.passcode)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
