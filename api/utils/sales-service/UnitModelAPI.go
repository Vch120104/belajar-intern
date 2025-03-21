package salesserviceapiutils

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

type UnitModelResponse struct {
	ModelId              int    `json:"model_id"`
	ModelCode            string `json:"model_code"`
	ModelName            string `json:"model_description"`
	ModelCodeDescription string `json:"model_code_description"`
}

type UnitModelMultiIdResponse struct {
	ModelId          int    `json:"model_id"`
	ModelCode        string `json:"model_code"`
	ModelDescription string `json:"model_description"`
}

type UnitModelParams struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	ModelId          int    `json:"model_id"`
	ModelCode        string `json:"model_code"`
	ModelDescription string `json:"model_description"`
	SortBy           string `json:"sort_by"`
	SortOf           string `json:"sort_of"`
}

func GetAllUnitModel(params UnitModelParams) ([]UnitModelResponse, *exceptions.BaseErrorResponse) {
	var getUnitModel []UnitModelResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "unit-model-list"

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

	err := utils.CallAPI("GET", url, nil, &getUnitModel)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit model due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit model service is temporarily unavailable"
		}

		return getUnitModel, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit model by ID"),
		}
	}

	return getUnitModel, nil
}

func GetUnitModelByCode(code string) (UnitModelResponse, *exceptions.BaseErrorResponse) {
	var getUnitModel UnitModelResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-model-by-code/" + code
	err := utils.CallAPI("GET", url, nil, &getUnitModel)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit model due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit model service is temporarily unavailable"
		}

		return getUnitModel, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit model by Code"),
		}
	}
	return getUnitModel, nil
}

func GetUnitModelById(id int) (UnitModelResponse, *exceptions.BaseErrorResponse) {
	var getUnitModel UnitModelResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getUnitModel)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit model due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit model service is temporarily unavailable"
		}

		return getUnitModel, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit model by ID"),
		}
	}
	return getUnitModel, nil
}

func GetUnitModelByMultiId(ids []int) ([]UnitModelMultiIdResponse, *exceptions.BaseErrorResponse) {
	var getUnitModel []UnitModelMultiIdResponse

	ids = utils.RemoveDuplicateIds(ids)

	var strIds string
	for _, id := range ids {
		if id != 0 {
			strIds += strconv.Itoa(id) + ","
		}
	}
	if strIds != "" {
		strIds = "[" + strIds[:len(strIds)-1] + "]"
	} else {
		strIds = "[]"
	}

	url := config.EnvConfigs.SalesServiceUrl + "unit-model-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, &getUnitModel)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve model due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "model service is temporarily unavailable"
		}

		return getUnitModel, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting model by multi ID"),
		}
	}
	return getUnitModel, nil
}
