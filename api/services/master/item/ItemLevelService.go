package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemLevelService interface {
	Save(masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptions.BaseErrorResponse)
	GetById(itemLevel int, itemLevelid int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptions.BaseErrorResponse)
	GetItemLevelDropDown(itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptions.BaseErrorResponse)
	GetItemLevelLookUp(filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemLevelLookUpbyId(itemLevelId int) (masteritemlevelpayloads.GetItemLevelLookUp, *exceptions.BaseErrorResponse)
	GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatus(itemLevel int, itemLevelId int) (bool, *exceptions.BaseErrorResponse)
}
