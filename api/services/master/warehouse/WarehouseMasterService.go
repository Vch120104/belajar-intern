package masterwarehouseservice

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseMasterService interface {
	WithTrx(Trxhandle *gorm.DB) WarehouseMasterService
	Save(masterwarehousepayloads.GetWarehouseMasterResponse) (bool, error)
	GetById(int) (masterwarehousepayloads.GetWarehouseMasterResponse, error)
	GetAll(request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) (pagination.Pagination, error)
	GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, error)
	GetWarehouseMasterByCode(Code string) ([]map[string]interface{}, error)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseMasterResponse, error)
}
