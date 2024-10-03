package core_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"project-wraith/pkg/core" // Update import according to your project's structure
	"project-wraith/pkg/modules/lics"
)

func TestActivate(t *testing.T) {
	t.Parallel()
	// Define test cases
	tests := []struct {
		name        string
		licenseKey  string
		mockLicense *lics.License
		expectedErr error
	}{
		{
			name:       "Activate a valid license",
			licenseKey: "valid-key",
			mockLicense: &lics.License{
				LicenseKey:         "valid-key",
				IsActive:           true,
				ExpiryDate:         time.Now().Add(24 * time.Hour), // Active and not expired
				IssuedAt:           time.Now(),
				MaxActivations:     5,
				CurrentActivations: 1,
			},
			expectedErr: nil,
		},
		{
			name:       "Activate an inactive license",
			licenseKey: "inactive-key",
			mockLicense: &lics.License{
				LicenseKey:         "inactive-key",
				IsActive:           false,
				ExpiryDate:         time.Now().Add(24 * time.Hour), // Inactive but not expired
				IssuedAt:           time.Now(),
				MaxActivations:     5,
				CurrentActivations: 1,
			},
			expectedErr: errors.New("license is either inactive or expired"),
		},
		{
			name:       "Activate an expired license",
			licenseKey: "expired-key",
			mockLicense: &lics.License{
				LicenseKey:         "expired-key",
				IsActive:           true,
				ExpiryDate:         time.Now().Add(-24 * time.Hour), // Active but expired
				IssuedAt:           time.Now(),
				MaxActivations:     5,
				CurrentActivations: 1,
			},
			expectedErr: errors.New("license is either inactive or expired"),
		},
		{
			name:        "License not found",
			licenseKey:  "not-found-key",
			mockLicense: nil, // No license found
			expectedErr: errors.New("license not found"),
		},
	}

	// Execute test cases
	for _, tc := range tests {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := new(lics.MockLicenseRepository)

			// Setup the expected calls and return values
			if tc.mockLicense != nil {
				repo.On("Get", mock.Anything).Return(tc.mockLicense, nil)
			} else {
				// Return nil license and an error for the "not found" scenario
				repo.On("Get", mock.Anything).Return((*lics.License)(nil), errors.New("license not found"))
			}

			err := core.Activate(repo, tc.licenseKey)

			// Check the error against the expected error
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t) // Assert that expectations were met
		})
	}
}
