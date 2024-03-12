package masterwarehouseservice

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
)

type WarehouseLocationService interface {
	Save(masterwarehousepayloads.GetWarehouseLocationResponse) (bool)
	GetById(int) (masterwarehousepayloads.GetWarehouseLocationResponse)
	GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) (pagination.Pagination)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseLocationResponse)
}
