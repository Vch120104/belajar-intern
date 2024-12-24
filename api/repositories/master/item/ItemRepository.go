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
	GetAllItemLookup(tx *gorm.DB, filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse)
	GetAllItemInventory(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
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
	GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdateItem(tx *gorm.DB, ItemId int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateItemDetail(tx *gorm.DB, Id int, itemDetail int, req masteritempayloads.ItemDetailUpdateRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse)
	GetPrincipalBrandParent(tx *gorm.DB, id int) ([]masteritempayloads.PrincipalBrandDropdownDescription, *exceptions.BaseErrorResponse)
	GetPrincipalBrandDropdown(tx *gorm.DB) ([]masteritempayloads.PrincipalBrandDropdownResponse, *exceptions.BaseErrorResponse)
	AddItemDetailByBrand(tx *gorm.DB, id string, itemId int) ([]masteritempayloads.ItemDetailResponse, *exceptions.BaseErrorResponse)
	GetAllItemSearch(tx *gorm.DB, filterCondition []utils.FilterCondition, itemIDs []string, supplierIDs []string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CheckItemCodeExist(tx *gorm.DB, itemCode string, itemGroupId int, commonPriceList bool, brandId int) (bool, int, int, *exceptions.BaseErrorResponse)
	GetPrincipalCatalog(tx *gorm.DB) ([]masteritempayloads.GetPrincipalCatalog, *exceptions.BaseErrorResponse)
	GetItemLatestId(tx *gorm.DB) (masteritempayloads.LatestItemAndLineTypeResponse, *exceptions.BaseErrorResponse)
	SaveItemToMappingItemOperation(tx *gorm.DB, response masteritempayloads.LatestItemAndLineTypeResponse) (bool, *exceptions.BaseErrorResponse)
}
