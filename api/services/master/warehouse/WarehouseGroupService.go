package masterwarehouseservice

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
)

type WarehouseGroupService interface {
	Save(masterwarehousepayloads.GetWarehouseGroupResponse) bool
	GetById(int) masterwarehousepayloads.GetWarehouseGroupResponse
	GetAll(request masterwarehousepayloads.GetAllWarehouseGroupRequest) []masterwarehousepayloads.GetWarehouseGroupResponse
	ChangeStatus(int) masterwarehousepayloads.GetWarehouseGroupResponse
}
