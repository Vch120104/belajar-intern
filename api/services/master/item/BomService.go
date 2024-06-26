package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type BomService interface {
	GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetBomMasterById(id int) (masteritempayloads.BomMasterRequest, *exceptions.BaseErrorResponse)
	SaveBomMaster(request masteritempayloads.BomMasterRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateBomMaster(id int, request masteritempayloads.BomMasterRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusBomMaster(id int) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetBomDetailById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	SaveBomDetail(request masteritempayloads.BomDetailRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateBomDetail(id int, request masteritempayloads.BomDetailRequest) (bool, *exceptions.BaseErrorResponse)
	DeleteByIds(ids []int) (bool, *exceptions.BaseErrorResponse)
	GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
