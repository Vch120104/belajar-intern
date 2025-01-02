package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type PDIStatusResponse struct {
	PDIStatusId   int    `json:"pdi_status_id"`
	PDIStatusCode string `json:"pdi_status_code"`
	PDIStatusName string `json:"pdi_status_description"`
}

func GetPDIStatusById(id int) (PDIStatusResponse, *exceptions.BaseErrorResponse) {
	var pdiStatus PDIStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "pdi-status/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &pdiStatus)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve pdi status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "PDI status service is temporarily unavailable"
		}

		return pdiStatus, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting pdi status by ID"),
		}
	}
	return pdiStatus, nil
}

func GetPDIStatusByCode(code string) (PDIStatusResponse, *exceptions.BaseErrorResponse) {
	var pdiStatus PDIStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "pdi-status-code/" + code

	err := utils.CallAPI("GET", url, nil, &pdiStatus)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve pdi status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "PDI status service is temporarily unavailable"
		}

		return pdiStatus, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting pdi status by ID"),
		}
	}
	return pdiStatus, nil
}
