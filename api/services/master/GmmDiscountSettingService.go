package masterservice

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
)

type GmmDiscountSettingService interface {
	GetAllGmmDiscountSetting() ([]masterpayloads.GmmDiscountSettingResponse, *exceptions.BaseErrorResponse)
}
