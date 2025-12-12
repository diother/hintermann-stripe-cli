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
			ReconciliationStatus: "completed",
		}, "",
		},
		"nilPayout":            {nil, "is nil"},
		"idMissing":            {&stripe.Payout{ID: ""}, "id is missing"},
		"createdIsNotPositive": {&stripe.Payout{ID: "test_id"}, "created is not positive"},
		"reconciliationStatusNotCompleted": {
			&stripe.Payout{ID: "test_id", Created: 123},
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
		"idMissing":     {&stripe.BalanceTransaction{Type: "payout"}, "id is missing"},
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

func TestValidateChargeTransaction(t *testing.T) {
	testCases := map[string]struct {
		input       *stripe.BalanceTransaction
		expectedErr string
	}{
		"validCharge": {&stripe.BalanceTransaction{
			Type:    "charge",
			ID:      "test_id",
			Created: 123,
			Amount:  100,
			Fee:     10,
			Net:     90,
			Source: &stripe.BalanceTransactionSource{
				Charge: &stripe.Charge{
					BillingDetails: &stripe.ChargeBillingDetails{
						Email: "test@gmail.com",
					},
				},
			},
		}, "",
		},
		"nilCharge":     {nil, "is nil"},
		"typeNotCharge": {&stripe.BalanceTransaction{}, "type is not charge or payment"},
		"idMissing":     {&stripe.BalanceTransaction{Type: "charge"}, "id is missing"},
		"createdNotPositive": {
			&stripe.BalanceTransaction{Type: "charge", ID: "test_id"},
			"created is not positive",
		},
		"amountNotPositive": {
			&stripe.BalanceTransaction{Type: "charge", ID: "test_id", Created: 123},
			"amount is not positive",
		},
		"feeNotPositive": {&stripe.BalanceTransaction{
			Type:    "charge",
			ID:      "test_id",
			Created: 123,
			Amount:  100,
		}, "fee is not positive",
		},
		"netNotPositive": {&stripe.BalanceTransaction{
			Type:    "charge",
			ID:      "test_id",
			Created: 123,
			Amount:  100,
			Fee:     10,
		}, "net is not positive",
		},
		"nilSource": {&stripe.BalanceTransaction{
			Type:    "charge",
			ID:      "test_id",
			Created: 123,
			Amount:  100,
			Fee:     10,
			Net:     90,
		}, "source is nil",
		},
		"nilChargeObject": {&stripe.BalanceTransaction{
			Type:    "charge",
			ID:      "test_id",
			Created: 123,
			Amount:  100,
			Fee:     10,
			Net:     90,
			Source:  &stripe.BalanceTransactionSource{},
		}, "charge object is nil",
		},
		"nilBillingDetails": {&stripe.BalanceTransaction{
			Type:    "charge",
			ID:      "test_id",
			Created: 123,
			Amount:  100,
			Fee:     10,
			Net:     90,
			Source:  &stripe.BalanceTransactionSource{Charge: &stripe.Charge{}},
		}, "billing details is nil",
		},
		"emailMissing": {&stripe.BalanceTransaction{
			Type:    "charge",
			ID:      "test_id",
			Created: 123,
			Amount:  100,
			Fee:     10,
			Net:     90,
			Source: &stripe.BalanceTransactionSource{
				Charge: &stripe.Charge{
					BillingDetails: &stripe.ChargeBillingDetails{},
				},
			},
		}, "email is missing",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateChargeTransaction(tc.input)
			if tc.expectedErr == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestValidateChargeTransactions(t *testing.T) {
	testCases := map[string]struct {
		input       []*stripe.BalanceTransaction
		expectedErr string
	}{
		"emptySlice": {
			input:       []*stripe.BalanceTransaction{},
			expectedErr: "slice is nil",
		},
		"validCharges": {[]*stripe.BalanceTransaction{
			{
				Type:    "charge",
				ID:      "test_id",
				Created: 123,
				Amount:  100,
				Fee:     10,
				Net:     90,
				Source: &stripe.BalanceTransactionSource{
					Charge: &stripe.Charge{
						BillingDetails: &stripe.ChargeBillingDetails{
							Email: "test@gmail.com",
							Name:  "John Doe",
						},
					},
				},
			},
			{
				Type:    "charge",
				ID:      "test_id",
				Created: 123,
				Amount:  200,
				Fee:     20,
				Net:     180,
				Source: &stripe.BalanceTransactionSource{
					Charge: &stripe.Charge{
						BillingDetails: &stripe.ChargeBillingDetails{
							Email: "test@gmail.com",
							Name:  "John Doe",
						},
					},
				},
			},
		}, "",
		},
		"invalidCharge": {
			[]*stripe.BalanceTransaction{{Type: "charge", ID: ""}},
			"index 0 id is missing",
		},
		"stripeFee": {
			[]*stripe.BalanceTransaction{{ID: "fee_1", Type: "stripe_fee"}},
			"stripe_fee transactions were not expected",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateChargeTransactions(tc.input)
			if tc.expectedErr == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestValidateMatchingSums(t *testing.T) {
	testCases := map[string]struct {
		payout        *stripe.BalanceTransaction
		charges       []*stripe.BalanceTransaction
		expectedGross int
		expectedFee   int
		expectedNet   int
		expectedErr   string
	}{
		"matchingSums": {
			payout: &stripe.BalanceTransaction{
				ID:     "po_1",
				Type:   "payout",
				Amount: -294,
			},
			charges: []*stripe.BalanceTransaction{
				{Amount: 100, Fee: 3},
				{Amount: 200, Fee: 3},
			},
			expectedGross: 300,
			expectedFee:   6,
			expectedNet:   294,
			expectedErr:   "",
		},
		"nonMatchingSums": {
			payout: &stripe.BalanceTransaction{
				ID:     "po_2",
				Type:   "payout",
				Amount: -295,
			},
			charges: []*stripe.BalanceTransaction{
				{Amount: 100, Fee: 3},
				{Amount: 200, Fee: 3},
			},
			expectedGross: 0,
			expectedFee:   0,
			expectedNet:   0,
			expectedErr:   "payout amount does not match total charges minus fees. amount 295 != net 294",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gross, fee, net, err := validateMatchingSums(tc.payout, tc.charges)

			if tc.expectedErr == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectedErr != "" {
				if err == nil || err.Error() != tc.expectedErr {
					t.Errorf("Expected error: %v, got: %v", tc.expectedErr, err)
				}
			}
			if gross != tc.expectedGross {
				t.Errorf("Expected gross %v, got %v", tc.expectedGross, gross)
			}
			if fee != tc.expectedFee {
				t.Errorf("Expected fee %v, got %v", tc.expectedFee, fee)
			}
			if net != tc.expectedNet {
				t.Errorf("Expected net %v, got %v", tc.expectedNet, net)
			}
		})
	}
}
