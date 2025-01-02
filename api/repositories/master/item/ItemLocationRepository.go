package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLocationRepository interface {
	DeleteItemLocation(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	AddItemLocation(*gorm.DB, int, masteritempayloads.ItemLocationDetailRequest) (masteritementities.ItemLocationDetail, *exceptions.BaseErrorResponse)
	GetAllItemLocationDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	PopupItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)

	GetAllItemLoc(tx *gorm.DB, filtercondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemLoc(tx *gorm.DB, id int) (masteritempayloads.ItemLocationGetByIdResponse, *exceptions.BaseErrorResponse)
	SaveItemLoc(tx *gorm.DB, req masteritempayloads.SaveItemlocation) (masteritementities.ItemLocation, *exceptions.BaseErrorResponse)
	DeleteItemLoc(tx *gorm.DB, ids []int) (bool, *exceptions.BaseErrorResponse)
	IsDuplicateItemLoc(tx *gorm.DB, warehouseId int, warehouseLocationId int, itemId int) (bool, error)
}
