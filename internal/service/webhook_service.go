package service

import (
	"fmt"

	"github.com/diother/hintermann-stripe-cli/internal/model"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
)

type Writer interface {
	WritePayoutAndDonations(p *model.Payout, ds []*model.Donation) error
}

type WebhookService struct {
	Repo Writer
}

func (s *WebhookService) HandlePayoutReconciliation(stripePayout *stripe.Payout) error {
	if err := validateStripePayout(stripePayout); err != nil {
		return fmt.Errorf("stripe payout invalid: %w", err)
	}
	payoutTransaction, chargeTransactions, err := fetchRelatedTransactions(stripePayout.ID)
	if err != nil {
		return fmt.Errorf("transactions fetch failed: %w", err)
	}
	if err := validatePayoutTransaction(payoutTransaction); err != nil {
		return fmt.Errorf("payout transaction invalid: %w", err)
	}
	if err := validateChargeTransactions(chargeTransactions); err != nil {
		return fmt.Errorf("charge transactions invalid: %w", err)
	}
	gross, fee, net, err := validateMatchingSums(payoutTransaction, chargeTransactions)
	if err != nil {
		return fmt.Errorf("matching sum validation failed: %w", err)
	}

	payout := model.FromStripePayoutAndTotals(stripePayout, gross, fee, net)
	donations := model.FromChargeTransactionsAndPayoutId(chargeTransactions, stripePayout.ID)

	if err := s.Repo.WritePayoutAndDonations(payout, donations); err != nil {
		return fmt.Errorf("failed to persist payout+donations: %w", err)
	}
	return nil
}

func fetchRelatedTransactions(id string) (*stripe.BalanceTransaction, []*stripe.BalanceTransaction, error) {
	params := &stripe.BalanceTransactionListParams{}
	params.Payout = &id
	params.AddExpand("data.source")

	iter := balancetransaction.List(params)

	var payout *stripe.BalanceTransaction
	if iter.Next() {
		payout = iter.BalanceTransaction()
	}

	var charges []*stripe.BalanceTransaction
	for iter.Next() {
		charges = append(charges, iter.BalanceTransaction())
	}

	if err := iter.Err(); err != nil {
		return nil, nil, err
	}
	return payout, charges, nil
}

func validateStripePayout(payout *stripe.Payout) error {
	if payout == nil {
		return fmt.Errorf("is nil")
	}
	if payout.ID == "" {
		return fmt.Errorf("id missing")
	}
	if payout.Created <= 0 {
		return fmt.Errorf("created invalid")
	}
	if payout.Status != "paid" {
		return fmt.Errorf("status invalid")
	}
	if payout.ReconciliationStatus != "completed" {
		return fmt.Errorf("reconciliation status invalid")
	}
	return nil
}

func validatePayoutTransaction(payout *stripe.BalanceTransaction) error {
	if payout == nil {
		return fmt.Errorf("is nil")
	}
	if payout.Type != "payout" {
		return fmt.Errorf("type invalid")
	}
	if payout.ID == "" {
		return fmt.Errorf("id missing")
	}
	if payout.Created <= 0 {
		return fmt.Errorf("created invalid")
	}
	if payout.Amount >= 0 {
		return fmt.Errorf("amount invalid")
	}
	if payout.Fee != 0 {
		return fmt.Errorf("fee invalid")
	}
	if payout.Net >= 0 {
		return fmt.Errorf("net invalid")
	}
	return nil
}

func validateChargeTransactions(charges []*stripe.BalanceTransaction) error {
	if len(charges) == 0 {
		return fmt.Errorf("slice is nil")
	}
	for i, charge := range charges {
		if err := validateChargeTransaction(charge); err != nil {
			if charge.Type == "stripe_fee" {
				return fmt.Errorf("stripe_fee transactions were not expected")
			}
			return fmt.Errorf("index %d %w", i, err)
		}
	}
	return nil
}

func validateChargeTransaction(charge *stripe.BalanceTransaction) error {
	if charge == nil {
		return fmt.Errorf("is nil")
	}
	if charge.Type != "charge" {
		return fmt.Errorf("type invalid")
	}
	if charge.ID == "" {
		return fmt.Errorf("id missing")
	}
	if charge.Created <= 0 {
		return fmt.Errorf("created invalid")
	}
	if charge.Amount <= 0 {
		return fmt.Errorf("amount invalid")
	}
	if charge.Fee <= 0 {
		return fmt.Errorf("fee invalid")
	}
	if charge.Net <= 0 {
		return fmt.Errorf("net invalid")
	}
	if charge.Source == nil {
		return fmt.Errorf("source is nil")
	}
	if charge.Source.Charge == nil {
		return fmt.Errorf("charge object is nil")
	}
	if charge.Source.Charge.BillingDetails == nil {
		return fmt.Errorf("billing details is nil")
	}
	if charge.Source.Charge.BillingDetails.Email == "" {
		return fmt.Errorf("email missing")
	}
	if charge.Source.Charge.BillingDetails.Name == "" {
		return fmt.Errorf("name missing")
	}
	return nil
}

func validateMatchingSums(payout *stripe.BalanceTransaction, charges []*stripe.BalanceTransaction) (int, int, int, error) {
	var gross, fee, net int

	for _, charge := range charges {
		gross += int(charge.Amount)
		fee += int(charge.Fee)
	}

	net = gross - fee
	payoutAmount := int(-payout.Amount)

	if payoutAmount != net {
		return 0, 0, 0, fmt.Errorf("payout amount does not match total charges minus fees. amount %v != net %v", payoutAmount, net)
	}
	return gross, fee, net, nil
}
