package repo

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/diother/hintermann-stripe-cli/internal/model"
)

const (
	donationsFile = "data/donations.csv"
	payoutsFile   = "data/payouts.csv"
)

type CSVRepo struct{}

func (r *CSVRepo) GetPayoutsByMonth(start time.Time) ([]*model.Payout, error) {
	payouts, err := r.loadPayouts()
	if err != nil {
		return nil, err
	}
	var filtered []*model.Payout
	for _, p := range payouts {
		created, err := time.Parse("2 Jan 2006", p.Created)
		if err != nil {
			return nil, fmt.Errorf("invalid time format for %s", p.Id)
		}
		end := start.AddDate(0, 1, -1)
		if (created.Equal(start) || created.After(start)) &&
			(created.Equal(end) || created.Before(end)) {
			filtered = append(filtered, p)
		}
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no payouts found for %d-%02d", start.Year(), start.Month())
	}
	return filtered, nil
}

func (r *CSVRepo) GetPayoutById(id string) (*model.Payout, error) {
	payouts, err := r.loadPayouts()
	if err != nil {
		return nil, err
	}
	for _, p := range payouts {
		if p.Id == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("payout not found: %s", id)
}

func (r *CSVRepo) GetDonationsByPayoutId(payoutId string) ([]*model.Donation, error) {
	donations, err := r.loadDonations()
	if err != nil {
		return nil, err
	}
	var filtered []*model.Donation
	for _, d := range donations {
		if d.PayoutId == payoutId {
			filtered = append(filtered, d)
		}
	}
	return filtered, nil
}

func (r *CSVRepo) loadDonations() ([]*model.Donation, error) {
	file, err := os.Open(donationsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	donations := make([]*model.Donation, len(records)-1)
	for i, record := range records[1:] {
		donations[i] = &model.Donation{
			Id:          record[0],
			Created:     record[1],
			ClientName:  record[2],
			ClientEmail: record[3],
			PayoutId:    record[4],
			Gross:       record[5],
			Fee:         record[6],
			Net:         record[7],
		}
	}
	return donations, nil
}

func (r *CSVRepo) loadPayouts() ([]*model.Payout, error) {
	file, err := os.Open(payoutsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	payouts := make([]*model.Payout, len(records)-1)
	for i, record := range records[1:] {
		payouts[i] = &model.Payout{
			Id:      record[0],
			Created: record[1],
			Gross:   record[2],
			Fee:     record[3],
			Net:     record[4],
		}
	}
	return payouts, nil
}
