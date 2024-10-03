package rules

import (
	"errors"
	"project-wraith/pkg/internal/domain"
	"project-wraith/pkg/modules/alchemy"
	"project-wraith/pkg/modules/tools"
	"time"
)

type UserRule interface {
	Login(model User) (*User, error)
	Register(model User) (*User, error)
	Edit(model User) error
	Get(model User) (*User, error)
	Remove(model User) error
}

type userRule struct {
	repo          domain.UserRepository
	encryptDbData bool
	dbDataSecret  string
	passSecret    string
}

func NewUserRule(
	repo domain.UserRepository,
	encryptDbData bool,
	dbDataSecret string,
	passSecret string) UserRule {
	return &userRule{
		repo:          repo,
		encryptDbData: encryptDbData,
		dbDataSecret:  dbDataSecret,
		passSecret:    passSecret,
	}
}

func (r userRule) Login(model User) (*User, error) {
	entity := domain.User{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Name:     model.Name,
		Phone:    model.Phone,
		Password: model.Password,
	}

	if r.encryptDbData {
		err := alchemy.Transmutation(&entity, r.dbDataSecret)
		if err != nil {
			return nil, err
		}
	}

	response, err := r.repo.Get(entity)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("user not found")
	}

	if r.encryptDbData {
		err = alchemy.Revert(&response, r.dbDataSecret)
		if err != nil {
			return nil, err
		}
	}

	if response.Password != tools.Sha512(r.passSecret, entity.Password) {
		return nil, errors.New("password incorrect")
	}

	result := &User{
		ID:       response.ID,
		Username: response.Username,
		Email:    response.Email,
		Name:     response.Name,
		Phone:    response.Phone,
		Password: response.Password,
	}

	return result, nil
}

func (r userRule) Register(model User) (*User, error) {
	entity := domain.User{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		Name:      model.Name,
		Phone:     model.Phone,
		Password:  tools.Sha512(r.passSecret, model.Password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	duplicates, err := r.repo.Duplicated(entity)
	if err != nil {
		return nil, err
	}

	if len(duplicates) > 0 {
		return nil, errors.New("user already exists")
	}

	if r.encryptDbData {
		err := alchemy.Transmutation(&entity, r.dbDataSecret)
		if err != nil {
			return nil, err
		}
	}

	err = r.repo.Create(entity)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r userRule) Edit(model User) error {
	entity := domain.User{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		Name:      model.Name,
		Phone:     model.Phone,
		Password:  tools.Sha512(r.passSecret, model.Password),
		UpdatedAt: time.Now(),
	}

	if r.encryptDbData {
		err := alchemy.Transmutation(&entity, r.dbDataSecret)
		if err != nil {
			return err
		}
	}

	err := r.repo.Update(entity)
	if err != nil {
		return err
	}

	return nil
}

func (r userRule) Get(model User) (*User, error) {
	entity := domain.User{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Phone:    model.Phone,
	}

	if r.encryptDbData {
		err := alchemy.Transmutation(&entity, r.dbDataSecret)
		if err != nil {
			return nil, err
		}
	}

	response, err := r.repo.Get(entity)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("user not found")
	}

	result := &User{
		ID:       response.ID,
		Username: response.Username,
		Email:    response.Email,
		Name:     response.Name,
		Phone:    response.Phone,
		Password: response.Password,
	}

	if r.encryptDbData {
		err = alchemy.Revert(&result, r.dbDataSecret)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r userRule) Remove(model User) error {
	entity := domain.User{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Phone:    model.Phone,
		Password: model.Password,
	}

	if r.encryptDbData {
		err := alchemy.Transmutation(&entity, r.dbDataSecret)
		if err != nil {
			return err
		}
	}

	response, err := r.repo.Get(entity)
	if err != nil {
		return nil
	}

	if response == nil {
		return errors.New("user not found")
	}

	if r.encryptDbData {
		err = alchemy.Revert(&response, r.dbDataSecret)
		if err != nil {
			return err
		}
	}

	if response.Password != tools.Sha512(r.passSecret, entity.Password) {
		return errors.New("password incorrect")
	}

	err = r.repo.Delete(entity.ID)
	if err != nil {
		return err
	}

	return nil
}
