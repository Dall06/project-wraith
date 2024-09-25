package lics

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

func TestLicenseRepository_Get(test *testing.T) {
	test.Parallel()

	mt := mtest.New(test, mtest.NewOptions().ClientType(mtest.Mock))

	testCases := []struct {
		name           string
		input          License
		mockResponse   bson.D
		expectedResult *License
		expectedErr    error
	}{
		{
			name: "Successfully retrieve a license",
			input: License{
				LicenseKey: "1234-5678-9101",
			},
			mockResponse: bson.D{
				{"license_key", "1234-5678-9101"},
				{"product_name", "Test License"},
				{"expiry_date", "2025-01-01"},
			},
			expectedResult: &License{
				LicenseKey:  "1234-5678-9101",
				ProductName: "Test License",
				ExpiryDate:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: nil,
		},
		{
			name: "License not found",
			input: License{
				LicenseKey: "not-found",
			},
			mockResponse:   nil, // No response for the 'not found' case
			expectedResult: nil,
			expectedErr:    errors.New("license not found"),
		},
		{
			name: "Database error",
			input: License{
				LicenseKey: "error-case",
			},
			mockResponse:   nil, // No response needed
			expectedResult: nil,
			expectedErr:    errors.New("some db error"), // This is just a placeholder
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mongoTest *mtest.T) {
			mongoTest.Parallel()

			collection := mongoTest.Coll
			repo := NewLicenseRepository(*collection, context.TODO())

			if tc.mockResponse != nil {
				mongoTest.AddMockResponses(mtest.CreateCursorResponse(1, "db.licenses", mtest.FirstBatch, tc.mockResponse))
			} else {
				mongoTest.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    11000,           // Code for error (could vary depending on case)
					Message: "some db error", // Message for the error case
				}))
			}

			result, err := repo.Get(tc.input)
			assert.Equal(test, tc.expectedErr, err)
			assert.Equal(test, tc.expectedResult, result)
		})
	}
}
