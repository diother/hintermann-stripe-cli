package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/diother/hintermann-stripe-cli/internal/handler"
	"github.com/diother/hintermann-stripe-cli/internal/repo"
	"github.com/diother/hintermann-stripe-cli/internal/service"
)

func main() {
	repo := &repo.CSVRepo{}
	service := &service.WebhookService{Repo: repo}
	handler := &handler.WebhookHandler{
		Service:      service,
		StripeSecret: "",
	}

	http.Handle("/webhook", handler)

	fmt.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
