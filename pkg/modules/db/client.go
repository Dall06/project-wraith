package db

import (
	"context"
	"fmt"
	"time"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Open() error
	Close() error
	Collection(coll string) *mongo.Collection
	Client() *mongo.Client
	Ctx() context.Context
}

type client struct {
	client *mongo.Client
	Db     *mongo.Database
	ctx    context.Context
	uri    string
	dbName string
}

// NewClient creates a new instance of Client.
func NewClient(uri, dbName string) Client {
	return &client{uri: uri, dbName: dbName}
}

// Open connects to the MongoDB instance.
func (mc *client) Open() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mc.uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	mc.client = client
	mc.Db = client.Database(mc.dbName)
	mc.ctx = ctx

	return nil
}

// Close disconnects from the MongoDB instance.
func (mc *client) Close() error {
	if mc.client == nil {
		return fmt.Errorf("client not connected")
	}
	return mc.client.Disconnect(context.Background())
}

func (mc *client) Collection(coll string) *mongo.Collection {
	return mc.Db.Collection(coll)
}

func (mc *client) Client() *mongo.Client {
	return mc.client
}

func (mc *client) Ctx() context.Context {
	return mc.ctx
}
