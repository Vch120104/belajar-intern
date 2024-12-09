package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type UnitColourResponse []struct {
	VariantColourId   int    `json:"colour_id"`
	VariantColourCode string `json:"colour_commercial_name"`
	VariantColourName string `json:"colour_police_name"`
}

func GetUnitColourByID(id int) (UnitColourResponse, *exceptions.BaseErrorResponse) {
	var unitColourResponse UnitColourResponse

	if id <= 0 {
		return unitColourResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid colour ID provided",
			Err:        errors.New("invalid colour ID provided"),
		}
	}

	url := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &unitColourResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit colour due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit colour service is temporarily unavailable"
		}

		return unitColourResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit colour by ID"),
		}
	}
	return unitColourResponse, nil
}

func GetUnitColourByBrandId(id int) (UnitColourResponse, *exceptions.BaseErrorResponse) {
	var unitColourResponse UnitColourResponse

	if id <= 0 {
		return unitColourResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid colour ID provided",
			Err:        errors.New("invalid colour ID provided"),
		}
	}

	url := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &unitColourResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit colour due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit colour service is temporarily unavailable"
		}

		return unitColourResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit colour by Brand ID"),
		}
	}
	return unitColourResponse, nil
}
