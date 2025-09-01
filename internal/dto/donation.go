package dto

import (
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type DonationDTO struct {
	ID          string
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
		ID:          donation.ID,
		Created:     donation.Created.Format("2006-01-02"),
		ClientName:  donation.ClientName,
		ClientEmail: donation.ClientEmail,
		PayoutID:    donation.PayoutID,
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
