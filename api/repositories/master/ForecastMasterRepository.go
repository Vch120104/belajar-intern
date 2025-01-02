package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ForecastMasterRepository interface {
	GetForecastMasterById(*gorm.DB, int) (masterpayloads.ForecastMasterResponse, *exceptions.BaseErrorResponse)
	SaveForecastMaster(*gorm.DB, masterpayloads.ForecastMasterResponse) (masterentities.ForecastMaster, *exceptions.BaseErrorResponse)
	ChangeStatusForecastMaster(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	GetAllForecastMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdateForecastMaster(tx *gorm.DB, req masterpayloads.ForecastMasterResponse, id int) (masterentities.ForecastMaster, *exceptions.BaseErrorResponse)
}
