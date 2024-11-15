package financeserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
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
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve closing period due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "closing period service is temporarily unavailable"
		}

		return PeriodResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting closing period by Company ID"),
		}
	}
	return PeriodResponse, nil
}
