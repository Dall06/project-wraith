package rules

import "github.com/stretchr/testify/mock"

type MockResetRepository struct {
	mock.Mock
}

func (m *MockResetRepository) Validate(reset Reset) error {
	args := m.Called(reset)
	return args.Error(0)
}

func (m *MockResetRepository) Start(reset Reset) error {
	args := m.Called(reset)
	return args.Error(0)
}
