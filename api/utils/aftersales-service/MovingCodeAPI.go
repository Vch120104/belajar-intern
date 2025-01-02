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

type MovingCodeResponse struct {
	MovingCodeId   int    `json:"moving_code_id"`
	MovingCode     string `json:"moving_code"`
	MovingCodeName string `json:"moving_code_description"`
}

type MovingCodeParams struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	MovingCodeId   int    `json:"moving_code_id"`
	MovingCode     string `json:"moving_code"`
	MovingCodeName string `json:"moving_code_description"`
	SortBy         string `json:"sort_by"`
	SortOf         string `json:"sort_of"`
}

func GetAllMovingCode(params MovingCodeParams) ([]MovingCodeResponse, *exceptions.BaseErrorResponse) {
	var getMovingCode []MovingCodeResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.AfterSalesServiceUrl + "moving-code"

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

	err := utils.CallAPI("GET", url, nil, &getMovingCode)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve moving code due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "moving code service is temporarily unavailable"
		}

		return getMovingCode, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting moving code by ID"),
		}
	}

	return getMovingCode, nil
}

func GetMovingCodeByCode(code string) (MovingCodeResponse, *exceptions.BaseErrorResponse) {
	var getMovingCode MovingCodeResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "moving-code-code/" + code

	err := utils.CallAPI("GET", url, nil, &getMovingCode)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve moving code due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "moving code service is temporarily unavailable"
		}

		return getMovingCode, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting moving code by ID"),
		}
	}
	return getMovingCode, nil
}

func GetMovingCodeByName(name string) (MovingCodeResponse, *exceptions.BaseErrorResponse) {
	var getMovingCode MovingCodeResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "moving-code-name/" + name

	err := utils.CallAPI("GET", url, nil, &getMovingCode)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve moving code due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "moving code service is temporarily unavailable"
		}

		return getMovingCode, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting moving code by ID"),
		}
	}
	return getMovingCode, nil
}

func GetMovingCodeById(id int) (MovingCodeResponse, *exceptions.BaseErrorResponse) {
	var getMovingCode MovingCodeResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "moving-code/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getMovingCode)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve moving code due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "moving code service is temporarily unavailable"
		}

		return getMovingCode, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting moving code by ID"),
		}
	}
	return getMovingCode, nil
}

func GetMovingCodeByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)
	var nonZeroIds []string
	for _, id := range ids {
		if id != 0 {
			nonZeroIds = append(nonZeroIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(nonZeroIds, ",") + "]"
	url := config.EnvConfigs.AfterSalesServiceUrl + "moving-code-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve moving code due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "moving code service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting moving code by ID"),
		}
	}
	return nil
}
