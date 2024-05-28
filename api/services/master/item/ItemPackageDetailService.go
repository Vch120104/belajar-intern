package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
)

type ItemPackageDetailService interface {
	GetItemPackageDetailByItemPackageId(itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)

	CreateItemPackageDetailByItemPackageId(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse)
	UpdateItemPackageDetailByItemPackageId(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse)
	GetItemPackageDetailById(itemPackageDetailId int) (masteritempayloads.ItemPackageDetailResponse, *exceptions.BaseErrorResponse)

	ChangeStatusItemPackageDetail(id int) (bool, *exceptions.BaseErrorResponse)
}
