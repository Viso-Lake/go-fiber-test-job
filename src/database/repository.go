package database

import (
	"fmt"
	"go-fiber-test-job/src/database/entities"
	"strings"

	"gorm.io/gorm"
)

func accountTableName() string {
	return entities.Account{}.TableName()
}

func getDb(tx *gorm.DB) *gorm.DB {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = DbConn
	}
	return db
}

///// Account queries

func GetAccountsAndTotal(search string, status entities.AccountStatus, orderParams map[string]string, offset int, count int) ([]*entities.Account, int64) {
	var total int64
	var accounts []*entities.Account
	query := getBaseAccountsQuery(status)
	totalQuery := getBaseAccountsQuery(status)

	search = strings.TrimSpace(search)
	if search != "" {
		searchPattern := "%" + strings.ReplaceAll(search, "%", "\\%") + "%"
		queryRaw := "account.address LIKE ? OR account.name LIKE ? OR account.memo LIKE ?"
		query = query.Where(queryRaw, searchPattern, searchPattern, searchPattern)
		totalQuery = totalQuery.Where(queryRaw, searchPattern, searchPattern, searchPattern)
	}

	allowedSortFields := map[string]bool{
		"id":         true,
		"updated_at": true,
		"address":    true,
		"name":       true,
		"rank":       true,
	}

	for key, value := range orderParams {
		if allowedSortFields[key] && (value == "ASC" || value == "DESC") {
			query = query.Order(fmt.Sprintf("account.%s %s", key, value))
		}
	}

	query.
		Limit(count).
		Offset(offset).
		Find(&accounts)
	totalQuery.Count(&total)
	return accounts, total
}

func getBaseAccountsQuery(status entities.AccountStatus) *gorm.DB {
	query := DbConn.Table(accountTableName() + " account")
	if status != "" {
		query = query.Where("account.status = ?", status)
	}
	return query
}

func IsAddressExists(tx *gorm.DB, address string) bool {
	db := getDb(tx)
	var account *entities.Account
	db.Table(accountTableName()+" account").
		Where("account.address = ?", address).
		First(&account)
	if account.Id != 0 {
		return true
	}
	return false
}

func GetAccountByAddress(address string) *entities.Account {
	var account *entities.Account
	DbConn.Table(accountTableName()+" account").
		Where("account.address = ?", address).
		First(&account)
	if account.Id == 0 {
		return nil
	}

	return account
}

func CreateAccount(tx *gorm.DB, newAccount *entities.Account) (*entities.Account, error) {
	err := tx.Create(newAccount).Error
	if err != nil {
		return nil, err
	}
	return newAccount, nil
}

func GetAccountsBatch(limit int) []*entities.Account {
	var accounts []*entities.Account
	DbConn.Table(accountTableName()+" account").
		Where("account.status = ?", entities.AccountStatusOn).
		Order("account.updated_at ASC").
		Limit(limit).
		Find(&accounts)
	return accounts
}

func GetAccountsByIds(accountIds []int64) []*entities.Account {
	var accounts []*entities.Account
	DbConn.Table(accountTableName()+" account").
		Where("account.id IN(?)", accountIds).
		Find(&accounts)
	return accounts
}

func UpdateAccount(tx *gorm.DB, account *entities.Account, updateData map[string]interface{}) error {
	db := getDb(tx)
	return db.Model(entities.Account{}).Where("id = ?", account.Id).Updates(updateData).Error
}
