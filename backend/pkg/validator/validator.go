package validator

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

func ValidatePhone(phone string) bool {
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	return phoneRegex.MatchString(cleaned)
}

func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

func ValidateMinLength(value string, min int) bool {
	return len(strings.TrimSpace(value)) >= min
}

func ValidateMaxLength(value string, max int) bool {
	return len(strings.TrimSpace(value)) <= max
}

func ValidateCoordinates(lat, lon float64) bool {
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

func ValidatePositive(value float64) bool {
	return value > 0
}
