package auth

import (
	"testing"
	"user-management/internal/domain"

	"github.com/stretchr/testify/assert"
)

var tokenGeneratingTestCases = map[domain.User]struct {
	token string
	err   error
}{
	{Username: "user", Password: "pass"}: {
		token: "dXNlcjpwYXNz",
		err:   nil,
	},
	{Username: "usew2r", Password: "pass:pass"}: {
		token: "dXNldzJyOnBhc3M6cGFzcw==",
		err:   nil,
	},
	{Username: "uwse0r", Password: "pass123._"}: {
		token: "dXdzZTByOnBhc3MxMjMuXw==",
		err:   nil,
	},
	{Username: "", Password: ""}: {
		token: "Og==",
		err:   nil,
	},
}

func TestTokenGenerating(t *testing.T) {
	token := NewToken()
	for user, expected := range tokenGeneratingTestCases {
		result, err := token.Generate(user)
		assert.Equal(t, expected.token, result, user)
		assert.Equal(t, expected.err, err)
	}
}

var tokenDecodingTestCases = map[string]struct {
	user domain.User
	err  error
}{
	"dXNlcjpwYXNz": {
		user: domain.NewUser("user", "pass"),
		err:  nil,
	},
	"kdfjkjdflks": {
		user: domain.User{},
		err:  InvalidToken,
	},
	"dXNlcnBhc3M=": {
		user: domain.User{},
		err:  InvalidToken,
	},
}

func TestTokenDecoding(t *testing.T) {
	tokenSvc := NewToken()
	for token, expected := range tokenDecodingTestCases {
		user, err := tokenSvc.Decode(token)
		assert.Equal(t, expected.user, user)
		if expected.err != nil {
			assert.ErrorIs(t, err, expected.err)
			continue
		}
		assert.Nil(t, err)
	}
}
