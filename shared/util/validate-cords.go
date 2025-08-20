package util

import "errors"

func ValidateCords(lat float64, lon float64) error {
	// Check if latitude is within the valid range of -90 to +90.
	if lat < -90 || lat > 90 {
		return errors.New("latitude must be between -90 and 90")
	}

	// Check if longitude is within the valid range of -180 to +180.
	if lon < -180 || lon > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	return nil
}
