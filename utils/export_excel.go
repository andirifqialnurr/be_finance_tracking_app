package utils

import (
	"bytes"
	"finance-tracking-app/services"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportMonthlyReportToExcel generates an Excel report for a given monthly report
func ExportMonthlyReportToExcel(report *services.MonthlyReport) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Monthly Report"
	f.SetSheetName("Sheet1", sheetName)

	// Set title
	title := fmt.Sprintf("Monthly Financial Report - %s %d",
		time.Month(report.Month).String(), report.Year)
	f.SetCellValue(sheetName, "A1", title)

	// Merge cells for title
	f.MergeCell(sheetName, "A1", "D1")

	// Style for title
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	f.SetCellStyle(sheetName, "A1", "D1", titleStyle)

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})

	// Currency style
	currencyStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 7, // Currency format
	})

	row := 3

	// Income Section
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Income Summary")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Income:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Income.Total)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row += 2

	// Income sources table
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Source")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Amount")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), headerStyle)
	row++

	for _, source := range report.Income.Sources {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), source.Source)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), source.Amount)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
		row++
	}
	row += 2

	// Expenses Section
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Expense Summary")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Expenses:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Expenses.Total)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row += 2

	// Expenses by category table
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Category")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Amount")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Count")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), headerStyle)
	row++

	for _, cat := range report.Expenses.ByCategory {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), cat.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), cat.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), cat.Count)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
		row++
	}
	row += 2

	// Budget Section
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Budget Summary")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Allocated:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Budget.TotalAllocated)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Spent:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Budget.TotalSpent)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Remaining:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Budget.TotalRemaining)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row += 2

	// Budget categories table
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Category")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Allocated")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Spent")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "Remaining")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	for _, cat := range report.Budget.CategoriesSummary {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), cat.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), cat.Allocated)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), cat.Spent)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), cat.Remaining)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("D%d", row), currencyStyle)
		row++
	}
	row += 2

	// Financial Summary
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Financial Summary")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Savings:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Savings)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Unallocated:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Unallocated)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Top Spending Category:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("%s (Rp %.2f)", report.TopCategory, report.TopCategoryAmount))
	row += 2

	// Insights
	if len(report.Insights) > 0 {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Insights")
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
		row++

		for _, insight := range report.Insights {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), insight)
			row++
		}
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 25)
	f.SetColWidth(sheetName, "B", "D", 15)

	// Output to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate Excel: %w", err)
	}

	return &buf, nil
}

// ExportYearlyReportToExcel generates an Excel report for a given yearly report
func ExportYearlyReportToExcel(report *services.YearlyReport) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Yearly Report"
	f.SetSheetName("Sheet1", sheetName)

	// Set title
	title := fmt.Sprintf("Yearly Financial Report - %d", report.Year)
	f.SetCellValue(sheetName, "A1", title)
	f.MergeCell(sheetName, "A1", "D1")

	// Style for title
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, "A1", "D1", titleStyle)

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})

	// Currency style
	currencyStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 7})

	row := 3

	// Annual Summary
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Annual Summary")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Income:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.TotalIncome)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Expenses:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.TotalExpenses)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Savings:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.TotalSavings)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row += 2

	// Monthly Breakdown
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Monthly Breakdown")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Month")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Income")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Expenses")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "Savings")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	for _, month := range report.MonthlyBreakdown {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), time.Month(month.Month).String())
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), month.Income)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), month.Expenses)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), month.Savings)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("D%d", row), currencyStyle)
		row++
	}
	row += 2

	// Category Summary
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Category Yearly Summary")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Category")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Total Spent")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "Avg Monthly")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "Peak Month")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	for _, cat := range report.CategoryYearlySummary {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), cat.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), cat.TotalSpent)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), cat.AverageMonthly)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), time.Month(cat.HighestMonth).String())
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), currencyStyle)
		row++
	}
	row += 2

	// Trends
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Trends & Insights")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Average Monthly Income:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Trends.AverageMonthlyIncome)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Average Monthly Expense:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Trends.AverageMonthlyExpense)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), currencyStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Highest Income Month:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), time.Month(report.Trends.HighestIncomeMonth).String())
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Highest Expense Month:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), time.Month(report.Trends.HighestExpenseMonth).String())
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Total Savings Rate:")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("%.2f%%", report.Trends.TotalSavingsRate))
	row++

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 30)
	f.SetColWidth(sheetName, "B", "D", 18)

	// Output to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate Excel: %w", err)
	}

	return &buf, nil
}
