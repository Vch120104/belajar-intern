package masterservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
)

type LookupService interface {
	ItemOprCode(linetypeId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
