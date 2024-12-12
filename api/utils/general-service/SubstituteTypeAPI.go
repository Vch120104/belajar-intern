package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type SubstituteTypeResponse struct {
	SubstituteTypeId   int    `json:"substitute_type_id"`
	SubstituteTypeCode string `json:"substitute_type_code"`
	SubstituteTypeName string `json:"substitute_type_name"`
}

func GetAllSubstituteType() ([]SubstituteTypeResponse, *exceptions.BaseErrorResponse) {
	var getSubstituteType []SubstituteTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "substitute-types"

	err := utils.CallAPI("GET", url, nil, &getSubstituteType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve substitute type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "substitute type service is temporarily unavailable"
		}

		return getSubstituteType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting substitute types"),
		}
	}

	return getSubstituteType, nil
}

func GetSubstituteTypeById(id int) (SubstituteTypeResponse, *exceptions.BaseErrorResponse) {
	var getSubstituteType SubstituteTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "substitute-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getSubstituteType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve substitute type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "substitute type service is temporarily unavailable"
		}

		return getSubstituteType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting substitute type by ID"),
		}
	}

	return getSubstituteType, nil
}

func GetSubstituteTypeByCode(code string) (SubstituteTypeResponse, *exceptions.BaseErrorResponse) {
	var getSubstituteType SubstituteTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "substitute-type/code/" + code

	err := utils.CallAPI("GET", url, nil, &getSubstituteType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve substitute type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "substitute type service is temporarily unavailable"
		}

		return getSubstituteType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting substitute type by code"),
		}
	}

	return getSubstituteType, nil
}
