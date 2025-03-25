package accountTests

import (
	"bytes"
	"encoding/json"
	"fmt"
	errorHelpers "go-fiber-test-job/src/common/error-helpers"
	"go-fiber-test-job/src/config"
	"go-fiber-test-job/src/database"
	"go-fiber-test-job/src/database/entities"
	accountModuleDto "go-fiber-test-job/src/modules/account/dto"
	arrayUtil "go-fiber-test-job/src/utils/array"
	numberUtil "go-fiber-test-job/src/utils/number"
	orderUtil "go-fiber-test-job/src/utils/order"
	timeUtil "go-fiber-test-job/src/utils/time"
	"go-fiber-test-job/test"
	"go-fiber-test-job/test/seeds"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountRoute(t *testing.T) {
	// GetAccounts
	validationGetAccountsTests(t)
	t.Run("TestGetAccountsRoute_SuccessNoParams", TestGetAccountsRoute_SuccessNoParams)
	t.Run("TestGetAccountsRoute_SuccessParamsOffsetAndCount", TestGetAccountsRoute_SuccessParamsOffsetAndCount)
	t.Run("TestGetAccountsRoute_SuccessParamsStatus", TestGetAccountsRoute_SuccessParamsStatus)
	t.Run("TestGetAccountsRoute_SuccessParamsOrderBy", TestGetAccountsRoute_SuccessParamsOrderBy)
	t.Run("TestGetAccountsRoute_SuccessParamsStatusAndOrderBy", TestGetAccountsRoute_SuccessParamsStatusAndOrderBy)
	t.Run("TestGetAccountsRoute_SuccessParamsOffsetAndCountAndStatusAndOrderBy", TestGetAccountsRoute_SuccessParamsOffsetAndCountAndStatusAndOrderBy)
	// CreateAccount
	validationCreateAccountTests(t)
	t.Run("TestCreateAccountRoute_FailAddressAlreadyExists", TestCreateAccountRoute_FailAddressAlreadyExists)
	t.Run("TestCreateAccountRoute_Success", TestCreateAccountRoute_Success)
}

func validationGetAccountsTests(t *testing.T) {
	validationTests := []struct {
		name         string
		params       accountModuleDto.GetAccountRequestDto
		expectedCode int
		expectedBody errorHelpers.ResponseBadRequestErrorHTTP
	}{
		{
			"FailInvalidOffsetMinValue",
			accountModuleDto.GetAccountRequestDto{Offset: -5},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Offset must be greater than or equal 0"},
		},
		{
			"FailInvalidCountMinValue",
			accountModuleDto.GetAccountRequestDto{Count: -1},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Count must be greater than or equal 1"},
		},
		{
			"FailInvalidCountMaxValue",
			accountModuleDto.GetAccountRequestDto{Count: 101},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Count must be less than or equal 100"},
		},
		{
			"FailInvalidStatus",
			accountModuleDto.GetAccountRequestDto{Status: "invalid status"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: fmt.Sprintf("%s must be one of the next values: %s", "Status", strings.Join(entities.AccountStatusList, ","))},
		},
		{
			"FailInvalidOrderBy",
			accountModuleDto.GetAccountRequestDto{OrderBy: "invalid order by"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "invalid order by parameter: invalid order by"},
		},
		{
			"FailInvalidOrderByMaxLength",
			accountModuleDto.GetAccountRequestDto{OrderBy: strings.Repeat("OrderBy", 255)},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "OrderBy must be shorter than or equal to 255 characters"},
		},
	}
	for _, validationTest := range validationTests {
		t.Run("TestGetAccountsRoute"+validationTest.name, func(t *testing.T) {
			type Params struct {
				Count   int                    `json:"count"`
				Offset  int                    `json:"offset"`
				Search  string                 `json:"search"`
				Status  entities.AccountStatus `json:"status"`
				OrderBy string                 `json:"orderBy"`
			}
			params := &Params{
				Count:   validationTest.params.Count,
				Offset:  validationTest.params.Offset,
				Search:  validationTest.params.Search,
				Status:  validationTest.params.Status,
				OrderBy: validationTest.params.OrderBy,
			}

			query := url.Values{}
			query.Add("Count", numberUtil.IntToString(params.Count))
			query.Add("Offset", numberUtil.IntToString(params.Offset))
			query.Add("Search", string(params.Search))
			query.Add("Status", string(params.Status))
			query.Add("OrderBy", params.OrderBy)

			u := &url.URL{
				Path:     fmt.Sprintf("/account"),
				RawQuery: query.Encode(),
			}

			request := httptest.NewRequest("GET", u.String(), nil)
			request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
			response, _ := test.TestApp.Test(request)
			assert.Equal(t, validationTest.expectedCode, response.StatusCode)

			// Read the response body and parse JSON
			var responseDto errorHelpers.ResponseBadRequestErrorHTTP
			err := json.NewDecoder(response.Body).Decode(&responseDto)
			assert.Nil(t, err)

			assert.NotNil(t, responseDto.Success, "Success parameter should exist")
			assert.NotNil(t, responseDto.Message, "Message parameter should exist")

			assert.Equal(t, validationTest.expectedBody.Success, responseDto.Success)
			assert.Equal(t, validationTest.expectedBody.Message, responseDto.Message)
		})
	}
}

func TestGetAccountsRoute_SuccessNoParams(t *testing.T) {
	accounts, total := database.GetAccountsAndTotal("", "", make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	responseDto := sendGetAccountRequest(t, url.Values{})

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, 0, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}
}

func TestGetAccountsRoute_SuccessParamsOffsetAndCount(t *testing.T) {
	type Params struct {
		Count  int `json:"count"`
		Offset int `json:"offset"`
	}
	params := &Params{
		Count:  2,
		Offset: 1,
	}

	query := url.Values{}
	query.Add("Count", numberUtil.IntToString(params.Count))
	query.Add("Offset", numberUtil.IntToString(params.Offset))
	responseDto := sendGetAccountRequest(t, query)

	accounts, total := database.GetAccountsAndTotal("", "", make(map[string]string), params.Offset, params.Count)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, params.Offset, responseDto.Offset)
	assert.Equal(t, params.Count, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}
}

func TestGetAccountsRoute_SuccessParamsStatus(t *testing.T) {
	type Params struct {
		Status entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Status: entities.AccountStatusOn,
	}

	query := url.Values{}
	query.Add("Status", string(params.Status))
	responseDto := sendGetAccountRequest(t, query)

	accounts, total := database.GetAccountsAndTotal("", params.Status, make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}
}

func TestGetAccountsRoute_SuccessParamsOrderBy(t *testing.T) {
	type Params struct {
		OrderBy string `json:"orderBy"`
	}
	params := &Params{
		OrderBy: "id DESC",
	}

	query := url.Values{}
	query.Add("OrderBy", params.OrderBy)
	responseDto := sendGetAccountRequest(t, query)

	orderParams, err := orderUtil.GetOrderByParamsSecure(params.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	assert.Nil(t, err, "GetOrderByParamsSecure вернул ошибку")

	accounts, total := database.GetAccountsAndTotal("", "", orderParams, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}

	assert.Equal(t, true, test.TestListSort(responseDto.List, params.OrderBy), "List is not sorted")
}

func TestGetAccountsRoute_SuccessParamsStatusAndOrderBy(t *testing.T) {
	type Params struct {
		Status  entities.AccountStatus `json:"status"`
		OrderBy string                 `json:"orderBy"`
	}
	params := &Params{
		Status:  entities.AccountStatusOff,
		OrderBy: "updated_at DESC",
	}

	query := url.Values{}
	query.Add("Status", string(params.Status))
	query.Add("OrderBy", params.OrderBy)

	orderParams, err := orderUtil.GetOrderByParamsSecure(params.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	assert.Nil(t, err, "GetOrderByParamsSecure вернул ошибку")

	accounts, total := database.GetAccountsAndTotal("", params.Status, orderParams, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	responseDto := sendGetAccountRequest(t, query)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}

	assert.Equal(t, true, test.TestListSort(responseDto.List, params.OrderBy), "List is not sorted")
}

func TestGetAccountsRoute_SuccessParamsOffsetAndCountAndStatusAndOrderBy(t *testing.T) {
	type Params struct {
		Count   int                    `json:"count"`
		Offset  int                    `json:"offset"`
		Search  string                 `json:"search"`
		Status  entities.AccountStatus `json:"status"`
		OrderBy string                 `json:"orderBy"`
	}
	params := &Params{
		Count:   2,
		Offset:  0,
		Search:  "sato",
		Status:  entities.AccountStatusOn,
		OrderBy: "updated_at ASC",
	}

	query := url.Values{}
	query.Add("Count", numberUtil.IntToString(params.Count))
	query.Add("Offset", numberUtil.IntToString(params.Offset))
	query.Add("Search", string(params.Search))
	query.Add("Status", string(params.Status))
	query.Add("OrderBy", params.OrderBy)

	orderParams, err := orderUtil.GetOrderByParamsSecure(params.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	assert.Nil(t, err, "GetOrderByParamsSecure вернул ошибку")

	accounts, total := database.GetAccountsAndTotal(params.Search, params.Status, orderParams, params.Offset, params.Count)

	responseDto := sendGetAccountRequest(t, query)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, params.Offset, responseDto.Offset)
	assert.Equal(t, params.Count, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}

	assert.Equal(t, true, test.TestListSort(responseDto.List, params.OrderBy), "List is not sorted")
}

func validationCreateAccountTests(t *testing.T) {
	validationTests := []struct {
		name         string
		params       accountModuleDto.PostCreateAccountRequestDto
		expectedCode int
		expectedBody errorHelpers.ResponseBadRequestErrorHTTP
	}{
		{
			"FailNoBody",
			accountModuleDto.PostCreateAccountRequestDto{},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Address format is wrong"},
		},
		{
			"FailInvalidAddress",
			accountModuleDto.PostCreateAccountRequestDto{Address: "invalid address"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Address format is wrong"},
		},
		{
			"FailInvalidStatus",
			accountModuleDto.PostCreateAccountRequestDto{Address: "14yqg2y3a6HMgW9MiF5tVPAH4Dr1uxGKFJ", Status: "invalid status"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: fmt.Sprintf("%s must be one of the next values: %s", "Status", strings.Join(entities.AccountStatusList, ","))},
		},
	}
	for _, validationTest := range validationTests {
		t.Run("TestCreateAccountRoute"+validationTest.name, func(t *testing.T) {
			type Params struct {
				Address string                 `json:"address"`
				Name    string                 `json:"name"`
				Rank    int                    `json:"rank"`
				Memo    string                 `json:"memo"`
				Status  entities.AccountStatus `json:"status"`
			}
			params := &Params{
				Name:    validationTest.params.Name,
				Rank:    validationTest.params.Rank,
				Memo:    validationTest.params.Memo,
				Address: validationTest.params.Address,
				Status:  validationTest.params.Status,
			}
			body, _ := json.Marshal(params)

			u := &url.URL{
				Path: fmt.Sprintf("/account"),
			}

			request := httptest.NewRequest("POST", u.String(), bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
			response, _ := test.TestApp.Test(request)
			assert.Equal(t, validationTest.expectedCode, response.StatusCode)

			// Read the response body and parse JSON
			var responseDto errorHelpers.ResponseBadRequestErrorHTTP
			err := json.NewDecoder(response.Body).Decode(&responseDto)
			assert.Nil(t, err)

			assert.NotNil(t, responseDto.Success, "Success parameter should exist")
			assert.NotNil(t, responseDto.Message, "Message parameter should exist")

			assert.Equal(t, validationTest.expectedBody.Success, responseDto.Success)
			assert.Equal(t, validationTest.expectedBody.Message, responseDto.Message)
		})
	}
}

func TestCreateAccountRoute_FailAddressAlreadyExists(t *testing.T) {
	accountInfo := seeds.ACCOUNTS.ACCOUNT_1
	type Params struct {
		Address string                 `json:"address"`
		Name    string                 `json:"name"`
		Rank    int                    `json:"rank"`
		Memo    string                 `json:"memo"`
		Status  entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Name:    accountInfo.Name,
		Rank:    accountInfo.Rank,
		Memo:    accountInfo.Memo,
		Address: accountInfo.Address,
		Status:  entities.AccountStatusOn,
	}
	body, _ := json.Marshal(params)

	u := &url.URL{
		Path: fmt.Sprintf("/account"),
	}

	assert.Equal(t, true, database.IsAddressExists(nil, params.Address), "Address must exists")

	request := httptest.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	response, _ := test.TestApp.Test(request)
	assert.Equal(t, http.StatusConflict, response.StatusCode)

	// Read the response body and parse JSON
	var responseDto errorHelpers.ResponseBadRequestErrorHTTP
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Success, "Success parameter should exist")
	assert.NotNil(t, responseDto.Message, "Message parameter should exist")

	assert.Equal(t, false, responseDto.Success)
	assert.Equal(t, "Address already exists", responseDto.Message)
}

func TestCreateAccountRoute_Success(t *testing.T) {
	start := timeUtil.GetUnixTime()
	type Params struct {
		Address string                 `json:"address"`
		Name    string                 `json:"name"`
		Rank    int                    `json:"rank"`
		Memo    string                 `json:"memo"`
		Status  entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Address: "32AaKxGbdhGMSGutcZjspFq9U89jJHW1um",
		Name:    "Satoshi 5",
		Rank:    89,
		Memo:    "memorandum text 5",
		Status:  entities.AccountStatusOn,
	}
	body, _ := json.Marshal(params)

	u := &url.URL{
		Path: fmt.Sprintf("/account"),
	}

	assert.Equal(t, false, database.IsAddressExists(nil, params.Address), "Address must not exists")

	request := httptest.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	response, _ := test.TestApp.Test(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.AccountDto
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Id, "Id parameter should exist")
	assert.NotNil(t, responseDto.Address, "Address parameter should exist")
	assert.NotNil(t, responseDto.Balance, "Balance parameter should exist")
	assert.NotNil(t, responseDto.Status, "Status parameter should exist")
	assert.NotNil(t, responseDto.CreatedAt, "CreatedAt parameter should exist")
	assert.NotNil(t, responseDto.UpdatedAt, "UpdatedAt parameter should exist")

	accountAfter := database.GetAccountByAddress(params.Address)
	assert.NotNil(t, accountAfter)

	assert.Equal(t, responseDto.Id, accountAfter.Id)
	assert.Equal(t, responseDto.Address, accountAfter.Address)
	assert.Equal(t, responseDto.Balance, accountAfter.Balance.String())
	assert.Equal(t, responseDto.Status, string(accountAfter.Status))
	assert.Equal(t, responseDto.CreatedAt, accountAfter.CreatedAt)
	assert.Equal(t, responseDto.UpdatedAt, accountAfter.UpdatedAt)
	assert.GreaterOrEqual(t, responseDto.CreatedAt, start)
	assert.GreaterOrEqual(t, responseDto.UpdatedAt, start)

	test.CompareAccount(t, accountAfter, responseDto)
}

func sendGetAccountRequest(t *testing.T, queryParams url.Values) accountModuleDto.GetAccountResponseDto {
	u := &url.URL{
		Path:     "/account",
		RawQuery: queryParams.Encode(),
	}

	request := httptest.NewRequest("GET", u.String(), nil)
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	response, _ := test.TestApp.Test(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseDto accountModuleDto.GetAccountResponseDto
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	return responseDto
}
