package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountRepository interface {
	GetAllDiscount(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetAllDiscountIsActive(*gorm.DB) ([]masterpayloads.DiscountResponse, error)
	GetDiscountById(*gorm.DB, int) (masterpayloads.DiscountResponse, error)
	GetDiscountByCode(*gorm.DB, string) (masterpayloads.DiscountResponse, error)
	SaveDiscount(*gorm.DB, masterpayloads.DiscountResponse) (bool, error)
	ChangeStatusDiscount(*gorm.DB, int) (bool, error)
}
