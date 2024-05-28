package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemImportService interface {
	GetAllItemImport(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	GetItemImportbyId(Id int) (any, *exceptions.BaseErrorResponse)
	SaveItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
	UpdateItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse)
}
