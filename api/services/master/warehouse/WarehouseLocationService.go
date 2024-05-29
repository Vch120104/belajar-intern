package masterwarehouseservice

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
)

type WarehouseLocationService interface {
	Save(masterwarehousepayloads.GetWarehouseLocationResponse) (bool, *exceptions.BaseErrorResponse)
	GetById(int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptions.BaseErrorResponse)
	GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatus(int) (bool, *exceptions.BaseErrorResponse)
}
