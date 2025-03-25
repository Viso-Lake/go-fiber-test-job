package validations

import (
	"github.com/go-playground/validator/v10"
	"go-fiber-test-job/src/database/entities"
	addressValidationUtil "go-fiber-test-job/src/utils/address-validation"
	"strings"
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
