package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/diother/hintermann-stripe-cli/internal/handler"
	"github.com/diother/hintermann-stripe-cli/internal/repo"
	"github.com/diother/hintermann-stripe-cli/internal/service"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	stripeKey := os.Getenv("STRIPE_SECRET")
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	dataDir := os.Getenv("DATA_DIR")

	if stripeKey == "" || webhookSecret == "" || dataDir == "" {
		log.Fatal("env variables are missing")
	}
	stripe.Key = stripeKey

	repo := &repo.CSVRepo{
		DonationsFile: filepath.Join(dataDir, "donations.csv"),
		PayoutsFile:   filepath.Join(dataDir, "payouts.csv"),
	}
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
