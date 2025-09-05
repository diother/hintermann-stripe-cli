package helper

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

func PayoutReportDir(payoutId string) string {
	return filepath.Join(distDir, "payout_reports", payoutId)
}

func PayoutReportPath(payoutId string) string {
	return filepath.Join(PayoutReportDir(payoutId), "payout_report.pdf")
}

func InvoicePath(payoutId, donationId string) string {
	return filepath.Join(PayoutReportDir(payoutId), "invoices",
		fmt.Sprintf("invoice_%s.pdf", donationId))
}

func EnsureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}
