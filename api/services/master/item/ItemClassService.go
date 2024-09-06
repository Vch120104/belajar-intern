package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemClassService interface {
	GetAllItemClass(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemClassDropDown() ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse)
	GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse)
	GetItemClassByCode(itemClassCode string) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse)
	SaveItemClass(req masteritempayloads.ItemClassResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemClass(Id int) (bool, *exceptions.BaseErrorResponse)
	GetItemClassDropDownbyGroupId(groupId int) ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse)
}
