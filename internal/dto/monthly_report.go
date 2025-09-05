package dto

import (
	"fmt"
	"time"
)

type MonthlyReportDTO struct {
	MonthStart string
	MonthEnd   string
	Issued     string
	Gross      string
	Fee        string
	Net        string
	Payouts    []*PayoutDTO
}

func FromMonthTotalsAndPayoutDTOs(start time.Time, gross, fee, net int, payoutDTOs []*PayoutDTO) *MonthlyReportDTO {
	end := start.AddDate(0, 1, -1)
	issued := start.AddDate(0, 1, 0)

	return &MonthlyReportDTO{
		MonthStart: start.Format("2 Jan 2006"),
		MonthEnd:   end.Format("2 Jan 2006"),
		Issued:     issued.Format("2 Jan 2006"),
		Gross:      fmt.Sprintf("%.2f lei", float64(gross)/100),
		Fee:        fmt.Sprintf("%.2f lei", float64(fee)/100),
		Net:        fmt.Sprintf("%.2f lei", float64(net)/100),
		Payouts:    payoutDTOs,
	}
}
