package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

)

type DiscountService interface {
	GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetAllDiscountIsActive() []masterpayloads.DiscountResponse
	GetDiscountById(Id int) masterpayloads.DiscountResponse
	GetDiscountByCode(Code string) masterpayloads.DiscountResponse
	SaveDiscount(req masterpayloads.DiscountResponse) bool
	ChangeStatusDiscount(Id int) bool
}
