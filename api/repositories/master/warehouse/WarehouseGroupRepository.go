package masterwarehouserepository

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	"gorm.io/gorm"
)

type WarehouseGroupRepository interface {
	WithTrx(trxHandle *gorm.DB) WarehouseGroupRepository
	Save(masterwarehousepayloads.GetWarehouseGroupResponse) (bool, error)
	GetById(int) (masterwarehousepayloads.GetWarehouseGroupResponse, error)
	GetAll(request masterwarehousepayloads.GetAllWarehouseGroupRequest) ([]masterwarehousepayloads.GetWarehouseGroupResponse, error)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseGroupResponse, error)
}
