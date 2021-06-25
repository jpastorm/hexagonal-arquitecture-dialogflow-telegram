package product

import "github.com/jpastorm/dialogflowbot/model"

// UseCase interfaces of use cases of User
type UseCase interface {
	Create(m *model.Product) error
	Update(m *model.Product) error
	Delete(ID uint) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.Product, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Products, error)
}

// Storage interface for db engines
type Storage interface {
	Create(m *model.Product) error
	Update(m *model.Product) error
	Delete(ID uint) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.Product, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Products, error)
}
