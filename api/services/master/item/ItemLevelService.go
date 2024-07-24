package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemLevelService interface {
	Save(masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptions.BaseErrorResponse)
	GetById(int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptions.BaseErrorResponse)
	GetItemLevelDropDown(itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptions.BaseErrorResponse)
	GetItemLevelLookUp(filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemLevelLookUpbyId(itemLevelId int) (masteritempayloads.GetItemLevelLookUp, *exceptions.BaseErrorResponse)
	GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatus(int) (bool, *exceptions.BaseErrorResponse)
}
