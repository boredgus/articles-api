package auth

import (
	"time"

	"github.com/sethvargo/go-password/password"
)

var PasscodeExpiresAfter = 50 * time.Minute

func GeneratePasscode() (string, error) {
	return password.Generate(6, 6, 0, false, false)
}
