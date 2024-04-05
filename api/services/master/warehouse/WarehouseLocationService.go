package masterwarehouseservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
)

type WarehouseLocationService interface {
	Save(masterwarehousepayloads.GetWarehouseLocationResponse) (bool,*exceptionsss_test.BaseErrorResponse)
	GetById(int) (masterwarehousepayloads.GetWarehouseLocationResponse,*exceptionsss_test.BaseErrorResponse)
	GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) (pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	ChangeStatus(int) (bool,*exceptionsss_test.BaseErrorResponse)
}
