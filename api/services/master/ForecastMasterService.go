package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ForecastMasterService interface {
	GetForecastMasterById(int) (masterpayloads.ForecastMasterResponse,*exceptionsss_test.BaseErrorResponse)
	SaveForecastMaster(masterpayloads.ForecastMasterResponse) (bool,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusForecastMaster(int) (bool,*exceptionsss_test.BaseErrorResponse)
	GetAllForecastMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int,*exceptionsss_test.BaseErrorResponse)
}
