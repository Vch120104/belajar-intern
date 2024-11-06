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
		return unitColourResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching unit color by ID",
			Err:        errors.New("error consuming external API while fetching unit color by ID"),
		}
	}
	return unitColourResponse, nil
}
