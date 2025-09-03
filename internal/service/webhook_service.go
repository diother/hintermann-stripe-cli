package service

import "github.com/diother/hintermann-stripe-cli/internal/model"

type Writer interface {
	AppendDonation(d *model.Donation) error
	AppendPayout(p *model.Payout) error
}

type WebhookService struct {
	Repo Writer
}

func (s *WebhookService) HandlePayoutReconciliation() error {
	return nil
}
