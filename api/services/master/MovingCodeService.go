package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	// "after-sales/api/utils"
)

type MovingCodeService interface {
	GetAllMovingCode(pagination.Pagination) pagination.Pagination
	SaveMovingCode(masterpayloads.MovingCodeRequest) bool
	ChangePriorityMovingCode(int) bool
	ChangeStatusMovingCode(int) bool
}
