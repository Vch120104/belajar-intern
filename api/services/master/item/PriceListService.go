package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type PriceListService interface {
	GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse)
	GetPriceListById(Id int) (masteritempayloads.PriceListGetbyId, *exceptions.BaseErrorResponse)
	SavePriceList(request masteritempayloads.SavePriceListMultiple) (int, *exceptions.BaseErrorResponse)
	ChangeStatusPriceList(Id int) (bool, *exceptions.BaseErrorResponse)
	GetAllPriceListNew(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	DeactivatePriceList(id string) (bool, *exceptions.BaseErrorResponse)
	ActivatePriceList(id string) (bool, *exceptions.BaseErrorResponse)
	DeletePriceList(id string) (bool, *exceptions.BaseErrorResponse)
	GenerateDownloadTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
	UploadFile(rows [][]string, uploadRequest masteritempayloads.PriceListUploadDataRequest) ([]string, *exceptions.BaseErrorResponse)
	ProcessUploadFile(upload []masteritempayloads.PriceListProcessdDataRequest) (bool, *exceptions.BaseErrorResponse)
	CheckPriceListItem(itemGroupId int, brandId int, currencyId int, date string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	Download(uploadRequest masteritempayloads.PriceListUploadDataRequest) (*excelize.File, *exceptions.BaseErrorResponse)
	Duplicate(itemGroupId int, brandId int, currencyId int, date string) ([]masteritempayloads.PriceListItemResponses, *exceptions.BaseErrorResponse)
}
