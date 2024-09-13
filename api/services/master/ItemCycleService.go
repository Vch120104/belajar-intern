package masterservice

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
)

type ItemCycleService interface {
	ItemCycleInsert(payloads masterpayloads.ItemCycleInsertPayloads) (bool, *exceptions.BaseErrorResponse)
}
