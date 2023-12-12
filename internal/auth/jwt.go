package auth

import (
	"errors"
	"fmt"
	"time"
	"user-management/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	Username string `json:"username"`
	UserOId  string `json:"user_oid"`
	Role     string `json:"role"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	JWTPayload
}

type JWT struct {
	secretKey []byte
	now       func() time.Time
}

var JWTSecretKey = []byte(config.GetConfig().JWTSecretKey)

func NewJWT() Token[JWTPayload] {
	return JWT{
		secretKey: JWTSecretKey,
		now:       time.Now,
	}
}

func (t JWT) Generate(payload JWTPayload) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256,
		JWTClaims{
			JWTPayload: JWTPayload{
				Username: payload.Username,
				UserOId:  payload.UserOId,
				Role:     payload.Role,
			},
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(t.now().UTC().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(t.now().UTC()),
			},
		},
	).SignedString(t.secretKey)
}

var UnexpectedSigningMethod = errors.New("unexpected signing method")

func (t JWT) Decode(token string) (JWTPayload, error) {
	var customClaims JWTClaims
	_, err := jwt.ParseWithClaims(token, &customClaims, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", UnexpectedSigningMethod, tkn.Header["alg"])
		}
		return t.secretKey, nil
	})
	if err != nil {
		return JWTPayload{}, err
	}
	return JWTPayload{
		Username: customClaims.Username,
		UserOId:  customClaims.UserOId,
		Role:     customClaims.Role,
	}, nil
}
