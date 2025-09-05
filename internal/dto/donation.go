package dto

import (
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/helper"
	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type DonationDTO struct {
	Id          string
	Created     string
	ClientName  string
	ClientEmail string
	PayoutId    string
	Gross       string
	Fee         string
	Net         string
}

func FromDonation(donation *model.Donation) *DonationDTO {
	g := helper.MustAtoi(donation.Gross)
	f := helper.MustAtoi(donation.Fee)
	n := helper.MustAtoi(donation.Net)

	return &DonationDTO{
		Id:          donation.Id,
		Created:     donation.Created,
		ClientName:  donation.ClientName,
		ClientEmail: donation.ClientEmail,
		PayoutId:    donation.PayoutId,
		Gross:       fmt.Sprintf("%.2f lei", float64(g)/100),
		Fee:         fmt.Sprintf("%.2f lei", float64(f)/100),
		Net:         fmt.Sprintf("%.2f lei", float64(n)/100),
	}
}

func FromDonations(donations []*model.Donation) []*DonationDTO {
	donationDTOs := make([]*DonationDTO, len(donations))
	for i, d := range donations {
		donationDTOs[i] = FromDonation(d)
	}
	return donationDTOs
}
