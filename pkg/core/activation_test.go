package core_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"project-wraith/pkg/core"
	"project-wraith/pkg/modules/lics"
	"testing"
	"time"
)

func TestActivate(test *testing.T) {
	test.Parallel()
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
				IssuedAt:           time.Now(),                     // Set the IssuedAt for clarity
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
		tc := tc
		test.Run(tc.name, func(t *testing.T) {
			repo := new(lics.MockLicenseRepository)

			// Setup the expected calls and return values
			if tc.mockLicense != nil {
				repo.On("Get", *tc.mockLicense).Return(tc.mockLicense, nil)
			} else {
				// Return nil license and an error for the "not found" scenario
				repo.On("Get", mock.Anything).Return(nil, errors.New("license not found"))
			}

			err := core.Activate(repo, tc.licenseKey)

			if tc.expectedErr != nil {
				assert.NotNil(test, err)
			} else {
				assert.NoError(test, err)
			}

		})
	}
}
