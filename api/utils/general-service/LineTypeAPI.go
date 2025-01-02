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

type LineTypeResponse struct {
	IsActive     bool   `json:"is_active"`
	LineTypeId   int    `json:"line_type_id"`
	LineTypeCode string `json:"line_type_code"`
	LineTypeName string `json:"line_type_name"`
}

type LineTypeListParams struct {
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	LineTypeCode string `json:"line_type_code"`
	LineTypeName string `json:"line_type_name"`
	IsActive     string `json:"is_active"`
	SortBy       string `json:"sort_by"`
	SortOf       string `json:"sort_of"`
}

func GetLineTypeById(id int) (LineTypeResponse, *exceptions.BaseErrorResponse) {
	var line LineTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &line)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve line type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "line type service is temporarily unavailable"
		}

		return line, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting line type by ID"),
		}
	}
	return line, nil
}

func GetLineTypeByCode(code string) (LineTypeResponse, *exceptions.BaseErrorResponse) {
	var line LineTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "line-type-code/" + code
	err := utils.CallAPI("GET", url, nil, &line)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve line type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "line type service is temporarily unavailable"
		}

		return line, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting line type by ID"),
		}
	}
	return line, nil
}

func GetLineTypeListByCode(param LineTypeListParams) ([]LineTypeResponse, *exceptions.BaseErrorResponse) {
	var line []LineTypeResponse
	if param.Limit == 0 {
		param.Limit = 1000000
	}

	// Make URL
	baseUrl := config.EnvConfigs.GeneralServiceUrl + "line-type-list?"
	queryParam := fmt.Sprintf("page=%d&limit=%d", param.Page, param.Limit)

	v := reflect.ValueOf(param)
	typeOfParam := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		strval, ok := value.(string)
		if ok && strval != "" {
			key := typeOfParam.Field(i).Tag.Get("json")
			value := strings.ReplaceAll(strval, " ", "%20")
			queryParam += "&" + key + "=" + value
		}
	}

	url := baseUrl + queryParam

	// Call external API
	err := utils.CallAPI("GET", url, nil, &line)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve line type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "line type service is temporarily unavailable"
		}

		return line, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting line type by ID"),
		}
	}
	return line, nil
}
