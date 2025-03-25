package accountModule

import (
	"github.com/gofiber/fiber/v2"
	accountModuleDto "go-fiber-test-job/src/modules/account/dto"
	orderUtil "go-fiber-test-job/src/utils/order"
)

// GetAccounts Get list of accounts
// @Summary Get list of accounts
// @Description Get list of account
// @Tags Account
// @Accept json
// @Produce json
// @Param offset query int false "This is paging offset. 0 by default" minimum(0) default(0)
// @Param count query int false "Max item count in single response. 100 by default" minimum(1) maximum(100) default(100)
// @Param status query string false "Account statuses: On, Off" Enums(On,Off) default(On)
// @Param orderBy query string false "Comma-separated sort order options (sort fields: id, updated, sort order: ASC,DESC)" default(id ASC)
// @Param X-API-Key header string true "Admin api key"
// @Success 200 {object} accountModuleDto.GetAccountResponseDto
// @Failure 400 {object} errorHelpers.ResponseBadRequestErrorHTTP{}
// @Failure 401 {object} errorHelpers.ResponseUnauthorizedErrorHTTP{}
// @Router /account [get]
func GetAccounts(c *fiber.Ctx) error {
	dto, err := accountModuleDto.CreateGetAccountRequestDto(c)
	if err != nil {
		return err
	}
	orderParams, err := orderUtil.GetOrderByParamsSecure(dto.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	if err != nil {
		return err
	}
	accounts, total := getAccounts(dto.Status, orderParams, dto.Offset, dto.Count)
	return c.Status(fiber.StatusOK).JSON(accountModuleDto.CreateGetAccountResponseDto(dto.Offset, dto.Count, total, accounts))
}

// CreateAccount Create new account
// @Summary Create new account
// @Description Create new account
// @Tags Account
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin api key"
// @Param request body accountModuleDto.PostCreateAccountRequestDto true "Request body"
// @Success 200 {object} accountModuleDto.AccountDto
// @Failure 400 {object} errorHelpers.ResponseBadRequestErrorHTTP{}
// @Failure 401 {object} errorHelpers.ResponseUnauthorizedErrorHTTP{}
// @Failure 409 {object} errorHelpers.ResponseConflictErrorHTTP{}
// @Router /account [post]
func CreateAccount(c *fiber.Ctx) error {
	dto, err := accountModuleDto.CreatePostCreateAccountRequestDto(c)
	if err != nil {
		return err
	}
	account, err := createAccount(dto.Address, dto.Status)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(accountModuleDto.CreatePostCreateAccountResponseDto(account))
}
