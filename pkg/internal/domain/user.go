package domain

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Get(user User) (*User, error)
	Create(user User) error
	Update(user User) error
	Delete(id string) error
	Duplicated(user User) ([]User, error)
}

type userRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserRepository(collection mongo.Collection, ctx context.Context) UserRepository {
	return &userRepository{
		collection: &collection,
		ctx:        ctx,
	}
}

func (r *userRepository) Get(user User) (*User, error) {
	filter := bson.M{}

	if user.ID != "" {
		filter["_id"] = user.ID
	}
	if user.Username != "" {
		filter["username"] = user.Username
	}
	if user.Email != "" {
		filter["email"] = user.Email
	}
	if user.Phone != "" {
		filter["phone"] = user.Phone
	}

	err := r.collection.FindOne(r.ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found (%e)", err)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(user User) error {
	_, err := r.collection.InsertOne(r.ctx, user)
	return err
}

func (r *userRepository) Update(user User) error {
	if user.ID == "" {
		return errors.New("user ID is required")
	}

	// Create a filter to match the user by ID
	filter := bson.M{"_id": user.ID}

	// Create an update document and add only non-empty fields
	toUpdate := bson.M{}

	if user.Name != "" {
		toUpdate["name"] = user.Name
	}
	if user.Email != "" {
		toUpdate["email"] = user.Email
	}
	if user.Phone != "" {
		toUpdate["phone"] = user.Phone
	}
	if user.Username != "" {
		toUpdate["username"] = user.Username
	}
	if user.Meta != nil {
		toUpdate["meta"] = user.Meta
	}
	if !user.UpdatedAt.IsZero() {
		toUpdate["updatedAt"] = user.UpdatedAt
	}
	if user.Status != "" {
		toUpdate["status"] = user.Status
	}
	if user.Password != "" {
		toUpdate["password"] = user.Password
	}

	if len(toUpdate) == 0 {
		return nil
	}

	update := bson.M{
		"$set": toUpdate,
	}

	// Perform the update operation
	_, err := r.collection.UpdateOne(
		r.ctx,
		filter,
		update,
		options.Update().SetUpsert(false),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) Delete(id string) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}
	return nil
}

func (r *userRepository) Duplicated(user User) ([]User, error) {
	filter := bson.M{}

	// Add filters based on provided user details
	if user.Username != "" {
		filter["username"] = user.Username
	}
	if user.Email != "" {
		filter["email"] = user.Email
	}
	if user.Phone != "" {
		filter["phone"] = user.Phone
	}

	// Use Find to retrieve potential duplicates
	cursor, err := r.collection.Find(r.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find duplicates: %w", err)
	}

	var duplicates []User

	// Decode the results into the duplicates slice
	for cursor.Next(r.ctx) {
		var existingUser User
		if err := cursor.Decode(&existingUser); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		duplicates = append(duplicates, existingUser)
	}

	// Check for any cursor error
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// Return an error if duplicates are found
	if len(duplicates) > 0 {
		return duplicates, nil
	}

	err = cursor.Close(r.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to close cursor: %w", err)
	}

	return nil, nil // No duplicates found
}
