package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatCurrency formats an amount with the given symbol
func FormatCurrency(amount float64, symbol string) string {
	return symbol + formatWithCommas(amount)
}

// ParseCurrency parses a currency string to float64
func ParseCurrency(s string) (float64, error) {
	// Remove common currency symbols and commas
	cleaned := strings.Map(func(r rune) rune {
		if r == ',' {
			return -1
		}
		return r
	}, s)
	cleaned = strings.TrimSpace(cleaned)

	// Remove leading non-numeric characters (currency symbols)
	for len(cleaned) > 0 && !isNumeric(rune(cleaned[0])) && cleaned[0] != '-' && cleaned[0] != '.' {
		cleaned = cleaned[1:]
	}

	return strconv.ParseFloat(cleaned, 64)
}

func formatWithCommas(amount float64) string {
	str := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := parts[1]

	var result strings.Builder
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}

	return result.String() + "." + decPart
}

func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}
