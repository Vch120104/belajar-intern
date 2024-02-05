package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ForecastMasterRepository interface {
	GetForecastMasterById(*gorm.DB, int) (masterpayloads.ForecastMasterResponse, error)
	SaveForecastMaster(*gorm.DB, masterpayloads.ForecastMasterResponse) (bool, error)
	ChangeStatusForecastMaster(*gorm.DB, int) (bool, error)
	GetAllForecastMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error)
}
