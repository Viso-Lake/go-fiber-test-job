package errorHelpers

import "github.com/gofiber/fiber/v2"

type ResponseBadRequestErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Bad request error"`
}

func NewResponseBadRequestErrorHTTP(message string) *ResponseBadRequestErrorHTTP {
	return &ResponseBadRequestErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondBadRequestError(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}
