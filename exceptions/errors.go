package exceptions

import "fmt"

type ValidationError interface {
	Description() interface{}
}

type validationError struct {
	Field string `json:"Field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (v validationError) Description() interface{} {

	return v
}

func NewValidationError(field string, tag string, value string) ValidationError {
	return &validationError{Field: field, Tag: tag, Value: value}
}

type DuplicateKeyError struct {
	obj string
}

func NewDuplicateKeyError(objectName string) *DuplicateKeyError {
	return &DuplicateKeyError{obj: objectName}
}

func (d DuplicateKeyError) Error() string {
	return fmt.Sprintf("%s already exists", d.obj)
}
