package rules

import (
	"errors"
	"project-wraith/src/modules/domain"
	"project-wraith/src/pkg/jwt"
	"time"
)

type ResetRule interface {
	Start(reset Reset) (*Reset, error)
	Validate(reset Reset) error
}

type resetRule struct {
	repo      *domain.UserRepository
	jwtSecret string
}

func NewResetRule(repository *domain.UserRepository, jwtSecret string) ResetRule {
	return &resetRule{
		repo:      repository,
		jwtSecret: jwtSecret,
	}
}

func (rr resetRule) Start(reset Reset) (*Reset, error) {
	entity := domain.User{
		Email: reset.Email,
		Phone: reset.Phone,
	}

	response, err := rr.repo.Get(entity)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("user not found")
	}

	result := &Reset{
		ID:       response.ID,
		Username: response.Username,
		Email:    response.Email,
		Phone:    response.Phone,
	}

	claims := map[string]interface{}{
		"id":    result.ID,
		"email": result.Email,
		"phone": result.Phone,
		"reset": true,
	}

	tkn, err := jwt.CreateJwtToken(rr.jwtSecret, 10*time.Minute, claims)
	if err != nil {
		return nil, err
	}

	result.Token = tkn

	return result, nil
}

func (rr resetRule) Validate(reset Reset) error {
	valid, err := jwt.ValidateJwtToken(reset.Token, rr.jwtSecret, nil)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid token")
	}

	return nil
}
