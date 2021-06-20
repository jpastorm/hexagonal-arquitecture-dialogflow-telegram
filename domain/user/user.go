package user

import "github.com/jpastorm/dialogflowbot/model"

// UseCase interfaces of use cases of User
type UseCase interface {
	Create(m *model.User) error
	Update(m *model.User) error
	Delete(ID uint) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.User, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error)
}

// Storage interface for db engines
type Storage interface {
	Create(m *model.User) error
	Update(m *model.User) error
	Delete(ID uint) error
	GetWhere(filter model.Fields, sort model.SortFields) (model.User, error)
	GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error)
}
