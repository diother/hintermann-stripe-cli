package handler

import (
	"fmt"
	"net/http"
)

type WebhookService interface {
	HandlePayoutReconciliation() error
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
	fmt.Println("Webhook hit!")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
