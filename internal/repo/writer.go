package repo

import "github.com/diother/hintermann-stripe-cli/internal/model"

func (r *CSVRepo) AppendDonation(d *model.Donation) error { return nil }
func (r *CSVRepo) AppendPayout(p *model.Payout) error     { return nil }
