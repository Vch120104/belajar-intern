package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ForecastMasterService interface {
	GetForecastMasterById(int) masterpayloads.ForecastMasterResponse
	SaveForecastMaster(masterpayloads.ForecastMasterResponse) bool
	ChangeStatusForecastMaster(int) bool
	GetAllForecastMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
}
