package model

import (
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v79"
)

type Payout struct {
	Id      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func FromStripePayoutAndTotals(payout *stripe.Payout, gross, fee, net int) *Payout {
	return &Payout{
		Id:      payout.ID,
		Created: time.Unix(payout.Created, 0).UTC().Format("2 Jan 2006"),
		Gross:   strconv.Itoa(gross),
		Fee:     strconv.Itoa(fee),
		Net:     strconv.Itoa(net),
	}
}
