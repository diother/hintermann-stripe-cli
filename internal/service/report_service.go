package service

import (
	"time"

	"github.com/diother/hintermann-stripe-cli/internal/dto"
	"github.com/diother/hintermann-stripe-cli/internal/helper"
	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type Reader interface {
	GetPayoutsByDateRange(start, end time.Time) ([]*model.Payout, error)
	GetPayoutById(id string) (*model.Payout, error)
	GetDonationsByPayoutId(payoutId string) ([]*model.Donation, error)
}

type ReportService struct {
	Repo Reader
}

func (s *ReportService) GetMonthlyReport(year int, month time.Month) (*dto.MonthlyReportDTO, error) {
	start, end := getDateRange(year, month)

	payouts, err := s.Repo.GetPayoutsByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	gross, fee, net := getMonthlyTotals(payouts)
	payoutDTOs := dto.FromPayouts(payouts)

	return dto.FromDateTotalsAndPayoutDTOs(start, gross, fee, net, payoutDTOs), nil
}

func (s *ReportService) GetPayoutReport(payoutId string) (*dto.PayoutReportDTO, []*dto.DonationDTO, error) {
	payout, err := s.Repo.GetPayoutById(payoutId)
	if err != nil {
		return nil, nil, err
	}

	donations, err := s.Repo.GetDonationsByPayoutId(payoutId)
	if err != nil {
		return nil, nil, err
	}

	donationDTOs := dto.FromDonations(donations)
	payoutReport := dto.FromPayoutWithDonations(payout, donations)

	return payoutReport, donationDTOs, nil
}

func getMonthlyTotals(payouts []*model.Payout) (int, int, int) {
	var gross, fee, net int
	for _, p := range payouts {
		gross += helper.MustAtoi(p.Gross)
		fee += helper.MustAtoi(p.Fee)
		net += helper.MustAtoi(p.Net)
	}
	return gross, fee, net
}

func getDateRange(year int, month time.Month) (time.Time, time.Time) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	return start, end
}
