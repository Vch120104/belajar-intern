package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemService interface {
	GetAllItem(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllItemLookup(filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse)
	GetItemById(Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemCode(string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveItem(masteritempayloads.ItemRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItem(Id int) (bool, *exceptions.BaseErrorResponse)
	GetAllItemDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetItemDetailById(int, int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse)
	AddItemDetail(int, masteritempayloads.ItemDetailRequest) *exceptions.BaseErrorResponse
	DeleteItemDetail(int, int) *exceptions.BaseErrorResponse
	GetUomTypeDropDown() ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse)
	GetUomDropDown(uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse)
	UpdateItem(int, masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateItemDetail(int, masteritempayloads.ItemDetailUpdateRequest) (bool, *exceptions.BaseErrorResponse)
	GetPrincipleBrandParent(code string) ([]masteritempayloads.PrincipleBrandDropdownDescription, *exceptions.BaseErrorResponse)
	GetPrincipleBrandDropdown() ([]masteritempayloads.PrincipleBrandDropdownResponse, *exceptions.BaseErrorResponse)
}
