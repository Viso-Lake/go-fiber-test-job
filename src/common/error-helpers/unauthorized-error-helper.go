package errorHelpers

import "github.com/gofiber/fiber/v2"

type ResponseUnauthorizedErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Unauthorized error"`
}

func NewResponseUnauthorizedErrorHTTP(message string) *ResponseUnauthorizedErrorHTTP {
	return &ResponseUnauthorizedErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondUnauthorizedError() error {
	return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
}
