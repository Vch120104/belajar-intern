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

type LandedCostTypeResponse struct {
	LandedCostTypeId          int    `json:"landed_cost_type_id"`
	LandedCostTypeCode        string `json:"landed_cost_type_code"`
	LandedCostTypeName        string `json:"landed_cost_type_name"`
	LandedCostTypeDescription string `json:"landed_cost_type_description"`
}

type LandedCostTypeParams struct {
	Page                      int    `json:"page"`
	Limit                     int    `json:"limit"`
	LandedCostTypeId          int    `json:"landed_cost_type_id"`
	LandedCostTypeCode        string `json:"landed_cost_type_code"`
	LandedCostTypeName        string `json:"landed_cost_type_name"`
	LandedCostTypeDescription string `json:"landed_cost_type_description"`
	SortBy                    string `json:"sort_by"`
	SortOf                    string `json:"sort_of"`
}

func GetAllLandedCostType(params LandedCostTypeParams) ([]LandedCostTypeResponse, *exceptions.BaseErrorResponse) {
	var getLandedCostType []LandedCostTypeResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "landed-cost-types"

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

	err := utils.CallAPI("GET", url, nil, &getLandedCostType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve landed cost type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "landed cost type service is temporarily unavailable"
		}

		return getLandedCostType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting landed cost type by ID"),
		}
	}

	return getLandedCostType, nil
}

func GetLandedCostTypeByCode(code string) (LandedCostTypeResponse, *exceptions.BaseErrorResponse) {
	var getLandedCostType LandedCostTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "landed-cost-type/" + code

	err := utils.CallAPI("GET", url, nil, &getLandedCostType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve landed cost type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "landed cost type service is temporarily unavailable"
		}

		return getLandedCostType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting landed cost type by ID"),
		}
	}
	return getLandedCostType, nil
}

func GetLandedCostTypeByName(name string) (LandedCostTypeResponse, *exceptions.BaseErrorResponse) {
	var getLandedCostType LandedCostTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "landed-cost-type/" + name

	err := utils.CallAPI("GET", url, nil, &getLandedCostType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve landed cost type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "landed cost type service is temporarily unavailable"
		}

		return getLandedCostType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting landed cost type by ID"),
		}
	}
	return getLandedCostType, nil
}

func GetLandedCostTypeById(id int) (LandedCostTypeResponse, *exceptions.BaseErrorResponse) {
	var getLandedCostType LandedCostTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "landed-cost-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getLandedCostType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve landed cost type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "landed cost type service is temporarily unavailable"
		}

		return getLandedCostType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting landed cost type by ID"),
		}
	}
	return getLandedCostType, nil
}

func GetLandedCostTypeByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)
	var nonZeroIds []string
	for _, id := range ids {
		if id != 0 {
			nonZeroIds = append(nonZeroIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(nonZeroIds, ",") + "]"
	url := config.EnvConfigs.GeneralServiceUrl + "landed-cost-type-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve landed cost type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "landed cost type service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting landed cost type by ID"),
		}
	}
	return nil
}
