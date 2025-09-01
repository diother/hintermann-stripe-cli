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

func NewDonation(id string, created time.Time, clientName, clientEmail, payoutID string, gross, fee, net int) *Donation {
	return &Donation{
		ID:          id,
		Created:     created,
		ClientName:  clientName,
		ClientEmail: clientEmail,
		PayoutID:    payoutID,
		Gross:       gross,
		Fee:         fee,
		Net:         net,
	}
}
