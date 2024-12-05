package generalserviceapiutils

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

type ShippingMethodResponse struct {
	ShippingMethodId              int    `json:"shipping_method_id"`
	ShippingMethodCode            string `json:"shipping_method_code"`
	ShippingMethodName            string `json:"shipping_method_description"`
	ShippingMethodCodeDescription string `json:"shipping_method_code_description"`
}

type ShippingMethodParams struct {
	Page                          int    `json:"page"`
	Limit                         int    `json:"limit"`
	ShippingMethodId              int    `json:"shipping_method_id"`
	ShippingMethodCode            string `json:"shipping_method_code"`
	ShippingMethodName            string `json:"shipping_method_description"`
	ShippingMethodCodeDescription string `json:"shipping_method_code_description"`
	SortBy                        string `json:"sort_by"`
	SortOf                        string `json:"sort_of"`
}

func GetAllShippingMethod(params ShippingMethodParams) ([]ShippingMethodResponse, *exceptions.BaseErrorResponse) {
	var getShippingMethod []ShippingMethodResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "shipping-methods"

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

	err := utils.CallAPI("GET", url, nil, &getShippingMethod)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve shipping method due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "shipping method service is temporarily unavailable"
		}

		return getShippingMethod, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting shipping method by ID"),
		}
	}

	return getShippingMethod, nil
}

func GetShippingMethodByCode(code string) (ShippingMethodResponse, *exceptions.BaseErrorResponse) {
	var getShippingMethod ShippingMethodResponse
	url := config.EnvConfigs.GeneralServiceUrl + "shipping-method-code/" + code

	err := utils.CallAPI("GET", url, nil, &getShippingMethod)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve shipping method due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "shipping method service is temporarily unavailable"
		}

		return getShippingMethod, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting shipping method by ID"),
		}
	}
	return getShippingMethod, nil
}

func GetShippingMethodByName(name string) (ShippingMethodResponse, *exceptions.BaseErrorResponse) {
	var getShippingMethod ShippingMethodResponse
	url := config.EnvConfigs.GeneralServiceUrl + "shipping-method-name/" + name

	err := utils.CallAPI("GET", url, nil, &getShippingMethod)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve shipping method due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "shipping method service is temporarily unavailable"
		}

		return getShippingMethod, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting shipping method by ID"),
		}
	}
	return getShippingMethod, nil
}

func GetShippingMethodById(id int) (ShippingMethodResponse, *exceptions.BaseErrorResponse) {
	var getShippingMethod ShippingMethodResponse
	url := config.EnvConfigs.GeneralServiceUrl + "shipping-method/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getShippingMethod)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve shipping method due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "shipping method service is temporarily unavailable"
		}

		return getShippingMethod, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting shipping method by ID"),
		}
	}
	return getShippingMethod, nil
}

func GetShippingMethodByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)
	var nonZeroIds []string
	for _, id := range ids {
		if id != 0 {
			nonZeroIds = append(nonZeroIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(nonZeroIds, ",") + "]"
	url := config.EnvConfigs.GeneralServiceUrl + "shipping-method-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve shipping method due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "shipping method service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting shipping method by ID"),
		}
	}
	return nil
}
