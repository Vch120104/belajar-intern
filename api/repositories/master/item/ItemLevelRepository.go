package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLevelRepository interface {
	Save(*gorm.DB, masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(*gorm.DB, int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptionsss_test.BaseErrorResponse)
	GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetItemLevelDropDown(tx *gorm.DB, itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemLevelLookUp(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
