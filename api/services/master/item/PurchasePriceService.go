package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type PurchasePriceService interface {
	GetAllPurchasePrice(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetPurchasePriceById(id int, pagination pagination.Pagination) (masteritempayloads.PurchasePriceResponse, *exceptions.BaseErrorResponse)
	UpdatePurchasePrice(id int, req masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
	SavePurchasePrice(masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)
	ChangeStatusPurchasePrice(Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse)

	GetAllPurchasePriceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetPurchasePriceDetailById(id int) (masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse)
	AddPurchasePrice(masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse)
	UpdatePurchasePriceDetail(Id int, req masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse)
	DeletePurchasePrice(id int, iddet []int) (bool, *exceptions.BaseErrorResponse)

	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
	//UploadPreviewFile(rows [][]string) ([]masteritempayloads.PurchasePriceDetailRequest, *exceptions.BaseErrorResponse)
	//ProcessDataUpload(req masteritempayloads.UploadRequest) (bool, *exceptions.BaseErrorResponse)
}
