package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemService interface {
	GetAllItemLookup(filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse)
	GetAllItemInventory(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemById(Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	GetItemCode(string) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse)
	SaveItem(masteritempayloads.ItemRequest) (masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse)
	ChangeStatusItem(Id int) (bool, *exceptions.BaseErrorResponse)
	GetAllItemDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemDetailById(int, int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse)
	AddItemDetail(id int, req masteritempayloads.ItemDetailRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse)
	DeleteItemDetails(id int, itemDetailIDs []int) (masteritempayloads.DeleteItemResponse, *exceptions.BaseErrorResponse)
	GetUomTypeDropDown() ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse)
	GetUomDropDown(uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse)
	UpdateItem(int, masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateItemDetail(id int, itemDetail int, req masteritempayloads.ItemDetailUpdateRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse)
	GetPrincipalBrandParent(id int) ([]masteritempayloads.PrincipalBrandDropdownDescription, *exceptions.BaseErrorResponse)
	GetPrincipalBrandDropdown() ([]masteritempayloads.PrincipalBrandDropdownResponse, *exceptions.BaseErrorResponse)
	AddItemDetailByBrand(id string, itemId int) ([]masteritempayloads.ItemDetailResponse, *exceptions.BaseErrorResponse)
	GetAllItemSearch(filterCondition []utils.FilterCondition, itemIDs []string, supplierIDs []string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CheckItemCodeExist(itemCode string, itemGroupId int, commonPriceList bool, brandId int) (bool, int, int, *exceptions.BaseErrorResponse)
	GetPrincipalCatalog() ([]masteritempayloads.GetPrincipalCatalog, *exceptions.BaseErrorResponse)
}
