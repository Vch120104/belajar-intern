package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountRepository interface {
	WithTrx(trxHandle *gorm.DB) DiscountRepository
	GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	GetAllDiscountIsActive() ([]masterpayloads.DiscountResponse, error)
	GetDiscountById(Id int) (masterpayloads.DiscountResponse, error)
	GetDiscountByCode(Code string) (masterpayloads.DiscountResponse, error)
	SaveDiscount(req masterpayloads.DiscountResponse) (bool, error)
	ChangeStatusDiscount(Id int) (bool, error)
}
