package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidPaginationParameter = errors.New("invalid pagination parameters")
)

type operatorField string

// Operators
var (
	Equals               operatorField = "="
	NotEqualTo           operatorField = "<>"
	LessThan             operatorField = "<"
	GreaterThan          operatorField = ">"
	LessThanOrEqualTo    operatorField = "<="
	GreaterThanOrEqualTo operatorField = ">="
	Ilike                operatorField = "ILIKE"
	In                   operatorField = "IN"
	IsNull               operatorField = "IS NULL"
	IsNotNull            operatorField = "IS NOT NULL"
)

// ChainingField is the keyword for chaining the next field
type ChainingField string

// Chaining
var (
	And ChainingField = "AND"
	Or  ChainingField = "OR"
)

type CustomFields struct {
	Filters    Fields
	Sorts      SortFields
	Pagination Pagination
}

// Field contains the information of a field for a query
type Field struct {
	Name        string        `json:"name"`
	Operator    operatorField `json:"operator"`
	Value       interface{}   `json:"value"`
	ChainingKey ChainingField `json:"chaining_key"`
	// Source sets the origin of the field, is used if a resource has more of one source,
	// this is useful generally when an infrastructure implementation used "Joins"
	Source string `json:"source"` // Optional
	// GroupOpen allows beginning a conditions group of fields and the infrastructure
	// implementation must include into group the field that sets the GroupOpen = true
	GroupOpen bool `json:"group_open"` // Optional
	// GroupClose allows ending a conditions group of fields, and the infrastructure
	// implementation must include into group the field that sets the GroupClose = true
	GroupClose bool `json:"group_close"` // Optional
}

// Fields slice of Field
type Fields []Field

// IsEmpty returns if the Fields is empty
func (fs Fields) IsEmpty() bool { return len(fs) == 0 }

// ValidateNames valida if the fields is allowed for query
func (fs Fields) ValidateNames(allowedFields []string) error {
	for _, field := range fs {
		var isAllowed bool
		for _, allowedField := range allowedFields {
			if strings.ToLower(allowedField) == strings.ToLower(field.Name) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return fmt.Errorf("the field %s is not allowed for query", field.Name)
		}
	}
	return nil
}

// SetOrAddFilter set or add filter in model.Fields
func (fs *Fields) SetOrAddFilter(name string, value interface{}) {
	for i, filter := range *fs {
		if filter.Name == name {
			(*fs)[i].Value = value
			return
		}
	}

	*fs = append(*fs, Field{Name: name, Value: value})
}

// FindField returns the Field and it returns true if field was found
func (fs Fields) FindField(inputField string) (Field, bool) {
	for _, field := range fs {
		if strings.ToLower(field.Name) == strings.ToLower(inputField) {
			return field, true
		}
	}
	return Field{}, false
}

// Error returns a error string with the field name and value
func (fs Fields) Error() string {
	err := "not found"
	if len(fs) == 0 {
		return err
	}
	for _, field := range fs {
		err = fmt.Sprintf("%s: %v, %s", field.Name, field.Value, err)
	}

	return err
}

// OrderField is the keyword for order the field
type OrderField string

// OrderFields
var (
	Asc  OrderField = "ASC"
	Desc OrderField = "DESC"
)

// SortField contains the information of the order of a field
type SortField struct {
	Name  string     `json:"name"`
	Order OrderField `json:"order"`
}

// SortFields slice of SortField
type SortFields []SortField

// IsEmpty returns if the SortFields is empty
func (ss SortFields) IsEmpty() bool { return len(ss) == 0 }

// ValidateNames valida if the fields is allowed for ordering
func (ss SortFields) ValidateNames(allowedFields []string) error {
	for _, field := range ss {
		var isAllowed bool
		for _, allowedField := range allowedFields {
			if strings.ToLower(allowedField) == strings.ToLower(field.Name) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return fmt.Errorf("the field %s is not allowed for ordering", field.Name)
		}
	}
	return nil
}

// Pagination contains the information of the pagination
type Pagination struct {
	Page     uint `json:"page"`
	Limit    uint `json:"limit"`
	MaxLimit uint
}
