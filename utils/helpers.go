package utils

import (
	"time"
)

// GetCurrentMonthYear returns current month and year
func GetCurrentMonthYear() (int, int) {
	now := time.Now()
	return int(now.Month()), now.Year()
}

// GetMonthStartEnd returns start and end time of a given month
func GetMonthStartEnd(month, year int) (time.Time, time.Time) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return start, end
}

// FormatCurrency formats a float64 to Indonesian Rupiah format
func FormatCurrency(amount float64) string {
	return formatNumber(int64(amount))
}

func formatNumber(n int64) string {
	if n < 0 {
		return "-" + formatNumber(-n)
	}

	if n < 1000 {
		return string(rune(n))
	}

	return formatNumber(n/1000) + "." + string(rune(n%1000))
}

// RoundToTwoDecimals rounds a float64 to 2 decimal places
func RoundToTwoDecimals(num float64) float64 {
	return float64(int(num*100)) / 100
}

// CalculatePercentage calculates percentage of part from total
func CalculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return RoundToTwoDecimals((part / total) * 100)
}
