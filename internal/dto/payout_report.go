package dto

type PayoutReportDTO struct {
	Payout    *PayoutDTO
	Donations []*DonationDTO
}

func NewPayoutReportDTO(payout *PayoutDTO, donations []*DonationDTO) *PayoutReportDTO {
	return &PayoutReportDTO{
		Payout:    payout,
		Donations: donations,
	}
}
