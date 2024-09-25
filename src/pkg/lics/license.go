package lics

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type License struct {
	ID                 string    `json:"id" bson:"_id,omitempty"`
	ProductName        string    `json:"product_name" bson:"product_name"`
	LicenseKey         string    `json:"license_key" bson:"license_key"`
	CustomerName       string    `json:"customer_name" bson:"customer_name"`
	CustomerEmail      string    `json:"customer_email" bson:"customer_email"`
	IssuedAt           time.Time `json:"issued_at" bson:"issued_at"`
	ExpiryDate         time.Time `json:"expiry_date" bson:"expiry_date"`
	MaxActivations     int       `json:"max_activations" bson:"max_activations"`
	CurrentActivations int       `json:"current_activations" bson:"current_activations"`
	IsActive           bool      `json:"is_active" bson:"is_active"`
}

type LicenseRepository interface {
	Get(lic License) (*License, error)
}

type licenseRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewLicenseRepository(collection mongo.Collection, ctx context.Context) LicenseRepository {
	return &licenseRepository{
		collection: &collection,
		ctx:        ctx,
	}
}

func (r *licenseRepository) Get(lic License) (*License, error) {
	filter := bson.M{"license_key": lic.LicenseKey}

	var result License
	err := r.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("license not found")
		}
		return nil, err
	}

	return &result, nil
}
