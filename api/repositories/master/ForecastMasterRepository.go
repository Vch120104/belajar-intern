package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type ForecastMasterRepository interface {
	GetForecastMasterById(*gorm.DB, int) (masterpayloads.ForecastMasterResponse, error)
	SaveForecastMaster(*gorm.DB, masterpayloads.ForecastMasterResponse) (bool, error)
	ChangeStatusForecastMaster(*gorm.DB, int) (bool, error)
}
