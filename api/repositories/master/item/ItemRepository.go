package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllItemLookup(tx *gorm.DB, filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse)
	GetItemById(tx *gorm.DB, Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemCode(*gorm.DB, string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveItem(*gorm.DB, masteritempayloads.ItemRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItem(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	GetUomTypeDropDown(tx *gorm.DB) ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse)
	GetUomDropDown(tx *gorm.DB, uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse)
	AddItemDetail(*gorm.DB, int, masteritempayloads.ItemDetailRequest) *exceptions.BaseErrorResponse
	DeleteItemDetail(*gorm.DB, int, int) *exceptions.BaseErrorResponse
	GetItemDetailById(*gorm.DB, int, int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse)
	GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateItem(tx *gorm.DB, ItemId int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateItemDetail(tx *gorm.DB, ItemId int, req masteritempayloads.ItemDetailUpdateRequest) (bool, *exceptions.BaseErrorResponse)
	GetPrincipleBrandParent(tx *gorm.DB, code string) ([]masteritempayloads.PrincipleBrandDropdownDescription, *exceptions.BaseErrorResponse)
	GetPrincipleBrandDropdown(tx *gorm.DB) ([]masteritempayloads.PrincipleBrandDropdownResponse, *exceptions.BaseErrorResponse)
}
