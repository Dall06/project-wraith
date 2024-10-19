package rules

import (
	"github.com/stretchr/testify/mock"
)

type MockUserRule struct {
	mock.Mock
}

func (m *MockUserRule) Login(model User) (*User, error) {
	args := m.Called(model)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRule) Register(model User) (*User, error) {
	args := m.Called(model)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRule) Edit(model User) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockUserRule) Get(model User) (*User, error) {
	args := m.Called(model)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRule) Disable(model User) error {
	args := m.Called(model)
	return args.Error(0)
}
