package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type MovingCodeRepository interface {
	GetAllMovingCode(tx *gorm.DB, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	PushMovingCodePriority(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	CreateMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse)
	GetMovingCodebyId(tx *gorm.DB, Id int) (any, *exceptions.BaseErrorResponse)
	ChangeStatusMovingCode(tx *gorm.DB, Id int) (any, *exceptions.BaseErrorResponse)
}
