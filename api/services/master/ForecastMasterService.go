package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ForecastMasterService interface {
	GetForecastMasterById(int) (masterpayloads.ForecastMasterResponse, *exceptions.BaseErrorResponse)
	SaveForecastMaster(masterpayloads.ForecastMasterResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusForecastMaster(int) (bool, *exceptions.BaseErrorResponse)
	GetAllForecastMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateForecastMaster(req masterpayloads.ForecastMasterResponse, id int)(bool,*exceptions.BaseErrorResponse)
}
