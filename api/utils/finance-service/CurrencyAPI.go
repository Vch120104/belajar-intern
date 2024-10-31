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

func GetCurrency(id int) (CurrencyResponse, *exceptions.BaseErrorResponse) {
	var getCurrency CurrencyResponse
	url := config.EnvConfigs.FinanceServiceUrl + "currency/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getCurrency)
	if err != nil {
		return getCurrency, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consume data currency external api",
			Err:        errors.New("error consume data currency  external api"),
		}
	}
	return getCurrency, nil
}
