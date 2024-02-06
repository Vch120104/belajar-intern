package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type FieldActionService interface {
	GetAllFieldAction(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	SaveFieldAction(req masterpayloads.FieldActionResponse) bool
	// GetFieldActionById(id int) masterpayloads.FieldActionResponse

}
