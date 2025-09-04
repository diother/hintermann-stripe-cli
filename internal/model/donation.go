package model

import (
	"time"

	"github.com/stripe/stripe-go/v79"
)

type Donation struct {
	Id          string
	Created     time.Time
	ClientName  string
	ClientEmail string
	PayoutId    string
	Gross       int
	Fee         int
	Net         int
}

func FromChargeTransactionAndPayoutId(charge *stripe.BalanceTransaction, payoutId string) *Donation {
	return &Donation{
		Id:          charge.ID,
		Created:     time.Unix(charge.Created, 0).UTC(),
		ClientName:  charge.Source.Charge.BillingDetails.Name,
		ClientEmail: charge.Source.Charge.BillingDetails.Email,
		PayoutId:    payoutId,
		Gross:       int(charge.Amount),
		Fee:         int(charge.Fee),
		Net:         int(charge.Net),
	}
}

func FromChargeTransactionsAndPayoutId(charges []*stripe.BalanceTransaction, payoutId string) []*Donation {
	donations := make([]*Donation, len(charges))
	for i, d := range charges {
		donations[i] = FromChargeTransactionAndPayoutId(d, payoutId)
	}
	return donations
}
