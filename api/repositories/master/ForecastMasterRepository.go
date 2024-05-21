package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ForecastMasterRepository interface {
	GetForecastMasterById(*gorm.DB, int) (masterpayloads.ForecastMasterResponse, *exceptions.BaseErrorResponse)
	SaveForecastMaster(*gorm.DB, masterpayloads.ForecastMasterResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusForecastMaster(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	GetAllForecastMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
