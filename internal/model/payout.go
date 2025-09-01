package model

import "time"

type Payout struct {
	ID      string
	Created time.Time
	Gross   int
	Fee     int
	Net     int
}
