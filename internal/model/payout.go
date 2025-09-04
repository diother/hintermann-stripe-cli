package model

import (
	"time"

	"github.com/stripe/stripe-go/v79"
)

type Payout struct {
	ID      string
	Created time.Time
	Gross   int
	Fee     int
	Net     int
}

func FromStripePayoutAndTotals(payout *stripe.Payout, gross, fee, net int) *Payout {
	return &Payout{
		ID:      payout.ID,
		Created: time.Unix(payout.Created, 0).UTC(),
		Gross:   gross,
		Fee:     fee,
		Net:     net,
	}
}
