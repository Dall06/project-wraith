package db

import (
	"context"
	"fmt"
	"time"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Client *mongo.Client
	Db     *mongo.Database
	Ctx    context.Context
	uri    string
	dbName string
}

// NewClient creates a new instance of Client.
func NewClient(uri, dbName string) *Client {
	return &Client{uri: uri, dbName: dbName}
}

// Open connects to the MongoDB instance.
func (mc *Client) Open() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mc.uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	mc.Client = client
	mc.Db = client.Database(mc.dbName)
	mc.Ctx = ctx

	return nil
}

// Close disconnects from the MongoDB instance.
func (mc *Client) Close() error {
	if mc.Client == nil {
		return fmt.Errorf("client not connected")
	}
	return mc.Client.Disconnect(context.Background())
}

func (mc *Client) Collection(coll string) *mongo.Collection {
	return mc.Db.Collection(coll)
}
