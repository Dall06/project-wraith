package guard

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"project-wraith/pkg/modules/tools"
)

type Credentials struct {
	Username string
	Password string
}

type Manticore interface {
	StingAndProwl(Credentials) error
}

type manticore struct {
	collection *mongo.Collection
	ctx        context.Context
	passSecret string
}

func NewManticore(collection mongo.Collection, ctx context.Context, passSecret string) Manticore {
	return &manticore{
		collection: &collection,
		ctx:        ctx,
		passSecret: passSecret,
	}
}

func (m *manticore) StingAndProwl(cred Credentials) error {
	filter := bson.M{"username": cred.Username}

	var result Credentials
	err := m.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}

		return err
	}

	if result.Password != tools.Sha512(m.passSecret, cred.Password) {
		return errors.New("password incorrect")
	}

	return nil
}
