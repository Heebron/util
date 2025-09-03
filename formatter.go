// Package util provides utility functions and types for common operations.
package util

import "github.com/leekchan/accounting"

// Pre-configured accounting formatters for various currency and number formats.
// These can be used directly to format numbers in different display styles.
var (
	// DollarsFormat formats numbers as dollars with 2 decimal places (e.g., $123.45)
	DollarsFormat = accounting.Accounting{Symbol: "$", Precision: 2}

	// CentsFormat formats numbers as cents with 2 decimal places (e.g., ¢123.45)
	CentsFormat = accounting.Accounting{Symbol: "¢", Precision: 2}

	// DollarsCentsFormat formats numbers as dollars with 4 decimal places for precise values (e.g., $123.4567)
	DollarsCentsFormat = accounting.Accounting{Symbol: "$", Precision: 4}

	// DollarsCentsFormat2 formats numbers as dollars with 2 decimal places (same as DollarsFormat)
	DollarsCentsFormat2 = accounting.Accounting{Symbol: "$", Precision: 2}

	// PercentFormat formats numbers as percentages with 2 decimal places (e.g., 12.34%)
	PercentFormat = accounting.Accounting{Symbol: "%", Precision: 2}

	// NumberFormat formats numbers with thousands separators and no decimal places (e.g., 1,234,567)
	NumberFormat = accounting.Accounting{Symbol: "", Precision: 0, Thousand: ",", Decimal: ".", Format: "", FormatNegative: "", FormatZero: ""}
)
