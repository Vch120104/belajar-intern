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
	GetAllPurchasePrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetPurchasePriceById(tx *gorm.DB, Id int, pagination pagination.Pagination) (masteritempayloads.PurchasePriceResponse, *exceptions.BaseErrorResponse)
	SavePurchasePrice(tx *gorm.DB, req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
	UpdatePurchasePrice(tx *gorm.DB, Id int, req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
	ChangeStatusPurchasePrice(tx *gorm.DB, Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)

	GetAllPurchasePriceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetPurchasePriceDetailById(tx *gorm.DB, Id int) (masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse)
	AddPurchasePrice(tx *gorm.DB, req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse)
	UpdatePurchasePriceDetail(tx *gorm.DB, Id int, req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse)
	DeletePurchasePrice(tx *gorm.DB, Id int, iddet []int) (bool, *exceptions.BaseErrorResponse)
	ActivatePurchasePriceDetail(tx *gorm.DB, Id int, iddet []int) (bool, *exceptions.BaseErrorResponse)
	DeactivatePurchasePriceDetail(tx *gorm.DB, Id int, iddet []int) (bool, *exceptions.BaseErrorResponse)
}
