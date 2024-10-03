package guard_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"project-wraith/pkg/modules/guard"
	"project-wraith/pkg/modules/tools"
	"testing"
)

func TestManticore(test *testing.T) {
	test.Parallel()

	mt := mtest.New(test, mtest.NewOptions().ClientType(mtest.Mock))

	testCases := []struct {
		name        string
		credentials guard.Credentials
		expectedErr error
		setupMocks  func(mongoTest *mtest.T)
	}{
		{
			name: "successful sting and prowl",
			credentials: guard.Credentials{
				Username: "user1",
				Password: "password123",
			},
			expectedErr: nil,
			setupMocks: func(mongoTest *mtest.T) {
				mongoTest.AddMockResponses(mtest.CreateCursorResponse(1, "db.users", mtest.FirstBatch, bson.D{
					{"username", "user1"},
					{"password", tools.Sha512("secret", "password123")}, // Adjust to match your hashing method
				}))
			},
		},
		{
			name: "incorrect password",
			credentials: guard.Credentials{
				Username: "user1",
				Password: "wrongpassword",
			},
			expectedErr: errors.New("password incorrect"),
			setupMocks: func(mongoTest *mtest.T) {
				mongoTest.AddMockResponses(mtest.CreateCursorResponse(1, "db.users", mtest.FirstBatch, bson.D{
					{"username", "user1"},
					{"password", tools.Sha512("secret", "password123")}, // Correct hashed password
				}))
			},
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mongoTest *mtest.T) {
			mongoTest.Parallel()
			tc.setupMocks(mongoTest)

			collection := mongoTest.Coll
			repo := guard.NewManticore(*collection, context.TODO(), "secret")

			err := repo.StingAndProwl(tc.credentials)
			if tc.expectedErr != nil {
				assert.EqualError(test, err, tc.expectedErr.Error())
			} else {
				assert.NoError(test, err)
			}
		})
	}
}
