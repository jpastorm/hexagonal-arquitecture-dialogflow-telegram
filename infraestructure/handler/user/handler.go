package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/jpastorm/dialogflowbot/domain/user"
	"github.com/jpastorm/dialogflowbot/infraestructure/request"
	"github.com/jpastorm/dialogflowbot/infraestructure/response"
	"github.com/jpastorm/dialogflowbot/model"
)

// User handler
type User struct {
	usecase user.UseCase
}

// NewHandler returns a User handler
func NewHandler(usecase user.UseCase) User {
	return User{usecase: usecase}
}

// Create handle the creation of a model.User
func (e User) Create(c echo.Context) error {
	data := model.User{}

	if err := c.Bind(&data); err != nil {
		return response.Failed("c.Bind()", response.BindFailed, err)
	}

	errData := model.NewError()
	if err := e.usecase.Create(&data); err != nil {
		if errors.As(err, &errData) {
			errResponse := response.Failed("useCase.Create()", response.Failure, err)
			errResponse.SetStatus(http.StatusBadRequest)
			errResponse.SetAPIMessage(errData.APIMessage())
			return errResponse
		}

		return response.Failed("useCase.Create()", response.UnexpectedError, err)
	}

	return c.JSON(response.Successfull(response.RecordCreated, data))
}

// Update handle the updating of a model.User
func (e User) Update(c echo.Context) error {
	data := model.User{}

	ID, err := request.ExtractIDFromURLParam(c)
	if err != nil {
		return response.Failed("request.ExtractIDFromURLParam()", response.InvalidParameter, err)
	}
	data.ID = uint(ID)

	if err := c.Bind(&data); err != nil {
		return response.Failed("c.Bind()", response.BindFailed, err)
	}

	errData := model.NewError()
	err = e.usecase.Update(&data)
	if errors.Is(err, model.ErrInvalidID) {
		return response.Failed("usecase.Update()", response.InvalidParameter, err)
	}
	if err != nil {
		if errors.As(err, &errData) {
			errResponse := response.Failed("useCase.Update()", response.Failure, err)
			errResponse.SetStatus(http.StatusBadRequest)
			errResponse.SetAPIMessage(errData.APIMessage())
			return errResponse
		}

		return response.Failed("useCase.Update()", response.UnexpectedError, err)
	}

	return c.JSON(response.Successfull(response.RecordUpdated, data))
}

// Delete handle the deleting of a model.User
func (e User) Delete(c echo.Context) error {
	ID, err := request.ExtractIDFromURLParam(c)
	if err != nil {
		return response.Failed("request.ExtractIDFromURLParam()", response.InvalidParameter, err)
	}

	if err := e.usecase.Delete(uint(ID)); err != nil {
		return response.Failed("usecase.Delete()", response.UnexpectedError, err)
	}

	return c.JSON(response.Successfull(response.RecordDeleted, nil))
}

// GetWhere this method get an ordered model.User
func (e User) GetWhere(c echo.Context) error {
	filters := c.QueryParam("filters")
	fields := model.Fields{}
	if filters != "" {
		err := json.Unmarshal([]byte(filters), &fields)
		if err != nil {
			return response.Failed("GetWhere().json.Unmarshal()", response.UnexpectedError, err)
		}
	}

	sorts := c.QueryParam("sorts")
	sortsFields := model.SortFields{}
	if sorts != "" {
		err := json.Unmarshal([]byte(sorts), &sortsFields)
		if err != nil {
			return response.Failed("GetWhere().json.Unmarshal()", response.UnexpectedError, err)
		}
	}

	data, err := e.usecase.GetWhere(fields, sortsFields)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(response.Successfull(response.RecordNotFound, nil))
	}
	if err != nil {
		return response.Failed("GetWhere().json.Unmarshal()", response.UnexpectedError, err)
	}

	return c.JSON(response.Successfull(response.Ok, data))
}

// GetAllWhere this method get all ordered model.EndpointJob
func (e User) GetAllWhere(c echo.Context) error {
	filters := c.QueryParam("filters")
	fields := model.Fields{}
	if filters != "" {
		err := json.Unmarshal([]byte(filters), &fields)
		if err != nil {
			return response.Failed("GetWhere().json.Unmarshal()", response.UnexpectedError, err)
		}
	}

	sorts := c.QueryParam("sorts")
	sortsFields := model.SortFields{}
	if sorts != "" {
		err := json.Unmarshal([]byte(sorts), &sortsFields)
		if err != nil {
			return response.Failed("GetWhere().json.Unmarshal()", response.UnexpectedError, err)
		}
	}

	pagination := c.QueryParam("pagination")
	pag := model.Pagination{}
	if pagination != "" {
		err := json.Unmarshal([]byte(pagination), &pag)
		if err != nil {
			return response.Failed("GetWhere().json.Unmarshal()", response.UnexpectedError, err)
		}
	}

	data, err := e.usecase.GetAllWhere(fields, sortsFields, pag)
	if err != nil {
		return response.Failed("GetWhere().json.Unmarshal()", response.UnexpectedError, err)
	}

	return c.JSON(response.Successfull(response.Ok, data))
}
