package errorHelpers

import "github.com/gofiber/fiber/v2"

type ResponseInternalErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Internal error"`
}

func NewResponseInternalErrorHTTP(message string) *ResponseInternalErrorHTTP {
	return &ResponseInternalErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondInternalError(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}
