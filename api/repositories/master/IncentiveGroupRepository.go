package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type IncentiveGroupRepository interface {
	WithTrx(trxHandle *gorm.DB) IncentiveGroupRepository
	GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, error)
}
