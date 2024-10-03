package rules

import "github.com/stretchr/testify/mock"

type MockResetRule struct {
	mock.Mock
}

func (m *MockResetRule) Validate(reset Reset) (*Reset, error) {
	args := m.Called(reset)
	if args.Get(0) != nil {
		return args.Get(0).(*Reset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockResetRule) Start(reset Reset) (*Reset, error) {
	args := m.Called(reset)
	if args.Get(0) != nil {
		return args.Get(0).(*Reset), args.Error(1)
	}
	return nil, args.Error(1)
}
