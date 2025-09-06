package repo

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/diother/hintermann-stripe-cli/internal/model"
)

func (r *CSVRepo) WritePayoutAndDonations(p *model.Payout, ds []*model.Donation) error {
	existingIds, err := readExistingPayoutIds(r.PayoutsFile)
	if err != nil {
		return fmt.Errorf("failed to read existing payout IDs: %w", err)
	}

	if _, exists := existingIds[p.Id]; exists {
		return nil
	}

	payoutRow := [][]string{
		{p.Id, p.Created, p.Gross, p.Fee, p.Net},
	}
	if err := appendWithTemp(r.PayoutsFile, payoutRow); err != nil {
		return fmt.Errorf("failed to append payout: %w", err)
	}

	donationRows := make([][]string, len(ds))
	for i, d := range ds {
		donationRows[i] = []string{
			d.Id,
			d.Created,
			d.ClientName,
			d.ClientEmail,
			d.PayoutId,
			d.Gross,
			d.Fee,
			d.Net,
		}
	}

	if err := appendWithTemp(r.DonationsFile, donationRows); err != nil {
		return fmt.Errorf("failed to append donations: %w", err)
	}

	return nil
}

func appendWithTemp(filename string, newRows [][]string) error {
	tmpFile := filename + ".tmp"

	// 1. create temp file
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	defer func() {
		w.Flush()
		f.Close()
		if err != nil {
			os.Remove(tmpFile)
		}
	}()

	// 2. copy existing CSV
	if _, err := os.Stat(filename); err == nil {
		existing, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer existing.Close()
		r := csv.NewReader(existing)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if err := w.Write(record); err != nil {
				return err
			}
		}
	}

	// 3. append new rows
	for _, row := range newRows {
		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	// 4. atomic replace
	return os.Rename(tmpFile, filename)
}

func readExistingPayoutIds(filename string) (map[string]struct{}, error) {
	ids := make(map[string]struct{})

	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return ids, nil
		}
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) > 0 {
			ids[record[0]] = struct{}{}
		}
	}
	return ids, nil
}
