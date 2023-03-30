package exception

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicateKey       = errors.New("object with the same name already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type ErrValidation interface {
	Display() interface{}
	Error() string
}

type errValidation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

func (v errValidation) Display() interface{} {

	return v
}

func (v errValidation) Error() string {
	return fmt.Sprintf("ValidationErr: {field:%s, description:%s}", v.Field, v.Description)
}

func NewValidationError(field string, description string) ErrValidation {
	return &errValidation{Field: field, Description: description}
}
