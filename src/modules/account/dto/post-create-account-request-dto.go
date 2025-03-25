package accountModuleDto

import (
	"fmt"
	errorHelpers "go-fiber-test-job/src/common/error-helpers"
	errorMessages "go-fiber-test-job/src/common/error-messages"
	"go-fiber-test-job/src/common/validations"
	"go-fiber-test-job/src/database/entities"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PostCreateAccountRequestDto struct {
	Name    string                 `json:"name" validate:"max=255" example:"Satoshi"`
	Rank    int                    `json:"rank" validate:"gte=0,lte=100" example:"50"`
	Memo    string                 `json:"memo" validate:"omitempty" example:"Memorandum text"`
	Address string                 `json:"address" validate:"AccountAddressValidation" example:"1JzfdUygUFk2M6KS3ngFMGRsy5vsH4N37a"`
	Status  entities.AccountStatus `json:"status" validate:"AccountStatusValidation" enums:"On,Off" example:"On"`
}

var postCreateAccountRequestDtoValidator *validator.Validate

func init() {
	postCreateAccountRequestDtoValidator = validator.New()
	_ = postCreateAccountRequestDtoValidator.RegisterValidation("AccountAddressValidation", validations.AccountAddressValidation)
	_ = postCreateAccountRequestDtoValidator.RegisterValidation("AccountStatusValidation", validations.AccountStatusValidation)
}

func validatePostCreateAccountRequestDto(dto *PostCreateAccountRequestDto) error {
	return postCreateAccountRequestDtoValidator.Struct(dto)
}

func CreatePostCreateAccountRequestDto(c *fiber.Ctx) (PostCreateAccountRequestDto, error) {
	var dto PostCreateAccountRequestDto
	// Parse body params into DTO
	if err := c.BodyParser(&dto); err != nil {
		errorMessage := PostCreateAccountRequestDtoQueryParseErrorMessage(err)
		return dto, errorHelpers.RespondBadRequestError(errorMessage)
	}
	// Validate the DTO
	if err := validatePostCreateAccountRequestDto(&dto); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := PostCreateAccountRequestDtoValidateErrorMessage(err)
			return dto, errorHelpers.RespondBadRequestError(errorMessage)
		}
	}
	return dto, nil
}

func PostCreateAccountRequestDtoQueryParseErrorMessage(err error) string {
	return errorMessages.DefaultQueryParseErrorMessage()
}

func PostCreateAccountRequestDtoValidateErrorMessage(err validator.FieldError) string {
	var errorMessage string
	if err.Field() == "Address" && err.Tag() == "AccountAddressValidation" {
		errorMessage = fmt.Sprintf("%s format is wrong", err.Field())
	} else if err.Field() == "Status" && err.Tag() == "AccountStatusValidation" {
		errorMessage = fmt.Sprintf("%s must be one of the next values: %s", err.Field(), strings.Join(entities.AccountStatusList, ","))
	} else if err.Field() == "Name" && err.Tag() == "max" {
		errorMessage = fmt.Sprintf("%s must be shorter than or equal to %s characters", err.Field(), err.Param())
	} else if err.Field() == "Rank" && err.Tag() == "gte" {
		errorMessage = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
	} else if err.Field() == "Rank" && err.Tag() == "lte" {
		errorMessage = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
	} else {
		errorMessage = errorMessages.DefaultFieldErrorMessage(err.Field())
	}
	return errorMessage
}
