package service

import (
	"encoding/json"
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/model"
)

type Writer interface {
	AppendDonation(d *model.Donation) error
	AppendPayout(p *model.Payout) error
}

type WebhookService struct {
	Repo         Writer
	StripeSecret string
}

func (s *WebhookService) HandlePayoutReconciliation(object json.RawMessage) error {
	fmt.Println(object)
	return nil
}
