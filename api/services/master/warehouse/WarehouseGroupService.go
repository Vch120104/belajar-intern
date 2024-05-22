package masterwarehouseservice

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarehouseGroupService interface {
	SaveWarehouseGroup(masterwarehousepayloads.GetWarehouseGroupResponse) (bool, *exceptions.BaseErrorResponse)
	GetByIdWarehouseGroup(int) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptions.BaseErrorResponse)
	GetAllWarehouseGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatusWarehouseGroup(int) (bool, *exceptions.BaseErrorResponse)
}
