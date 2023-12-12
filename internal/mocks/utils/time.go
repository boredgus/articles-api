package mocks

import "time"

func MockTimeNow(fakeValue time.Time) func() time.Time {
	return func() time.Time {
		return fakeValue
	}
}
