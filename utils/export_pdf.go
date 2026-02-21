package utils

import (
	"bytes"
	"finance-tracking-app/services"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// ExportMonthlyReportToPDF generates a PDF report for a given monthly report
func ExportMonthlyReportToPDF(report *services.MonthlyReport) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font for title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Monthly Financial Report - %s %d",
		time.Month(report.Month).String(), report.Year))
	pdf.Ln(12)

	// Income Section
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Income Summary", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total Income: Rp %.2f", report.Income.Total))
	pdf.Ln(8)

	// Income sources
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(90, 7, "Source")
	pdf.Cell(90, 7, "Amount")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 10)
	for _, source := range report.Income.Sources {
		pdf.Cell(90, 6, source.Source)
		pdf.Cell(90, 6, fmt.Sprintf("Rp %.2f", source.Amount))
		pdf.Ln(6)
	}
	pdf.Ln(5)

	// Expenses Section
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Expense Summary", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total Expenses: Rp %.2f", report.Expenses.Total))
	pdf.Ln(8)

	// Expenses by category
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(70, 7, "Category")
	pdf.Cell(60, 7, "Amount")
	pdf.Cell(50, 7, "Count")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 10)
	for _, cat := range report.Expenses.ByCategory {
		pdf.Cell(70, 6, cat.Category)
		pdf.Cell(60, 6, fmt.Sprintf("Rp %.2f", cat.Amount))
		pdf.Cell(50, 6, fmt.Sprintf("%d", cat.Count))
		pdf.Ln(6)
	}
	pdf.Ln(5)

	// Budget Section
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Budget Summary", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total Allocated: Rp %.2f", report.Budget.TotalAllocated))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Total Spent: Rp %.2f", report.Budget.TotalSpent))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Total Remaining: Rp %.2f", report.Budget.TotalRemaining))
	pdf.Ln(10)

	// Category budgets
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 7, "Category")
	pdf.Cell(40, 7, "Allocated")
	pdf.Cell(40, 7, "Spent")
	pdf.Cell(50, 7, "Remaining")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 10)
	for _, cat := range report.Budget.CategoriesSummary {
		pdf.Cell(50, 6, cat.Category)
		pdf.Cell(40, 6, fmt.Sprintf("%.2f", cat.Allocated))
		pdf.Cell(40, 6, fmt.Sprintf("%.2f", cat.Spent))
		pdf.Cell(50, 6, fmt.Sprintf("%.2f", cat.Remaining))
		pdf.Ln(6)
	}
	pdf.Ln(5)

	// Financial Summary
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Financial Summary", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Savings: Rp %.2f", report.Savings))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Unallocated: Rp %.2f", report.Unallocated))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Top Spending Category: %s (Rp %.2f)",
		report.TopCategory, report.TopCategoryAmount))
	pdf.Ln(10)

	// Insights
	if len(report.Insights) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(0, 10, "Insights", "1", 0, "L", true, 0, "")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, insight := range report.Insights {
			pdf.MultiCell(0, 6, "- "+insight, "", "", false)
		}
	}

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return &buf, nil
}

// ExportYearlyReportToPDF generates a PDF report for a given yearly report
func ExportYearlyReportToPDF(report *services.YearlyReport) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font for title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Yearly Financial Report - %d", report.Year))
	pdf.Ln(12)

	// Overall Summary
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Annual Summary", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total Income: Rp %.2f", report.TotalIncome))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Total Expenses: Rp %.2f", report.TotalExpenses))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Total Savings: Rp %.2f", report.TotalSavings))
	pdf.Ln(10)

	// Monthly Breakdown
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Monthly Breakdown", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(40, 7, "Month")
	pdf.Cell(45, 7, "Income")
	pdf.Cell(45, 7, "Expenses")
	pdf.Cell(50, 7, "Savings")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 10)
	for _, month := range report.MonthlyBreakdown {
		pdf.Cell(40, 6, time.Month(month.Month).String())
		pdf.Cell(45, 6, fmt.Sprintf("%.2f", month.Income))
		pdf.Cell(45, 6, fmt.Sprintf("%.2f", month.Expenses))
		pdf.Cell(50, 6, fmt.Sprintf("%.2f", month.Savings))
		pdf.Ln(6)
	}
	pdf.Ln(5)

	// Category Summary
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Category Yearly Summary", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 7, "Category")
	pdf.Cell(45, 7, "Total Spent")
	pdf.Cell(45, 7, "Avg Monthly")
	pdf.Cell(40, 7, "Peak Month")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 10)
	for _, cat := range report.CategoryYearlySummary {
		pdf.Cell(50, 6, cat.Category)
		pdf.Cell(45, 6, fmt.Sprintf("%.2f", cat.TotalSpent))
		pdf.Cell(45, 6, fmt.Sprintf("%.2f", cat.AverageMonthly))
		pdf.Cell(40, 6, time.Month(cat.HighestMonth).String())
		pdf.Ln(6)
	}
	pdf.Ln(10)

	// Trends
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 10, "Trends & Insights", "1", 0, "L", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Average Monthly Income: Rp %.2f", report.Trends.AverageMonthlyIncome))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Average Monthly Expense: Rp %.2f", report.Trends.AverageMonthlyExpense))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Highest Income Month: %s", time.Month(report.Trends.HighestIncomeMonth).String()))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Highest Expense Month: %s", time.Month(report.Trends.HighestExpenseMonth).String()))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Total Savings Rate: %.2f%%", report.Trends.TotalSavingsRate))
	pdf.Ln(10)

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return &buf, nil
}
