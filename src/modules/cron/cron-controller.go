package cronModule

import (
	"go-fiber-test-job/src/common/dto"

	"github.com/gofiber/fiber/v2"
)

// UpdateAccountsBalances Update accounts balances
// @Summary Update accounts balances
// @Description Update accounts balances
// @Tags Cron
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Cron api key"
// @Success 201 {object} dto.SuccessDto
// @Failure 401 {object} errorHelpers.ResponseUnauthorizedErrorHTTP{}
// @Router /cron/account-balance [post]
func UpdateAccountsBalances(c *fiber.Ctx) error {
	updateAccountsBalances()
	return c.Status(fiber.StatusOK).JSON(dto.CreateSuccessDto())
}
