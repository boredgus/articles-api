package models

import (
	"errors"
	"fmt"
	"testing"
	"user-management/internal/auth"
	"user-management/internal/domain"
	authMocks "user-management/internal/mocks/auth"
	repoMocks "user-management/internal/mocks/repo"
	"user-management/internal/models/repo"

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
	pswdMock := authMocks.NewPassword(t)
	setup := func(res mockedRes) func() {
		repoCall := repoMock.EXPECT().
			Create(mock.Anything).Return(res.createErr).Once()
		pswdCall := pswdMock.EXPECT().
			Hash(mock.Anything).Return("", res.hashErr).Once()

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
			wantErr:   InvalidUserErr,
		},
		{
			name:      "invalid api key on protected user creation",
			mockedRes: mockedRes{},
			args:      args{user: validUser},
			wantErr:   InvalidAPIKeyErr,
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
			err := user{repo: repoMock, pswd: pswdMock}.Create(tt.args.user)
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
	pswdMock := authMocks.NewPassword(t)
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
			gotToken, err := user{
				repo:  repoMock,
				token: tokenMock,
				pswd:  pswdMock,
			}.Authorize(validUser.Username, validUser.Password)
			assert.Equal(t, gotToken, tt.wantToken)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_user_Exists(t *testing.T) {
	type mockedRes struct {
		user        repo.User
		repoErr     error
		isPswdValid bool
	}
	repoMock := repoMocks.NewUserRepository(t)
	pswdMock := authMocks.NewPassword(t)
	userData := repo.User{OId: "o_id", Username: "username", Password: "pass"}
	setup := func(res mockedRes) func() {
		repoCall := repoMock.EXPECT().
			GetByOId(userData.OId).Return(res.user, res.repoErr).Once()
		pswdCall := pswdMock.EXPECT().
			Compare(res.user.Password, mock.Anything).Return(res.isPswdValid).Once()

		return func() {
			repoCall.Unset()
			pswdCall.Unset()
		}
	}
	someError := errors.New("some error")
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name:      "there is no user with such oid",
			mockedRes: mockedRes{repoErr: someError},
			wantErr:   someError,
		},
		{
			name:      "password is invalid",
			mockedRes: mockedRes{user: userData, isPswdValid: false},
			wantErr:   InvalidAuthParameterErr,
		},
		{
			name:      "success",
			mockedRes: mockedRes{user: userData, isPswdValid: true},
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := user{repo: repoMock, pswd: pswdMock}.Exists(userData.OId, userData.Password)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
