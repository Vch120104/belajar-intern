package masteritemservice

import (
	"after-sales/api/payloads/pagination"
)

type ItemPackageDetailService interface {
	GetItemPackageDetailByItemPackageId(itemPackageId int, pages pagination.Pagination) pagination.Pagination
}
