package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type LookupRepository interface {
	ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
