package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLocationRepository interface {
	DeleteItemLocation(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	AddItemLocation(*gorm.DB, int, masteritempayloads.ItemLocationDetailRequest) *exceptions.BaseErrorResponse
	SaveItemLocation(*gorm.DB, masteritempayloads.ItemLocationRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetItemLocationById(tx *gorm.DB, Id int) (masteritempayloads.ItemLocationRequest, *exceptions.BaseErrorResponse)
	GetAllItemLocationDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	PopupItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
