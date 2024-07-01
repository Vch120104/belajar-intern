package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemImportService interface {
	GetAllItemImport(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	GetItemImportbyId(Id int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse)
	SaveItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
	UpdateItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
	GetItemImportbyItemIdandSupplierId(itemId int, supplierId int) (masteritempayloads.ItemImportByIdResponse, *exceptions.BaseErrorResponse)
}
