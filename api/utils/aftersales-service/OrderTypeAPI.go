package aftersalesserviceapiutils

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

type OrderTypeResponse struct {
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeCode string `json:"order_type_code"`
	OrderTypeName string `json:"order_type_name"`
}

type OrderTypeParams struct {
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeCode string `json:"order_type_code"`
	OrderTypeName string `json:"order_type_name"`
	SortBy        string `json:"sort_by"`
	SortOf        string `json:"sort_of"`
}

func GetAllOrderType(params OrderTypeParams) ([]OrderTypeResponse, *exceptions.BaseErrorResponse) {
	var getOrderType []OrderTypeResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.AfterSalesServiceUrl + "order-type"

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

	err := utils.CallAPI("GET", url, nil, &getOrderType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "order type service is temporarily unavailable"
		}

		return getOrderType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting order type by ID"),
		}
	}

	return getOrderType, nil
}

func GetOrderTypeByCode(code string) (OrderTypeResponse, *exceptions.BaseErrorResponse) {
	var getOrderType OrderTypeResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "order-type-code/" + code

	err := utils.CallAPI("GET", url, nil, &getOrderType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "order type service is temporarily unavailable"
		}

		return getOrderType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting order type by ID"),
		}
	}
	return getOrderType, nil
}

func GetOrderTypeByName(name string) (OrderTypeResponse, *exceptions.BaseErrorResponse) {
	var getOrderType OrderTypeResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "order-type-name/" + name

	err := utils.CallAPI("GET", url, nil, &getOrderType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "order type service is temporarily unavailable"
		}

		return getOrderType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting order type by ID"),
		}
	}
	return getOrderType, nil
}

func GetOrderTypeById(id int) (OrderTypeResponse, *exceptions.BaseErrorResponse) {
	var getOrderType OrderTypeResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "order-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getOrderType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "order type service is temporarily unavailable"
		}

		return getOrderType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting order type by ID"),
		}
	}
	return getOrderType, nil
}

func GetOrderTypeByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)
	var nonZeroIds []string
	for _, id := range ids {
		if id != 0 {
			nonZeroIds = append(nonZeroIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(nonZeroIds, ",") + "]"
	url := config.EnvConfigs.AfterSalesServiceUrl + "order-type-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "order type service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting order type by ID"),
		}
	}
	return nil
}
