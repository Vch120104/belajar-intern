package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"github.com/xuri/excelize/v2"
)

type BomService interface {
	GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomMasterById(id int, pages pagination.Pagination) (masteritempayloads.BomMasterResponseDetail, *exceptions.BaseErrorResponse)
	SaveBomMaster(request masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	UpdateBomMaster(id int, request masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	ChangeStatusBomMaster(id int) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	SaveBomDetail(request masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	UpdateBomDetail(id int, request masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	DeleteByIds(ids []int) (bool, *exceptions.BaseErrorResponse)
	GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)

	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
}
