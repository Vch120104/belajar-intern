package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllItemLookup(tx *gorm.DB, filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse)
	GetItemById(tx *gorm.DB, Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemCode(*gorm.DB, string) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	SaveItem(*gorm.DB, masteritempayloads.ItemRequest) (masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse)
	ChangeStatusItem(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	GetUomTypeDropDown(tx *gorm.DB) ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse)
	GetUomDropDown(tx *gorm.DB, uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse)
	AddItemDetail(tx *gorm.DB, ItemId int, req masteritempayloads.ItemDetailRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse)
	DeleteItemDetails(tx *gorm.DB, ItemId int, itemDetailIDs []int) (masteritempayloads.DeleteItemResponse, *exceptions.BaseErrorResponse)
	GetItemDetailById(*gorm.DB, int, int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse)
	GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateItem(tx *gorm.DB, ItemId int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateItemDetail(tx *gorm.DB, Id int, itemDetail int, req masteritempayloads.ItemDetailUpdateRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse)
	GetPrincipleBrandParent(tx *gorm.DB, id int) ([]masteritempayloads.PrincipleBrandDropdownDescription, *exceptions.BaseErrorResponse)
	GetPrincipleBrandDropdown(tx *gorm.DB) ([]masteritempayloads.PrincipleBrandDropdownResponse, *exceptions.BaseErrorResponse)
	AddItemDetailByBrand(tx *gorm.DB, id string, itemId int) ([]masteritempayloads.ItemDetailResponse, *exceptions.BaseErrorResponse)
	GetAllItemSearch(tx *gorm.DB, filterCondition []utils.FilterCondition, itemIDs []string, supplierIDs []string, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CheckItemCodeExist(tx *gorm.DB, itemCode string, itemGroupId int, commonPriceList bool, brandId int) (bool, int, int, *exceptions.BaseErrorResponse)
	GetCatalogCode(tx *gorm.DB) ([]masteritempayloads.GetCatalogCode, *exceptions.BaseErrorResponse)
	GetAllItemListTransLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
