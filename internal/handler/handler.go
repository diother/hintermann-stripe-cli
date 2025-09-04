package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

type WebhookService interface {
	HandlePayoutReconciliation(payout *stripe.Payout) error
}

type WebhookHandler struct {
	Service       WebhookService
	WebhookSecret string
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.WebhookSecret)
	if err != nil {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		log.Println("invalid signature:", err)
		return
	}

	if event.Type != "payout.reconciliation_completed" {
		http.Error(w, "unrecognized event type", http.StatusBadRequest)
		log.Println("unrecognized event:", event.Type)
		return
	}

	payout := &stripe.Payout{}
	if err = json.Unmarshal(event.Data.Raw, payout); err != nil {
		http.Error(w, "unrecognized data object", http.StatusBadRequest)
		log.Println("unrecognized data object:", err)
		return
	}

	if err := h.Service.HandlePayoutReconciliation(payout); err != nil {
		http.Error(w, "service error", http.StatusInternalServerError)
		log.Println("service error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
