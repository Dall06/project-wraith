package db

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockClient struct {
	mock.Mock
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (mc *MockClient) Open() error {
	args := mc.Called()
	return args.Error(0)
}

func (mc *MockClient) Close() error {
	args := mc.Called()
	return args.Error(0)
}

func (mc *MockClient) Collection(coll string) *mongo.Collection {
	args := mc.Called(coll)
	if col, ok := args.Get(0).(*mongo.Collection); ok {
		return col
	}
	return nil
}

func (mc *MockClient) Ctx() context.Context {
	args := mc.Called()
	if ctx, ok := args.Get(0).(context.Context); ok {
		return ctx
	}
	return nil
}

func (mc *MockClient) Client() *mongo.Client {
	args := mc.Called()
	return args.Get(0).(*mongo.Client)
}
