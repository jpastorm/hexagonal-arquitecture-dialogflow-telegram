package model

import (
	"errors"
	"reflect"
)

// Errors
var (
	ErrNilPointer        = errors.New("el modelo recibido es nulo")
	ErrInvalidID         = errors.New("el id recibido no es valido")
	ErrModelNotFound     = errors.New("el modelo no fue encontrado")
	ErrParameterNotFound = errors.New("el parámetro no se encuentra configurado")
)

// Errors SQL
var (
	ErrUnique     = errors.New("Unique violation")
	ErrForeignKey = errors.New("Foreign key violation")
	ErrNotNull    = errors.New("Not Null violation")
)

// ValidateStructNil returns an error if the model is nil
func ValidateStructNil(i interface{}) error {
	// omit struct type
	if reflect.ValueOf(i).Kind() == reflect.Struct {
		return nil
	}

	// Type: nil, Value: nil
	if i == nil {
		return ErrNilPointer
	}

	// Type: StructPointer, Value: nil
	// example: Type: *Cashbox, Value: nil
	if reflect.ValueOf(i).IsNil() {
		return ErrNilPointer
	}

	// Type: StructPointer, Value: ZeroValue
	// example: Type: *CashBox, Value: &CashBox{}
	return nil
}
