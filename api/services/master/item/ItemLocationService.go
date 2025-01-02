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
	GetAllItemLocationDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	PopupItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	AddItemLocation(int, masteritempayloads.ItemLocationDetailRequest) (masteritementities.ItemLocationDetail, *exceptions.BaseErrorResponse)
	DeleteItemLocation(id int) *exceptions.BaseErrorResponse

	GetAllItemLoc(filtercondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemLoc(id int) (masteritempayloads.ItemLocationGetByIdResponse, *exceptions.BaseErrorResponse)
	SaveItemLoc(req masteritempayloads.SaveItemlocation) (masteritementities.ItemLocation, *exceptions.BaseErrorResponse)
	DeleteItemLoc(ids []int) (bool, *exceptions.BaseErrorResponse)

	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
	UploadPreviewFile(rows [][]string) ([]masteritempayloads.UploadItemLocationResponse, *exceptions.BaseErrorResponse)
	UploadProcessFile(uploadPreview []masteritempayloads.UploadItemLocationResponse) ([]masteritementities.ItemLocation, *exceptions.BaseErrorResponse)
}
