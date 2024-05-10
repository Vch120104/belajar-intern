package masterwarehouseservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarehouseGroupService interface {
	SaveWarehouseGroup(masterwarehousepayloads.GetWarehouseGroupResponse) (bool,*exceptionsss_test.BaseErrorResponse)
	GetByIdWarehouseGroup(int) (masterwarehousepayloads.GetWarehouseGroupResponse,*exceptionsss_test.BaseErrorResponse)
	GetAllWarehouseGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusWarehouseGroup(int) (bool,*exceptionsss_test.BaseErrorResponse)
}
