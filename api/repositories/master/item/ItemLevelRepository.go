package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemLevelRepository interface {
	Save(*gorm.DB, masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(*gorm.DB, int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptionsss_test.BaseErrorResponse)
	GetAll(tx *gorm.DB, request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
