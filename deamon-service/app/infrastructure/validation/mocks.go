package validation

import (
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/delfare/deamon-service/domain"
)

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(event domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
