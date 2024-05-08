package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type MovingCodeRepository interface {
	GetAllMovingCode(tx *gorm.DB, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	PushMovingCodePriority(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	CreateMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeListRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeListRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetMovingCodebyId(tx *gorm.DB, Id int) (any, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusMovingCode(tx *gorm.DB, Id int) (any, *exceptionsss_test.BaseErrorResponse)
}
