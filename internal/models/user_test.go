package models

import (
	"a-article/internal/auth"
	"a-article/internal/domain"
	authMocks "a-article/internal/mocks/auth"
	repoMocks "a-article/internal/mocks/repo"
	"a-article/internal/models/repo"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Create(t *testing.T) {
	type args struct {
		user domain.User
	}
	type mockedRes struct {
		createErr error
		hashErr   error
	}
	repoMock := repoMocks.NewUserRepository(t)
	pswdMock := authMocks.NewCryptor(t)
	setup := func(res mockedRes) func() {
		repoCall := repoMock.EXPECT().
			Create(mock.Anything).Return(res.createErr).Once()
		pswdCall := pswdMock.EXPECT().
			Encrypt(mock.Anything).Return("", res.hashErr).Once()

		return func() {
			repoCall.Unset()
			pswdCall.Unset()
		}
	}
	validUser := domain.User{Username: "username", Password: "PASsword/123"}
	hashErr := fmt.Errorf("hash error")
	tests := []struct {
		name      string
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
			name:      "password hashing failed",
			mockedRes: mockedRes{hashErr: hashErr},
			args:      args{user: validUser},
			wantErr:   hashErr,
		},
		{
			name:      "username is duplicated",
			mockedRes: mockedRes{createErr: UsernameDuplicationErr},
			args:      args{user: validUser},
			wantErr:   UsernameDuplicationErr,
		},
		{
			name:      "success",
			mockedRes: mockedRes{},
			args:      args{user: validUser},
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := (&user{repo: repoMock, crptr: pswdMock}).RequestSignup(tt.args.user)
			if err != nil {
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
	tokenErr := fmt.Errorf("token error")
	tests := []struct {
		name       string
		mockedRes  mockedRes
		wantUserId string
		wantToken  string
		wantErr    error
	}{
		{
			name:       "no user with such username",
			mockedRes:  mockedRes{user: repo.User{}, repoErr: InvalidAuthParameterErr},
			wantUserId: "",
			wantToken:  "",
			wantErr:    InvalidAuthParameterErr,
		},
		{
			name:       "invalid password",
			mockedRes:  mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: false},
			wantUserId: "",
			wantToken:  "",
			wantErr:    InvalidAuthParameterErr,
		},
		{
			name:       "failed to generate token",
			mockedRes:  mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: true, token: "", tokenErr: tokenErr},
			wantUserId: userFromRepo.OId,
			wantToken:  "",
			wantErr:    tokenErr,
		},
		{
			name:       "success",
			mockedRes:  mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: true, token: userToken, tokenErr: nil},
			wantUserId: userFromRepo.OId,
			wantToken:  userToken,
			wantErr:    nil,
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
