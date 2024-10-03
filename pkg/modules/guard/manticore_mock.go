package guard

import (
	"github.com/stretchr/testify/mock"
)

type MockManticore struct {
	mock.Mock
}

func (m *MockManticore) StingAndProwl(url string) error {
	return m.Called().Error(0)
}
