package common

import "time"

type Clock interface {
	Now() time.Time
}

type StandardClock struct{}

func (StandardClock) Now() time.Time {
	return time.Now().UTC()
}
