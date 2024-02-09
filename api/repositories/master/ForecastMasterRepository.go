package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type ForecastMasterRepository interface {
	GetForecastMasterById(*gorm.DB, int) (masterpayloads.ForecastMasterResponse, error)
}
