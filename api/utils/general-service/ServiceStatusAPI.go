package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type ServiceStatusResponse struct {
	ServiceStatusId          int    `json:"service_status_id"`
	ServiceStatusCode        string `json:"service_status_code"`
	ServiceStatusDescription string `json:"service_status_description"`
	IsActive                 bool   `json:"is_active"`
}

func GetServiceStatusById(id int) (ServiceStatusResponse, *exceptions.BaseErrorResponse) {
	var serviceStatus ServiceStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "service-status/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &serviceStatus)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve service status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Service status service is temporarily unavailable"
		}

		return serviceStatus, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting service status by ID"),
		}
	}
	return serviceStatus, nil
}
