package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type LineTypeResponse struct {
	LineTypeId   int    `json:"line_type_id"`
	LineTypeCode string `json:"line_type_code"`
	LineTypeName string `json:"line_type_name"`
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
