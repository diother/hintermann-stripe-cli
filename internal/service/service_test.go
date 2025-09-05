package service

import (
	"testing"

	"github.com/diother/hintermann-stripe-cli/internal/model"
	"github.com/stripe/stripe-go/v79"
)

func TestGetMonthlyTotals(t *testing.T) {
	testCases := map[string]struct {
		input         []*model.Payout
		expectedPanic bool
		expectedGross int
		expectedFee   int
		expectedNet   int
	}{
		"validStrings": {
			input: []*model.Payout{
				{Gross: "100", Fee: "10", Net: "90"},
				{Gross: "200", Fee: "20", Net: "180"},
			},
			expectedPanic: false,
			expectedGross: 300,
			expectedFee:   30,
			expectedNet:   270,
		},
		"badString": {
			input: []*model.Payout{
				{Gross: "bad", Fee: "10", Net: "90"},
			},
			expectedPanic: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic, got none")
					}
				}()
				getMonthlyTotals(tc.input)
			} else {
				gross, fee, net := getMonthlyTotals(tc.input)
				if gross != tc.expectedGross || fee != tc.expectedFee || net != tc.expectedNet {
					t.Errorf("wanted (%d, %d, %d), got (%d, %d, %d)",
						tc.expectedGross, tc.expectedFee, tc.expectedNet,
						gross, fee, net)
				}
			}
		})
	}
}

func TestValidateStripePayout(t *testing.T) {
	testCases := map[string]struct {
		input       *stripe.Payout
		expectedErr string
	}{
		"validPayout": {&stripe.Payout{
			ID:                   "test_id",
			Created:              123,
			Status:               "paid",
			ReconciliationStatus: "completed",
		}, "",
		},
		"nilPayout":            {nil, "is nil"},
		"idMissing":            {&stripe.Payout{ID: ""}, "id missing"},
		"createdIsNotPositive": {&stripe.Payout{ID: "test_id"}, "created is not positive"},
		"statusNotPaid": {
			&stripe.Payout{ID: "test_id", Created: 123},
			"status is not paid",
		},
		"reconciliationStatusNotCompleted": {
			&stripe.Payout{ID: "test_id", Created: 123, Status: "paid"},
			"reconciliation status is not completed",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateStripePayout(tc.input)
			if tc.expectedErr == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestValidatePayoutTransaction(t *testing.T) {
	testCases := map[string]struct {
		input       *stripe.BalanceTransaction
		expectedErr string
	}{
		"validPayout": {&stripe.BalanceTransaction{
			Type:    "payout",
			ID:      "test_id",
			Created: 123,
			Amount:  -100,
			Fee:     0,
			Net:     -100,
		}, "",
		},
		"nilPayout":     {nil, "is nil"},
		"typeNotPayout": {&stripe.BalanceTransaction{}, "type is not payout"},
		"idMissing":     {&stripe.BalanceTransaction{Type: "payout"}, "id missing"},
		"createdNotPositive": {
			&stripe.BalanceTransaction{Type: "payout", ID: "test_id"},
			"created is not positive",
		},
		"amountNotNegative": {
			&stripe.BalanceTransaction{Type: "payout", ID: "test_id", Created: 123},
			"amount is not negative",
		},
		"feeNot0": {&stripe.BalanceTransaction{
			Type:    "payout",
			ID:      "test_id",
			Created: 123,
			Amount:  -100,
			Fee:     10,
		}, "fee is not 0",
		},
		"netNotNegative": {&stripe.BalanceTransaction{
			Type:    "payout",
			ID:      "test_id",
			Created: 123,
			Amount:  -100,
			Fee:     0,
		}, "net is not negative",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validatePayoutTransaction(tc.input)
			if tc.expectedErr == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}
