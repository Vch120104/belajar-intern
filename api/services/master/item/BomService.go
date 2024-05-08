package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type BomService interface {
	GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetBomMasterById(id int) (masteritempayloads.BomMasterRequest, *exceptionsss_test.BaseErrorResponse)
	SaveBomMaster(masteritempayloads.BomMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusBomMaster(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetBomDetailById(id int) ([]masteritempayloads.BomDetailListResponse, *exceptionsss_test.BaseErrorResponse)
	GetBomDetailByIds(id int) ([]masteritempayloads.BomDetailListResponse, *exceptionsss_test.BaseErrorResponse)
	SaveBomDetail(masteritempayloads.BomDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	DeleteByIds(ids []int) (bool, *exceptionsss_test.BaseErrorResponse)
}
