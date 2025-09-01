package models

type PayoutReport struct {
	Payout    *Payout
	Donations []*Donation
}

func NewPayoutReport(payout *Payout, donations []*Donation) *PayoutReport {
	return &PayoutReport{
		Payout:    payout,
		Donations: donations,
	}
}
