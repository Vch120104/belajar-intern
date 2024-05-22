package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
)

type MovingCodeService interface {
	GetAllMovingCode(pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	PushMovingCodePriority(Id int) (bool, *exceptions.BaseErrorResponse)
	CreateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse)
	GetMovingCodebyId(Id int) (any, *exceptions.BaseErrorResponse)
	ChangeStatusMovingCode(Id int) (any, *exceptions.BaseErrorResponse)
}
