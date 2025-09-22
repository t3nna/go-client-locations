package util

import (
	"errors"
	"regexp"
	"strconv"
)

func ValidateCords(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}

	if longitude < -180 || longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	latStr := regexp.MustCompile(`^-?\d+(\.\d{1,8})?$`)
	lonStr := regexp.MustCompile(`^-?\d+(\.\d{1,8})?$`)

	if !latStr.MatchString(formatFloat(latitude)) || !lonStr.MatchString(formatFloat(longitude)) {
		return errors.New("coordinates must have up to 8 decimal places")
	}

	return nil
}
func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
