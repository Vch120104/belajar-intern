package financeserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type CurrencyParams struct {
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	CurrencyCode string `json:"currency_code"`
}

type CurrencyResponse struct {
	CurrencyId   int    `json:"currency_id"`
	CurrencyCode string `json:"currency_code"`
	CurrencyName string `json:"currency_name"`
}

func GetAllCurrency(params CurrencyParams) ([]CurrencyResponse, *exceptions.BaseErrorResponse) {
	var getCurrency []CurrencyResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.FinanceServiceUrl + "currency-code"

	queryParams := fmt.Sprintf("page=%d&limit=%d", params.Page, params.Limit)

	v := reflect.ValueOf(params)
	typeOfParams := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		if strVal, ok := value.(string); ok && strVal != "" {
			key := typeOfParams.Field(i).Tag.Get("json")
			value := strings.ReplaceAll(strVal, " ", "%20")
			queryParams += "&" + key + "=" + value
		}
	}

	url := baseURL + "?" + queryParams

	fmt.Println(url)
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
