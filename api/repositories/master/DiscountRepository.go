package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	exceptions "after-sales/api/exceptions"

	"gorm.io/gorm"
)

type DiscountRepository interface {
	GetAllDiscount(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllDiscountIsActive(*gorm.DB) ([]masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse)
	GetDiscountById(*gorm.DB, int) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse)
	GetDiscountByCode(*gorm.DB, string) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse)
	SaveDiscount(*gorm.DB, masterpayloads.DiscountResponse) (bool, *exceptions.BaseErrorResponse)
	UpdateDiscount(tx *gorm.DB, id int, req masterpayloads.DiscountUpdate) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusDiscount(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
