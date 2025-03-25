package accountModuleDto

import (
	"go-fiber-test-job/src/database/entities"
)

func CreatePostCreateAccountResponseDto(account *entities.Account) AccountDto {
	return CreateAccountDto(account)
}
