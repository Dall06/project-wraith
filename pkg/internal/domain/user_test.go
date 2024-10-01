package domain_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"project-wraith/pkg/internal/domain"
	"testing"
)

func TestUserRepository(test *testing.T) {
	test.Parallel()

	mt := mtest.New(test, mtest.NewOptions().ClientType(mtest.Mock))

	testCases := []struct {
		name        string
		action      string
		user        domain.User
		query       domain.User
		expectedErr error
	}{
		{
			name:   "Create new user",
			action: "create",
			user: domain.User{
				ID:       "1",
				Username: "newuser",
				Email:    "newuser@example.com",
			},
			expectedErr: nil,
		},
		{
			name:   "Get user by ID",
			action: "get",
			user: domain.User{
				ID:       "2",
				Username: "getuser",
				Email:    "getuser@example.com",
			},
			query:       domain.User{ID: "2"},
			expectedErr: nil,
		},
		{
			name:   "Update user",
			action: "update",
			user: domain.User{
				ID:       "3",
				Username: "updateuser",
				Email:    "updateuser@example.com",
			},
			query: domain.User{
				ID:       "3",
				Username: "updatedusername",
			},
			expectedErr: nil,
		},
		{
			name:   "Delete user",
			action: "delete",
			user: domain.User{
				ID: "4",
			},
			query:       domain.User{ID: "4"},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mongoTest *mtest.T) {
			mongoTest.Parallel()

			collection := mongoTest.Coll
			repo := domain.NewUserRepository(*collection, context.TODO())

			switch tc.action {
			case "create":
				mongoTest.AddMockResponses(mtest.CreateSuccessResponse())
				err := repo.Create(tc.user)
				assert.Equal(test, tc.expectedErr, err)

			case "get":
				mongoTest.AddMockResponses(mtest.CreateCursorResponse(1, "db.users", mtest.FirstBatch, bson.D{
					{"_id", tc.user.ID},
					{"username", tc.user.Username},
					{"email", tc.user.Email},
				}))
				result, err := repo.Get(tc.query)
				assert.Equal(test, tc.expectedErr, err)
				assert.Equal(test, &tc.user, result)

			case "update":
				mongoTest.AddMockResponses(mtest.CreateSuccessResponse())
				err := repo.Update(tc.query)
				assert.Equal(test, tc.expectedErr, err)

			case "delete":
				mongoTest.AddMockResponses(mtest.CreateSuccessResponse(
					bson.E{Key: "n", Value: int32(1)},
				))
				err := repo.Delete(tc.query.ID)
				assert.Equal(test, tc.expectedErr, err)
			}
		})
	}
}
