package masterwarehouseservice

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	"gorm.io/gorm"
)

type WarehouseGroupService interface {
	WithTrx(Trxhandle *gorm.DB) WarehouseGroupService
	Save(masterwarehousepayloads.GetWarehouseGroupResponse) (bool, error)
	GetById(int) (masterwarehousepayloads.GetWarehouseGroupResponse, error)
	GetAll(request masterwarehousepayloads.GetAllWarehouseGroupRequest) ([]masterwarehousepayloads.GetWarehouseGroupResponse, error)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseGroupResponse, error)
}
