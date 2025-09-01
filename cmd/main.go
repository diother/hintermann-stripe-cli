package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/diother/hintermann-stripe-cli/internal/pdfgen"
	"github.com/diother/hintermann-stripe-cli/internal/repo"
	"github.com/diother/hintermann-stripe-cli/internal/service"
)

func main() {
	monthly := flag.Bool("monthly", false, "Generate monthly report")
	payoutID := flag.String("payout", "", "Generate payout report by ID")
	year := flag.Int("year", time.Now().Year(), "Year for monthly report")
	month := flag.Int("month", int(time.Now().Month()), "Month for monthly report")
	flag.Parse()

	repo := &repo.CSVRepo{
		DonationsFile: "data/donations.csv",
		PayoutsFile:   "data/payouts.csv",
	}
	service := &service.ReportService{Reader: repo}

	if *monthly {
		report, err := service.GetMonthlyReport(*year, time.Month(*month))
		if err != nil {
			log.Fatal(err)
		}
		path, err := pdfgen.GenerateMonthlyReport(report, *year, time.Month(*month))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Monthly report generated:", path)
	} else if *payoutID != "" {
		payoutReport, donationDTOs, err := service.GetPayoutReport(*payoutID)
		if err != nil {
			log.Fatal(err)
		}
		path, err := pdfgen.GeneratePayoutReport(payoutReport)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Payout report generated:", path)

		for _, inv := range donationDTOs {
			pdfgen.GenerateInvoice(inv)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Invoice generated:", path)
		}
	} else {
		fmt.Println("No action specified. Use -monthly or -payout flags.")
	}
}
