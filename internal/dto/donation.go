package dto

import (
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type DonationDTO struct {
	Id          string
	Created     string
	ClientName  string
	ClientEmail string
	PayoutID    string
	Gross       string
	Fee         string
	Net         string
}

func FromDonation(donation *model.Donation) *DonationDTO {
	return &DonationDTO{
		Id:          donation.Id,
		Created:     donation.Created.Format("2 Jan 2006"),
		ClientName:  donation.ClientName,
		ClientEmail: donation.ClientEmail,
		PayoutID:    donation.PayoutId,
		Gross:       fmt.Sprintf("%.2f lei", float64(donation.Gross)/100),
		Fee:         fmt.Sprintf("%.2f lei", float64(donation.Fee)/100),
		Net:         fmt.Sprintf("%.2f lei", float64(donation.Net)/100),
	}
}

func FromDonations(donations []*model.Donation) []*DonationDTO {
	donationDTOs := make([]*DonationDTO, len(donations))
	for i, d := range donations {
		donationDTOs[i] = FromDonation(d)
	}
	return donationDTOs
}
