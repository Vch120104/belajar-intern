package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type UnitColourData struct {
	ColourId             int    `json:"colour_id"`
	ColourCommercialName string `json:"colour_commercial_name"`
	ColourPoliceName     string `json:"colour_police_name"`
}

type UnitColourResponse struct {
	StatusCode int              `json:"status_code"`
	Message    string           `json:"message"`
	Data       []UnitColourData `json:"data"`
}

type UnitColourDetailData struct {
	BrandId              int    `json:"brand_id"`
	ColourCode           string `json:"colour_code"`
	ColourCommercialName string `json:"colour_commercial_name"`
	ColourPoliceName     string `json:"colour_police_name"`
	ColourId             int    `json:"colour_id"`
	IsActive             bool   `json:"is_active"`
	BrandName            string `json:"brand_name"`
}

func GetUnitColorById(colourId int) (UnitColourDetailData, *exceptions.BaseErrorResponse) {
	var unitColourDetailResponse UnitColourDetailData

	if colourId <= 0 {
		return unitColourDetailResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid colour ID provided",
			Err:        errors.New("invalid colour ID provided"),
		}
	}

	url := config.EnvConfigs.SalesServiceUrl + "unit-colour/" + strconv.Itoa(colourId)
	//fmt.Println("Requesting URL:", url)

	err := utils.CallAPI("GET", url, nil, &unitColourDetailResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit colour details due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit colour service is temporarily unavailable"
		}

		return unitColourDetailResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit colour details"),
		}
	}

	return unitColourDetailResponse, nil
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
	fmt.Println("Requesting URL:", url)

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

	// Validasi status_code dari respons
	if unitColourResponse.StatusCode != http.StatusOK {
		return unitColourResponse, &exceptions.BaseErrorResponse{
			StatusCode: unitColourResponse.StatusCode,
			Message:    unitColourResponse.Message,
			Err:        errors.New("unexpected response status code from API"),
		}
	}

	return unitColourResponse, nil
}
