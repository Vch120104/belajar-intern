package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PurchasePriceRepository interface {
	DeletePurchasePrice(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse
	AddPurchasePrice(*gorm.DB, masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse)
	SavePurchasePrice(*gorm.DB, masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
	GetAllPurchasePrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetPurchasePriceById(tx *gorm.DB, Id int) (masteritempayloads.PurchasePriceRequest, *exceptions.BaseErrorResponse)
	GetAllPurchasePriceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetPurchasePriceDetailById(tx *gorm.DB, Id int, pages pagination.Pagination) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ChangeStatusPurchasePrice(tx *gorm.DB, Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
}
