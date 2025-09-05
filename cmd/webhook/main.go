package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/diother/hintermann-stripe-cli/internal/handler"
	"github.com/diother/hintermann-stripe-cli/internal/repo"
	"github.com/diother/hintermann-stripe-cli/internal/service"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	stripeKey := os.Getenv("STRIPE_SECRET")
	webhookSecret := os.Getenv("WEBHOOK_SECRET")

	if stripeKey == "" || webhookSecret == "" {
		log.Fatal("stripe keys missing")
	}
	stripe.Key = stripeKey

	repo := &repo.CSVRepo{}
	service := &service.WebhookService{Repo: repo}
	handler := &handler.WebhookHandler{
		Service:       service,
		WebhookSecret: webhookSecret,
	}
	http.Handle("/webhook", handler)

	fmt.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
