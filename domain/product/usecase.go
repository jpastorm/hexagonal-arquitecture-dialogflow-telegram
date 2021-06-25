package product

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jpastorm/dialogflowbot/model"
)

var allowedFieldsForQuery = []string{
	"id",
	"name",
	"price",
}

// Product implement the UseCase interface
type Product struct {
	storage Storage
}

// New returns a new Product pointer
func New(s Storage) Product {
	return Product{storage: s}
}

// Create creates a new model.Product
func (c Product) Create(m *model.Product) error {
	if err := model.ValidateStructNil(m); err != nil {
		return fmt.Errorf("Product: %w", err)
	}

	if err := c.storage.Create(m); err != nil {
		return fmt.Errorf("Product: %w", err)
	}

	return nil
}

// Update this method updates a model.Product by ID
func (c Product) Update(m *model.Product) error {
	if err := model.ValidateStructNil(m); err != nil {
		return fmt.Errorf("Product: %w", err)
	}

	if !m.HasID() {
		return model.ErrInvalidID
	}

	if err := c.storage.Update(m); err != nil {
		return fmt.Errorf("Product: could not update the record with id %d, %w", m.ID, err)
	}

	return nil
}

// Delete this method deletes a model.Product by ID
func (c Product) Delete(ID uint) error {
	if err := c.storage.Delete(ID); err != nil {
		return fmt.Errorf("Product: could not delete the record %d, %w", ID, err)
	}

	return nil
}

// GetWhere get a model.Product by the conditions of the  fields
func (c Product) GetWhere(filter model.Fields, sort model.SortFields) (model.Product, error) {
	if err := filter.ValidateNames(allowedFieldsForQuery); err != nil {
		return model.Product{}, fmt.Errorf("Product: %w", err)
	}

	if err := sort.ValidateNames(allowedFieldsForQuery); err != nil {
		return model.Product{}, fmt.Errorf("Product: %w", err)
	}

	m, err := c.storage.GetWhere(filter, sort)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Product{}, fmt.Errorf(" Product: %s", filter.Error())
	}

	if err != nil {
		return model.Product{}, fmt.Errorf("Product: %w", err)
	}

	return m, nil
}

// GetAllWhere get a model.Product by the conditions of the fields
func (c Product) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Products, error) {
	if err := filter.ValidateNames(allowedFieldsForQuery); err != nil {
		return nil, fmt.Errorf("Product: %w", err)
	}

	if err := sort.ValidateNames(allowedFieldsForQuery); err != nil {
		return nil, fmt.Errorf("Product: %w", err)
	}

	ms, err := c.storage.GetAllWhere(filter, sort, pag)
	if err != nil {
		return nil, fmt.Errorf("Product: %w", err)
	}

	return ms, nil
}
