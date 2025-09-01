package dto

import (
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type PayoutDTO struct {
	ID      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func FromPayout(payout *model.Payout) *PayoutDTO {
	return &PayoutDTO{
		ID:      payout.ID,
		Created: payout.Created.Format("2006-01-02"),
		Gross:   fmt.Sprintf("%.2f lei", float64(payout.Gross)/100),
		Fee:     fmt.Sprintf("%.2f lei", float64(payout.Fee)/100),
		Net:     fmt.Sprintf("%.2f lei", float64(payout.Net)/100),
	}
}

func FromPayouts(payouts []*model.Payout) []*PayoutDTO {
	payoutDTOs := make([]*PayoutDTO, len(payouts))
	for i, p := range payouts {
		payoutDTOs[i] = FromPayout(p)
	}
	return payoutDTOs
}
