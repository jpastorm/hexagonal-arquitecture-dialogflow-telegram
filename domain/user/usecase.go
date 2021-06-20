package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jpastorm/dialogflowbot/model"
)

var allowedFieldsForQuery = []string{
	"id",
	"name",
	"telegram_id",
}


// User implement the UseCase interface
type User struct {
	storage Storage
}

// New returns a new User pointer
func New(s Storage) User {
	return User{storage: s}
}

// Create creates a new model.User
func (c User) Create(m *model.User) error {
	if err := model.ValidateStructNil(m); err != nil {
		return fmt.Errorf("User: %w", err)
	}

	if err := c.storage.Create(m); err != nil {
		return fmt.Errorf("User: %w", err)
	}

	return nil
}

// Update this method updates a model.User by ID
func (c User) Update(m *model.User) error {
	if err := model.ValidateStructNil(m); err != nil {
		return fmt.Errorf("User: %w", err)
	}

	if !m.HasID() {
		return model.ErrInvalidID
	}

	if err := c.storage.Update(m); err != nil {
		return fmt.Errorf("User: could not update the record with id %d, %w", m.ID, err)
	}

	return nil
}

// Delete this method deletes a model.User by ID
func (c User) Delete(ID uint) error {
	if err := c.storage.Delete(ID); err != nil {
		return fmt.Errorf("User: could not delete the record %d, %w", ID, err)
	}

	return nil
}

// GetWhere get a model.User by the conditions of the  fields
func (c User) GetWhere(filter model.Fields, sort model.SortFields) (model.User, error) {
	if err := filter.ValidateNames(allowedFieldsForQuery); err != nil {
		return model.User{}, fmt.Errorf("User: %w", err)
	}

	if err := sort.ValidateNames(allowedFieldsForQuery); err != nil {
		return model.User{}, fmt.Errorf("User: %w", err)
	}

	m, err := c.storage.GetWhere(filter, sort)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf(" User: %s", filter.Error())
	}

	if err != nil {
		return model.User{}, fmt.Errorf("User: %w", err)
	}

	return m, nil
}

// GetAllWhere get a model.User by the conditions of the fields
func (c User) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error) {
	if err := filter.ValidateNames(allowedFieldsForQuery); err != nil {
		return nil, fmt.Errorf("User: %w", err)
	}

	if err := sort.ValidateNames(allowedFieldsForQuery); err != nil {
		return nil, fmt.Errorf("User: %w", err)
	}

	ms, err := c.storage.GetAllWhere(filter, sort, pag)
	if err != nil {
		return nil, fmt.Errorf("User: %w", err)
	}

	return ms, nil
}

