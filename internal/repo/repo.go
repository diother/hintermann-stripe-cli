package repo

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type CSVRepo struct {
	DonationsFile string
	PayoutsFile   string
}

func (r *CSVRepo) GetPayoutsByDateRange(start, end time.Time) ([]*model.Payout, error) {
	payouts, err := r.loadPayouts()
	if err != nil {
		return nil, err
	}
	var filtered []*model.Payout
	for _, p := range payouts {
		if (p.Created.Equal(start) || p.Created.After(start)) &&
			(p.Created.Equal(end) || p.Created.Before(end)) {
			filtered = append(filtered, p)
		}
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no payouts found for %d-%02d", start.Year(), start.Month())
	}
	return filtered, nil
}

func (r *CSVRepo) GetPayoutByID(id string) (*model.Payout, error) {
	payouts, err := r.loadPayouts()
	if err != nil {
		return nil, err
	}
	for _, p := range payouts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("payout not found: %s", id)
}

func (r *CSVRepo) GetDonationsByPayoutID(payoutID string) ([]*model.Donation, error) {
	donations, err := r.loadDonations()
	if err != nil {
		return nil, err
	}
	var filtered []*model.Donation
	for _, d := range donations {
		if d.PayoutID == payoutID {
			filtered = append(filtered, d)
		}
	}
	return filtered, nil
}

func (r *CSVRepo) loadDonations() ([]*model.Donation, error) {
	file, err := os.Open(r.DonationsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var donations []*model.Donation
	for i, record := range records {
		if i == 0 { // skip header
			continue
		}
		created, err := parseDate(record[1])
		if err != nil {
			return nil, fmt.Errorf("bad date at line %d: %w", i+1, err)
		}
		gross, _ := parseInt(record[5])
		fee, _ := parseInt(record[6])
		net, _ := parseInt(record[7])

		donations = append(donations, &model.Donation{
			ID:          record[0],
			Created:     created,
			ClientName:  record[2],
			ClientEmail: record[3],
			PayoutID:    record[4],
			Gross:       gross,
			Fee:         fee,
			Net:         net,
		})
	}

	return donations, nil
}

func (r *CSVRepo) loadPayouts() ([]*model.Payout, error) {
	file, err := os.Open(r.PayoutsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var payouts []*model.Payout
	for i, record := range records {
		if i == 0 { // skip header
			continue
		}
		created, err := parseDate(record[1])
		if err != nil {
			return nil, fmt.Errorf("bad date at line %d: %w", i+1, err)
		}
		gross, _ := parseInt(record[2])
		fee, _ := parseInt(record[3])
		net, _ := parseInt(record[4])

		payouts = append(payouts, &model.Payout{
			ID:      record[0],
			Created: created,
			Gross:   gross,
			Fee:     fee,
			Net:     net,
		})
	}

	return payouts, nil
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value) // strict YYYY-MM-DD
}

func parseInt(value string) (int, error) {
	return strconv.Atoi(value)
}
