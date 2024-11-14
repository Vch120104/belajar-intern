package financeserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
)

type OpenPeriodPayloadResponse struct {
	PeriodYear        string `json:"period_year"`
	PeriodMonth       string `json:"period_month"`
	CurrentPeriodDate string `json:"current_period_date"`
}

func GetOpenPeriodByCompany(companyId int, moduleCode string) (OpenPeriodPayloadResponse, *exceptions.BaseErrorResponse) {
	var PeriodResponse OpenPeriodPayloadResponse
	PeriodUrl := fmt.Sprintf("%sclosing-period-company/current-period?company_id=%s&closing_module_detail_code=%s", config.EnvConfigs.FinanceServiceUrl, strconv.Itoa(companyId), moduleCode)
	if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
		return PeriodResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to Period Response data from external service",
			Err:        err,
		}
	}
	return PeriodResponse, nil
}
