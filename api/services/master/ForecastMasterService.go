package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ForecastMasterService interface {
	GetForecastMasterById(int) (masterpayloads.ForecastMasterResponse, *exceptions.BaseErrorResponse)
	SaveForecastMaster(masterpayloads.ForecastMasterResponse) (masterentities.ForecastMaster, *exceptions.BaseErrorResponse)
	ChangeStatusForecastMaster(int) (bool, *exceptions.BaseErrorResponse)
	GetAllForecastMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdateForecastMaster(req masterpayloads.ForecastMasterResponse, id int) (masterentities.ForecastMaster, *exceptions.BaseErrorResponse)
}
