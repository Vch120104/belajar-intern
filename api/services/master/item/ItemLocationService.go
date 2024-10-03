package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type ItemLocationService interface {
	DeleteItemLocation(id int) *exceptions.BaseErrorResponse
	GetItemLocationById(id int) (masteritempayloads.ItemLocationRequest, *exceptions.BaseErrorResponse)
	SaveItemLocation(masteritempayloads.ItemLocationRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)

	GetAllItemLocationDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	PopupItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	AddItemLocation(int, masteritempayloads.ItemLocationDetailRequest) *exceptions.BaseErrorResponse

	GetAllItemLoc(filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetByIdItemLoc(id int) (masteritempayloads.ItemLocationGetByIdResponse, *exceptions.BaseErrorResponse)
	SaveItemLoc(req masteritempayloads.SaveItemlocation) (masteritementities.ItemLocation, *exceptions.BaseErrorResponse)
	DeleteItemLoc(ids []int) (bool, *exceptions.BaseErrorResponse)
	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
	UploadPreviewFile(rows [][]string) ([]masteritempayloads.UploadItemLocationResponse, *exceptions.BaseErrorResponse)
	UploadProcessFile(uploadPreview []masteritempayloads.UploadItemLocationResponse) ([]masteritementities.ItemLocation, *exceptions.BaseErrorResponse)
}
