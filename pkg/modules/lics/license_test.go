package lics_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"project-wraith/pkg/modules/lics"
	"testing"
	"time"
)

func TestLicenseRepository(test *testing.T) {
	test.Parallel()

	mt := mtest.New(test, mtest.NewOptions().ClientType(mtest.Mock))

	testCases := []struct {
		name        string
		action      string
		license     lics.License
		query       lics.License
		licenseKey  string
		expectedErr error
	}{
		{
			name:       "get license",
			action:     "get",
			licenseKey: "1234-5678-9101",
			license: lics.License{
				LicenseKey:         "1234-5678-9101",
				IsActive:           true,
				ExpiryDate:         time.Now().Add(24 * time.Hour).UTC(), // Active and not expired
				IssuedAt:           time.Now().UTC(),                     // Set the IssuedAt for clarity
				MaxActivations:     5,
				CurrentActivations: 1,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mongoTest *mtest.T) {
			mongoTest.Parallel()

			collection := mongoTest.Coll
			repo := lics.NewLicenseRepository(*collection, context.TODO())

			switch tc.action {
			case "get":
				mongoTest.AddMockResponses(mtest.CreateCursorResponse(1, "db.users", mtest.FirstBatch, bson.D{
					{"_id", tc.license.ID},
					{"license_key", tc.license.LicenseKey},
					{"customer_name", tc.license.CustomerName},
					{"customer_email", tc.license.CustomerEmail},
					{"issued_at", tc.license.IssuedAt},
					{"expiry_date", tc.license.ExpiryDate},
					{"max_activations", tc.license.MaxActivations},
					{"current_activations", tc.license.CurrentActivations},
					{"is_active", tc.license.IsActive},
					{"product_name", tc.license.ProductName},
				}))
				result, err := repo.Get(tc.query)
				assert.Equal(test, tc.expectedErr, err)
				assert.IsType(test, &tc.license, result)
				assert.Equal(test, tc.license.ID, result.ID)
			}
		})
	}
}
