package service

import (
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/model"
	"github.com/stripe/stripe-go/v79"
)

type Writer interface {
	AppendDonation(d *model.Donation) error
	AppendPayout(p *model.Payout) error
}

type WebhookService struct {
	Repo Writer
}

func (s *WebhookService) HandlePayoutReconciliation(stripePayout *stripe.Payout) error {
	if err := validatePayout(stripePayout); err != nil {
		return fmt.Errorf("payout validation error: %w", err)
	}
	return nil
}

func validatePayout(payout *stripe.Payout) error {
	if payout == nil {
		return fmt.Errorf("missing")
	}
	if payout.ID == "" {
		return fmt.Errorf("id missing")
	}
	if payout.Status != "paid" {
		return fmt.Errorf("status invalid")
	}
	if payout.ReconciliationStatus != "completed" {
		return fmt.Errorf("reconciliation status invalid")
	}
	return nil
}
