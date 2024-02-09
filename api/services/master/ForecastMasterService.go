package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
)

type ForecastMasterService interface {
	GetForecastMasterById(int) masterpayloads.ForecastMasterResponse
}
