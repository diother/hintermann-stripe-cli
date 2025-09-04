package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/diother/hintermann-stripe-cli/internal/handler"
	"github.com/diother/hintermann-stripe-cli/internal/repo"
	"github.com/diother/hintermann-stripe-cli/internal/service"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	stripe.Key = ""

	repo := &repo.CSVRepo{}
	service := &service.WebhookService{Repo: repo}
	handler := &handler.WebhookHandler{
		Service:       service,
		WebhookSecret: "",
	}

	http.Handle("/webhook", handler)

	fmt.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
