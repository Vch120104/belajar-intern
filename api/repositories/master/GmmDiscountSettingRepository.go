package masterrepository

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type GmmDiscountSettingRepository interface {
	GetAllGmmDiscountSetting(tx *gorm.DB) ([]masterpayloads.GmmDiscountSettingResponse, *exceptions.BaseErrorResponse)
}
