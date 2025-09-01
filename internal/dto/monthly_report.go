package dto

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
