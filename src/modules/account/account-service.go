package accountModule

import (
	errorHelpers "go-fiber-test-job/src/common/error-helpers"
	"go-fiber-test-job/src/database"
	"go-fiber-test-job/src/database/entities"
	"gorm.io/gorm"
)

func getAccounts(status entities.AccountStatus, orderParams map[string]string, offset int, count int) ([]*entities.Account, int64) {
	return database.GetAccountsAndTotal(status, orderParams, offset, count)
}

func createAccount(address string, status entities.AccountStatus) (*entities.Account, error) {
	var account *entities.Account
	transactionError := database.DbConn.Transaction(func(tx *gorm.DB) error {
		if database.IsAddressExists(tx, address) {
			return errorHelpers.RespondConflictError("Address already exists")
		}
		newAccount, err := database.CreateAccount(tx, entities.CreateAccount(address, status))
		if err != nil {
			return err
		}
		account = newAccount
		return nil
	}, database.DefaultTxOptions)
	if transactionError != nil {
		return nil, transactionError
	}
	return account, nil
}
