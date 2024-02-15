package masterwarehouserepository

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	"gorm.io/gorm"
)

type WarehouseGroupRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseGroupResponse) (bool, error)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseGroupResponse, error)
	GetAll(*gorm.DB, masterwarehousepayloads.GetAllWarehouseGroupRequest) ([]masterwarehousepayloads.GetWarehouseGroupResponse, error)
	ChangeStatus(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseGroupResponse, error)
}
