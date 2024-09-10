package masterrepository

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"gorm.io/gorm"
)

type ItemCycleRepository interface {
	InsertItemCycle(db *gorm.DB, payloads masterpayloads.ItemCycleInsertPayloads) (bool, *exceptions.BaseErrorResponse)
}
