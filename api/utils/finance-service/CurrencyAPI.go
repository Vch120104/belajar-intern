package financeserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type CurrencyResponse struct {
	CurrencyId   int    `json:"currency_id"`
	CurrencyCode string `json:"currency_code"`
	CurrencyName string `json:"currency_name"`
}

func GetCurrencyId(id int) (CurrencyResponse, *exceptions.BaseErrorResponse) {
	var getCurrency CurrencyResponse
	url := config.EnvConfigs.FinanceServiceUrl + "currency-code/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getCurrency)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve currency due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "currency service is temporarily unavailable"
		}

		return getCurrency, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting currency by ID"),
		}
	}
	return getCurrency, nil
}
