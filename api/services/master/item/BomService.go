package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type BomService interface {
	GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
	GetBomMasterById(id int) masteritempayloads.BomMasterRequest
	SaveBomMaster(masteritempayloads.BomMasterRequest) bool
	ChangeStatusBomMaster(Id int) bool
	GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
}
