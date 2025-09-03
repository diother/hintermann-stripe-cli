package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type WebhookService interface {
	HandlePayoutReconciliation(object json.RawMessage) error
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

	sigHeader := r.Header.Get("Stripe-Signature")
	if !validateStripeSignature(body, sigHeader, h.WebhookSecret) {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	var event struct {
		Type string `json:"type"`
		Data struct {
			Object json.RawMessage `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "service error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if event.Type != "payout.reconciliation_completed" {
		http.Error(w, "unrecognized event type", http.StatusBadRequest)
		log.Println("unrecognized event:", event.Type)
		return
	}

	if err := h.Service.HandlePayoutReconciliation(event.Data.Object); err != nil {
		http.Error(w, "service error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func validateStripeSignature(payload []byte, sigHeader, secret string) bool {
	var timestamp, signature string
	parts := strings.Split(sigHeader, ",")
	for _, p := range parts {
		if strings.HasPrefix(p, "t=") {
			timestamp = strings.TrimPrefix(p, "t=")
		}
		if strings.HasPrefix(p, "v1=") {
			signature = strings.TrimPrefix(p, "v1=")
		}
	}

	signedPayload := timestamp + "." + string(payload)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}
