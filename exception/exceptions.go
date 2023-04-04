package exception

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicateKey       = errors.New("object with the same name already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("password is required and must be at leat 4 characters long")
	ErrRecordNotFound     = errors.New("not found")
	ErrPasswordSame       = errors.New("password isn't new")
	ErrMalFormedJWT       = errors.New("missed or malformed token")
)

type ErrValidation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

func (v ErrValidation) Error() string {
	return fmt.Sprintf("ValidationErr: {field:%s, description:%s}", v.Field, v.Description)
}

func NewValidationError(field string, description string) ErrValidation {
	return ErrValidation{Field: field, Description: description}
}
