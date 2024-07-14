package argv

import (
	"fmt"
	"github.com/oddbit-project/blueprint/utils"
)

const (
	// global errors
	ErrEmptyArgs             = utils.Error("empty argument list")
	ErrInvalidDest           = utils.Error("dest must be a ptr")
	ErrInvalidDestType       = utils.Error("invalid argument type; dest must be a struct")
	ErrInvalidParameterCount = utils.Error("invalid parameter count")

	// field error types
	ErrTypeReadOnly     = 1
	ErrTypeMissingValue = 2
	ErrTypeInvalidValue = 3
	ErrTypeNotSupported = 4
)

// field validation errors
type FieldError struct {
	FieldName  string
	ErrorType  int
	FieldError error
}

func ErrReadOnly(fieldName string) FieldError {
	return FieldError{
		FieldName:  fieldName,
		ErrorType:  ErrTypeReadOnly,
		FieldError: nil,
	}
}

func ErrMissingValue(fieldName string) FieldError {
	return FieldError{
		FieldName:  fieldName,
		ErrorType:  ErrTypeMissingValue,
		FieldError: nil,
	}
}

func ErrInvalidValue(fieldName string, fieldError error) FieldError {
	return FieldError{
		FieldName:  fieldName,
		ErrorType:  ErrTypeInvalidValue,
		FieldError: fieldError,
	}
}

func ErrNotSupported(fieldName string) FieldError {
	return FieldError{
		FieldName:  fieldName,
		ErrorType:  ErrTypeNotSupported,
		FieldError: nil,
	}
}

func (e FieldError) Error() string {
	switch e.ErrorType {
	case ErrTypeReadOnly:
		return fmt.Sprintf("field %s is not settable", e.FieldName)
	case ErrTypeMissingValue:
		return fmt.Sprintf("value for arg '%s' is missing", e.FieldName)
	case ErrTypeNotSupported:
		return fmt.Sprintf("non-supported type on arg %s", e.FieldName)
	default:
		return fmt.Sprintf("error parsing arg %s: %s", e.FieldName, e.FieldError.Error())
	}
}
