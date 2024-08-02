package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
)

type MovingCodeService interface {
	GetAllMovingCode(companyId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	PushMovingCodePriority(companyId int, Id int) (bool, *exceptions.BaseErrorResponse)
	CreateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateMovingCode(req masterpayloads.MovingCodeListUpdate) (bool, *exceptions.BaseErrorResponse)
	GetMovingCodebyId(Id int) (masterpayloads.MovingCodeResponse, *exceptions.BaseErrorResponse)
	ChangeStatusMovingCode(Id int) (any, *exceptions.BaseErrorResponse)
	GetDropdownMovingCode(companyId int) ([]masterpayloads.MovingCodeDropDown, *exceptions.BaseErrorResponse)
	DeactiveMovingCode(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateMovingCode(id string) (bool, *exceptions.BaseErrorResponse)
}
