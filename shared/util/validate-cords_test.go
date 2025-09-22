package util

import (
	"testing"
)

func TestValidateCords(t *testing.T) {
	tests := []struct {
		name        string
		latitude    float64
		longitude   float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid coordinates - normal values",
			latitude:    51.11822470,
			longitude:   16.99071172,
			expectError: false,
		},
		{
			name:        "valid coordinates - edge values",
			latitude:    90.0,
			longitude:   180.0,
			expectError: false,
		},
		{
			name:        "valid coordinates - negative edge values",
			latitude:    -90.0,
			longitude:   -180.0,
			expectError: false,
		},
		{
			name:        "valid coordinates - zero values",
			latitude:    0.0,
			longitude:   0.0,
			expectError: false,
		},
		{
			name:        "valid coordinates - 8 decimal places",
			latitude:    51.12345678,
			longitude:   16.12345678,
			expectError: false,
		},
		{
			name:        "invalid latitude - too high",
			latitude:    91.0,
			longitude:   16.0,
			expectError: true,
			errorMsg:    "latitude must be between -90 and 90",
		},
		{
			name:        "invalid latitude - too low",
			latitude:    -91.0,
			longitude:   16.0,
			expectError: true,
			errorMsg:    "latitude must be between -90 and 90",
		},
		{
			name:        "invalid longitude - too high",
			latitude:    51.0,
			longitude:   181.0,
			expectError: true,
			errorMsg:    "longitude must be between -180 and 180",
		},
		{
			name:        "invalid longitude - too low",
			latitude:    51.0,
			longitude:   -181.0,
			expectError: true,
			errorMsg:    "longitude must be between -180 and 180",
		},
		{
			name:        "invalid coordinates - too many decimal places",
			latitude:    51.123456789,
			longitude:   16.123456789,
			expectError: true,
			errorMsg:    "coordinates must have up to 8 decimal places",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCords(tt.latitude, tt.longitude)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
