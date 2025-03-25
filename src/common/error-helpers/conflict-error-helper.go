package errorHelpers

import (
	"github.com/gofiber/fiber/v2"
)

type ResponseConflictErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Conflict error"`
}

func NewResponseConflictErrorHTTP(message string) *ResponseConflictErrorHTTP {
	return &ResponseConflictErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondConflictError(message string) error {
	return fiber.NewError(fiber.StatusConflict, message)
}
