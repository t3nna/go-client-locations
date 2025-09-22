package util

import (
	"testing"
)

func TestValidateUserName(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid username - 4 characters",
			username:    "user",
			expectError: false,
		},
		{
			name:        "valid username - 16 characters",
			username:    "user123456789012",
			expectError: false,
		},
		{
			name:        "valid username - alphanumeric",
			username:    "user123",
			expectError: false,
		},
		{
			name:        "invalid username - too short",
			username:    "ab",
			expectError: true,
			errorMsg:    "username must be between 4 and 16 characters",
		},
		{
			name:        "invalid username - too long",
			username:    "user1234567890123",
			expectError: true,
			errorMsg:    "username must be between 4 and 16 characters",
		},
		{
			name:        "invalid username - contains special characters",
			username:    "user@123",
			expectError: true,
			errorMsg:    "username can only contain alphanumeric characters",
		},
		{
			name:        "invalid username - contains spaces",
			username:    "user 123",
			expectError: true,
			errorMsg:    "username can only contain alphanumeric characters",
		},
		{
			name:        "invalid username - empty string",
			username:    "",
			expectError: true,
			errorMsg:    "username must be between 4 and 16 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserName(tt.username)

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
