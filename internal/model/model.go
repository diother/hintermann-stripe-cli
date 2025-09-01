package model

import "time"

type Donation struct {
	ID          string
	Created     time.Time
	ClientName  string
	ClientEmail string
	PayoutID    string
	Gross       int
	Fee         int
	Net         int
}

type Payout struct {
	ID      string
	Created time.Time
	Gross   int
	Fee     int
	Net     int
}
