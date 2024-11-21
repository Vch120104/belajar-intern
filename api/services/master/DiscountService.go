package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DiscountService interface {
	GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllDiscountIsActive() ([]masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse)
	GetDiscountById(Id int) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse)
	GetDiscountByCode(Code string) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse)
	SaveDiscount(req masterpayloads.DiscountResponse) (bool, *exceptions.BaseErrorResponse)
	UpdateDiscount(id int, req masterpayloads.DiscountUpdate) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusDiscount(Id int) (bool, *exceptions.BaseErrorResponse)
}
