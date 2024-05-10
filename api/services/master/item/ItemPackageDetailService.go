package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
)

type ItemPackageDetailService interface {
	GetItemPackageDetailByItemPackageId(itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)

	CreateItemPackageDetailByItemPackageId(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateItemPackageDetailByItemPackageId(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptionsss_test.BaseErrorResponse)
	GetItemPackageDetailById(itemPackageDetailId int) (masteritempayloads.ItemPackageDetailResponse, *exceptionsss_test.BaseErrorResponse)

	ChangeStatusItemPackageDetail(id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
