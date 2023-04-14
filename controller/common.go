package controller

import (
	"github.com/gofiber/fiber/v2"
)

const (
	BadRequest   = fiber.StatusBadRequest
	Conflict     = fiber.StatusConflict
	Created      = fiber.StatusCreated
	OK           = fiber.StatusOK
	NotFound     = fiber.StatusNotFound
	Unauthorized = fiber.StatusUnauthorized
)

// Map is alias for fiber.Map.
type Map = fiber.Map

// Message contains message returned to user.
type Message struct {
	Message string
}

// NewMessage returns  new Message object.
func NewMessage(message string) Message {
	return Message{Message: message}
}

// Contains ErrMessages  returned to user.
type ErrMessage struct {
	Error string
}

// NewErrMessage returns new ErrMessage object.
func NewErrMessage(errMessage string) ErrMessage {
	return ErrMessage{Error: errMessage}
}
