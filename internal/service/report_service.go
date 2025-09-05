package service

import (
	"time"

	"github.com/diother/hintermann-stripe-cli/internal/dto"
	"github.com/diother/hintermann-stripe-cli/internal/helper"
	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type Reader interface {
	GetPayoutsByMonth(start time.Time) ([]*model.Payout, error)
	GetPayoutById(id string) (*model.Payout, error)
	GetDonationsByPayoutId(payoutId string) ([]*model.Donation, error)
}

type ReportService struct {
	Repo Reader
}

func (s *ReportService) GetMonthlyReport(year int, month time.Month) (*dto.MonthlyReportDTO, error) {
	start := getMonthStart(year, month)

	payouts, err := s.Repo.GetPayoutsByMonth(start)
	if err != nil {
		return nil, err
	}

	gross, fee, net := getMonthlyTotals(payouts)
	payoutDTOs := dto.FromPayouts(payouts)

	return dto.FromMonthTotalsAndPayoutDTOs(start, gross, fee, net, payoutDTOs), nil
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

func getMonthStart(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
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
