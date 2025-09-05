package dto

import (
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/helper"
	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type PayoutDTO struct {
	Id      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func FromPayout(payout *model.Payout) *PayoutDTO {
	g := helper.MustAtoi(payout.Gross)
	f := helper.MustAtoi(payout.Fee)
	n := helper.MustAtoi(payout.Net)

	return &PayoutDTO{
		Id:      payout.Id,
		Created: payout.Created,
		Gross:   fmt.Sprintf("%.2f lei", float64(g)/100),
		Fee:     fmt.Sprintf("%.2f lei", float64(f)/100),
		Net:     fmt.Sprintf("%.2f lei", float64(n)/100),
	}
}

func FromPayouts(payouts []*model.Payout) []*PayoutDTO {
	payoutDTOs := make([]*PayoutDTO, len(payouts))
	for i, p := range payouts {
		payoutDTOs[i] = FromPayout(p)
	}
	return payoutDTOs
}
