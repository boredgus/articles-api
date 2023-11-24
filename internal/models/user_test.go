package models

import (
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
	type fields struct {
		repo repo.UserRepository
		pswd auth.Password
	}
	type args struct {
		user domain.User
	}
	type mockedRes struct {
		createErr error
		hashErr   error
	}
	setup := func(mocks *fields, res mockedRes) func() {
		repoCall := mocks.repo.(*repoMocks.UserRepository).EXPECT().
			Create(mock.Anything).Return(res.createErr).Once()
		pswdCall := mocks.pswd.(*authMocks.Password).EXPECT().
			Hash(mock.Anything).Return("", res.hashErr).Once()

		return func() {
			repoCall.Unset()
			pswdCall.Unset()
		}
	}
	valueOfFields := fields{
		repo: repoMocks.NewUserRepository(t),
		pswd: authMocks.NewPassword(t),
	}
	validUser := domain.NewUser("username", "PASsword/123")
	hashErr := fmt.Errorf("hash error")
	tests := []struct {
		name      string
		fields    fields
		mockedRes mockedRes
		args      args
		wantErr   error
	}{
		{
			name:      "invalid credentials",
			fields:    valueOfFields,
			mockedRes: mockedRes{},
			args:      args{user: domain.NewUser("qw", "er")},
			wantErr:   InvalidAuthParameterErr,
		},
		{
			name:      "password hashing failed",
			fields:    valueOfFields,
			mockedRes: mockedRes{hashErr: hashErr},
			args:      args{user: validUser},
			wantErr:   hashErr,
		},
		{
			name:      "username is duplicated",
			fields:    valueOfFields,
			mockedRes: mockedRes{createErr: UsernameDuplicationErr},
			args:      args{user: validUser},
			wantErr:   UsernameDuplicationErr,
		},
		{
			name:      "success",
			fields:    valueOfFields,
			mockedRes: mockedRes{},
			args:      args{user: validUser},
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.fields, tt.mockedRes)
			defer cleanSetup()
			err := user{
				repo: tt.fields.repo,
				pswd: tt.fields.pswd,
			}.Create(tt.args.user)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestUserService_Authorize(t *testing.T) {
	type fields struct {
		repo  repo.UserRepository
		token auth.Token
		pswd  auth.Password
	}
	type args struct {
		user domain.User
	}
	type mockedRes struct {
		user        repo.User
		repoErr     error
		isPswdValid bool
		token       string
		tokenErr    error
	}
	setup := func(mocks *fields, res mockedRes) func() {
		repoCall := mocks.repo.(*repoMocks.UserRepository).EXPECT().
			Get(mock.Anything).Return(res.user, res.repoErr).Once()
		pswdCall := mocks.pswd.(*authMocks.Password).EXPECT().
			Compare(mock.Anything, mock.Anything).Return(res.isPswdValid).Once()
		tokenCall := mocks.token.(*authMocks.Token).EXPECT().
			Generate(mock.Anything).Return(res.token, res.tokenErr).Once()
		return func() {
			repoCall.Unset()
			pswdCall.Unset()
			tokenCall.Unset()
		}
	}
	valueOfFields := fields{
		repo:  repoMocks.NewUserRepository(t),
		pswd:  authMocks.NewPassword(t),
		token: authMocks.NewToken(t),
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
		fields     fields
		mockedRes  mockedRes
		args       args
		wantUserId string
		wantToken  string
		wantErr    error
	}{
		{
			name:       "no user with such username",
			fields:     valueOfFields,
			mockedRes:  mockedRes{user: repo.User{}, repoErr: InvalidAuthParameterErr},
			args:       args{user: validUser},
			wantUserId: "",
			wantToken:  "",
			wantErr:    InvalidAuthParameterErr,
		},
		{
			name:       "invalid password",
			fields:     valueOfFields,
			mockedRes:  mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: false},
			args:       args{user: validUser},
			wantUserId: "",
			wantToken:  "",
			wantErr:    InvalidAuthParameterErr,
		},
		{
			name:       "failed to generate token",
			fields:     valueOfFields,
			mockedRes:  mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: true, token: "", tokenErr: tokenErr},
			args:       args{user: validUser},
			wantUserId: userFromRepo.OId,
			wantToken:  "",
			wantErr:    tokenErr,
		},
		{
			name:       "success",
			fields:     valueOfFields,
			mockedRes:  mockedRes{user: userFromRepo, repoErr: nil, isPswdValid: true, token: userToken, tokenErr: nil},
			args:       args{user: validUser},
			wantUserId: userFromRepo.OId,
			wantToken:  userToken,
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(&tt.fields, tt.mockedRes)
			defer cleanSetup()
			u := user{
				repo:  tt.fields.repo,
				token: tt.fields.token,
				pswd:  tt.fields.pswd,
			}
			gotUserId, gotToken, err := u.Authorize(tt.args.user)
			assert.Equal(t, gotUserId, tt.wantUserId)
			assert.Equal(t, gotToken, tt.wantToken)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
