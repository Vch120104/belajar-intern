package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type IncentiveGroupService interface {
	WithTrx(trxHandle *gorm.DB) IncentiveGroupService
	GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, error)
}
