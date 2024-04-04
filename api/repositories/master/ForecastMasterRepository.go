package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ForecastMasterRepository interface {
	GetForecastMasterById(*gorm.DB, int) (masterpayloads.ForecastMasterResponse, *exceptionsss_test.BaseErrorResponse)
	SaveForecastMaster(*gorm.DB, masterpayloads.ForecastMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusForecastMaster(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllForecastMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
