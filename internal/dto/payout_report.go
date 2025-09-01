package dto

import "github.com/diother/hintermann-stripe-cli/internal/model"

type PayoutReportDTO struct {
	Payout    *PayoutDTO
	Donations []*DonationDTO
}

func FromPayoutWithDonations(payout *model.Payout, donations []*model.Donation) *PayoutReportDTO {
	return &PayoutReportDTO{
		Payout:    FromPayout(payout),
		Donations: FromDonations(donations),
	}
}
