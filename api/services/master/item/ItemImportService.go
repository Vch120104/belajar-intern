package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type ItemImportService interface {
	GetAllItemImport(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	GetItemImportbyId(Id int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse)
	SaveItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
	UpdateItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
	GetItemImportbyItemIdandSupplierId(itemId int, supplierId int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse)
	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
	UploadPreviewFile(rows [][]string) ([]masteritempayloads.ItemImportUploadResponse, *exceptions.BaseErrorResponse)
	ProcessDataUpload(req masteritempayloads.ItemImportUploadRequest) (bool, *exceptions.BaseErrorResponse)
}
