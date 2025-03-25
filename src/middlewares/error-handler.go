package middleware

import (
	"github.com/gofiber/fiber/v2"
	errorHelpers "go-fiber-test-job/src/common/error-helpers"
)

func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			// Check if the error is a Fiber error
			if fiberErr, ok := err.(*fiber.Error); ok {
				switch fiberErr.Code {
				case fiber.StatusBadRequest:
					return c.Status(fiber.StatusBadRequest).JSON(errorHelpers.NewResponseBadRequestErrorHTTP(fiberErr.Message))
				case fiber.StatusUnauthorized:
					return c.Status(fiber.StatusUnauthorized).JSON(errorHelpers.NewResponseUnauthorizedErrorHTTP(fiberErr.Message))
				case fiber.StatusNotFound:
					return c.Status(fiber.StatusNotFound).JSON(errorHelpers.NewResponseNotFoundErrorHTTP(fiberErr.Message))
				case fiber.StatusConflict:
					return c.Status(fiber.StatusConflict).JSON(errorHelpers.NewResponseConflictErrorHTTP(fiberErr.Message))
				case fiber.StatusInternalServerError:
					return c.Status(fiber.StatusInternalServerError).JSON(errorHelpers.NewResponseInternalErrorHTTP(fiberErr.Message))
				default:
					return c.Status(fiber.StatusInternalServerError).JSON(errorHelpers.NewResponseInternalErrorHTTP(fiberErr.Message))
				}
			}
			// Return a consistent error response
			return c.Status(fiber.StatusInternalServerError).JSON(errorHelpers.NewResponseInternalErrorHTTP(err.Error()))
		}
		return nil
	}
}
