package validation

import (
	"errors"

	"github.com/yakob-abada/delfare/deamon-service/domain"
)

type Validator interface {
	Validate(event domain.Event) error
}

type SecurityEventValidator struct{}

func NewSecurityEventValidator() *SecurityEventValidator {
	return &SecurityEventValidator{}
}

func (v *SecurityEventValidator) Validate(event domain.Event) error {
	// Validate Criticality
	if event.Criticality < 0 || event.Criticality > 10 {
		return errors.New("invalid criticality value: must be between 0 and 10")
	}

	// Validate Message (optional)
	if event.Message == "" {
		return errors.New("event message cannot be empty")
	}

	return nil
}
