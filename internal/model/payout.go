package model

import "time"

type Payout struct {
	ID      string
	Created time.Time
	Gross   int
	Fee     int
	Net     int
}

func NewPayout(id string, created time.Time, gross, fee, net int) *Payout {
	return &Payout{
		ID:      id,
		Created: created,
		Gross:   gross,
		Fee:     fee,
		Net:     net,
	}
}
