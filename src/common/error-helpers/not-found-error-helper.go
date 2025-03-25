package errorHelpers

import "github.com/gofiber/fiber/v2"

type ResponseNotFoundErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Not found error"`
}

func NewResponseNotFoundErrorHTTP(message string) *ResponseNotFoundErrorHTTP {
	return &ResponseNotFoundErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondNotFoundError(message string) error {
	return fiber.NewError(fiber.StatusNotFound, message)
}
