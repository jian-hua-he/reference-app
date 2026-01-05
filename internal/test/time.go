package test

import "time"

func FakeNow() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}
