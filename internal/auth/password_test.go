package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var hashTestCases = map[string]struct {
	err        bool
	hashLength int
}{
	"":           {err: false, hashLength: 60},
	"pa":         {err: false, hashLength: 60},
	"p12p":       {err: false, hashLength: 60},
	"pO90lkjlkj": {err: false, hashLength: 60},
	".8Jjhkj8/":  {err: false, hashLength: 60},
	".8Jjhkj8hfkjhfjksdhfkjsdhfkjsdhfkjdshfkjdshfkjsdhfkjsdhkjfhsdkjfhsdkjfksjf/": {
		err: true, hashLength: 0},
}

func TestPasswordHashing(t *testing.T) {
	pswdSvc := NewPassword()
	for pswd, expected := range hashTestCases {
		result, err := pswdSvc.Hash(pswd)
		assert.Equal(t, expected.hashLength, len(result))
		if expected.err {
			assert.NotNil(t, err)
			continue
		}
		assert.Nil(t, err)
	}
}

var compareTestCases = map[string]struct {
	hash    string
	isEqual bool
}{
	"": {
		hash:    "$2a$10$tXFSjq1P0MSLBzizJ05HeulzDwbaGwrUYXMnsgiC9Uwbyba1G8hfe",
		isEqual: true,
	},
	"ps": {
		hash:    "$2a$10$JOvaQRSoJjEB5zugUS/b4Ow4MIKvb8mwx1pCF1eajghg1s1.6DULK",
		isEqual: true,
	},
	"ps12": {
		hash:    "$2a$10$9Ze92fxD.8kUPVB4jWacw.o9j9jlgn6h4AVHEpLeIa.EFuhmEkNn.",
		isEqual: true,
	},
	"po.*hj": {
		hash:    "$2a$10$NC9dnDCF56JLL9Z9MTBLr.d8hUBmk1sAQZwirwYsYaqEULV5HKc/m",
		isEqual: true,
	},
	"7879jhjkKJ": {
		hash:    "$2a$10$NNFG7QLAijyl5jdupLSWauKlz1/FGZidbhT9gnWK6pAAzbt8ri4xG",
		isEqual: true,
	},
	"123": {
		hash:    "$2a$10$NNFG7QLAijyl5jdupLSWauKlz1/FGZidbhT9gnWK6pAAzbt8ri4xG",
		isEqual: false,
	},
	"123_": {
		hash:    "",
		isEqual: false,
	},
}

func TestPasswordComparing(t *testing.T) {
	pswdSvc := NewPassword()
	for pswd, hash := range compareTestCases {
		assert.Equal(t, hash.isEqual, pswdSvc.Compare(hash.hash, pswd))
	}
}
