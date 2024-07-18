package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PurchasePriceService interface {
	DeletePurchasePrice(id int) *exceptions.BaseErrorResponse
	GetPurchasePriceById(id int) (masteritempayloads.PurchasePriceRequest, *exceptions.BaseErrorResponse)
	SavePurchasePrice(masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
	GetAllPurchasePrice(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllPurchasePriceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetPurchasePriceDetailById(id int, pages pagination.Pagination) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	AddPurchasePrice(masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse)
	ChangeStatusPurchasePrice(Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
}
