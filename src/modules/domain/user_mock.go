package domain

import (
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Get(user User) (*User, error) {
	args := m.Called(user)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Create(user User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) Update(user User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	return m.Called(id).Error(0)
}
