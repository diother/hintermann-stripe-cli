package model

import (
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v79"
)

type Donation struct {
	Id          string
	Created     string
	ClientName  string
	ClientEmail string
	PayoutId    string
	Gross       string
	Fee         string
	Net         string
}

func FromChargeTransactionAndPayoutId(charge *stripe.BalanceTransaction, payoutId string) *Donation {
	return &Donation{
		Id:          charge.ID,
		Created:     time.Unix(charge.Created, 0).UTC().Format("2 Jan 2006"),
		ClientName:  charge.Source.Charge.BillingDetails.Name,
		ClientEmail: charge.Source.Charge.BillingDetails.Email,
		PayoutId:    payoutId,
		Gross:       strconv.Itoa(int(charge.Amount)),
		Fee:         strconv.Itoa(int(charge.Fee)),
		Net:         strconv.Itoa(int(charge.Net)),
	}
}

func FromChargeTransactionsAndPayoutId(charges []*stripe.BalanceTransaction, payoutId string) []*Donation {
	donations := make([]*Donation, len(charges))
	for i, d := range charges {
		donations[i] = FromChargeTransactionAndPayoutId(d, payoutId)
	}
	return donations
}
