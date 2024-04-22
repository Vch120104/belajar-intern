package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type MovingCodeRepository interface {
	GetAllMovingCode(*gorm.DB, pagination.Pagination) (pagination.Pagination, error)
	GetMovingCodeById(*gorm.DB, int) (masterpayloads.MovingCodeResponse, error)
	GetMovingCodeByPriority(*gorm.DB, float64) (masterpayloads.MovingCodeResponse, error)
	SaveMovingCode(*gorm.DB, masterpayloads.MovingCodeRequest) (bool, error)
	IncreasePriorityMovingCode(*gorm.DB, int) (bool, error)
	DecreasePriorityMovingCode(*gorm.DB, int) (bool, error)
	ChangeStatusMovingCode(*gorm.DB, int) (bool, error)
}
