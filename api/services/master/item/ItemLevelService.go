package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemLevelService interface {
	Save(masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptionsss_test.BaseErrorResponse)
	GetItemLevelDropDown(itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemLevelLookUp(filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(int) (bool, *exceptionsss_test.BaseErrorResponse)
}
