package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const distDir = "dist"

func MonthlyReportPath(year int, month time.Month) string {
	filename := fmt.Sprintf("monthly_report_%d_%02d.pdf", year, month)
	return filepath.Join(distDir, "monthly_reports", filename)
}

func PayoutReportDir(payoutID string) string {
	return filepath.Join(distDir, "payout_reports", payoutID)
}

func PayoutReportPath(payoutID string) string {
	return filepath.Join(PayoutReportDir(payoutID), "payout_report.pdf")
}

func InvoicePath(payoutID, donationID string) string {
	return filepath.Join(PayoutReportDir(payoutID), "invoices",
		fmt.Sprintf("invoice_%s.pdf", donationID))
}

func EnsureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}
