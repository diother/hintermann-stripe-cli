package models

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

func NewDonationDTO(id, created, clientName, clientEmail, payoutID, gross, fee, net string) *DonationDTO {
	return &DonationDTO{
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

type PayoutDTO struct {
	ID      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func NewPayoutDTO(id, created, gross, fee, net string) *PayoutDTO {
	return &PayoutDTO{
		ID:      id,
		Created: created,
		Gross:   gross,
		Fee:     fee,
		Net:     net,
	}
}

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

type MonthlyReportDTO struct {
	MonthStart string
	MonthEnd   string
	Issued     string
	Gross      string
	Fee        string
	Net        string
	Payouts    []*PayoutDTO
}

func NewMonthlyReportDTO(monthStart, monthEnd, issued, gross, fee, net string, payouts []*PayoutDTO) *MonthlyReportDTO {
	return &MonthlyReportDTO{
		MonthStart: monthStart,
		MonthEnd:   monthEnd,
		Issued:     issued,
		Gross:      gross,
		Fee:        fee,
		Net:        net,
		Payouts:    payouts,
	}
}
