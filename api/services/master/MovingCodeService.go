package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
)

type MovingCodeService interface {
	GetAllMovingCode(pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	PushMovingCodePriority(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	CreateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetMovingCodebyId(Id int) (any, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusMovingCode(Id int) (any, *exceptionsss_test.BaseErrorResponse)
}
