package validations

import (
	"go-fiber-test-job/src/database/entities"
	addressValidationUtil "go-fiber-test-job/src/utils/address-validation"
	searchValidationUtil "go-fiber-test-job/src/utils/search-validation"
	"strings"

	"github.com/go-playground/validator/v10"
)

func AccountStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch entities.AccountStatus(status) {
	case entities.AccountStatusOn, entities.AccountStatusOff:
		return true
	}
	return false
}

func AccountAddressValidation(fl validator.FieldLevel) bool {
	address := fl.Field().String()
	if addressValidationUtil.IsValidAddress(address) {
		return true
	}
	return false
}

func NotEmpty(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	return strings.TrimSpace(str) != ""
}

func AccountSearchValidation(fl validator.FieldLevel) bool {
	search := fl.Field().String()
	if searchValidationUtil.IsValidSearch(search) {
		return true
	}
	return false
}
