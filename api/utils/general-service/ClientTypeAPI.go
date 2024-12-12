package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type ClientTypeResponse struct {
	ClientTypeId          int    `json:"client_type_id"`
	ClientTypeCode        string `json:"client_type_code"`
	ClientTypeDescription string `json:"client_type_description"`
	ClientFlagListId      int    `json:"client_flag_list_id"`
	ClientGroupId         int    `json:"client_group_id"`
	IsActive              bool   `json:"is_active"`
}

func GetClientTypeById(id int) (ClientTypeResponse, *exceptions.BaseErrorResponse) {
	var client ClientTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "client-type/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &client)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve client type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "client type service is temporarily unavailable"
		}

		return client, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting client type by ID"),
		}
	}
	return client, nil
}

func GetClientTypeByCode(code string) (ClientTypeResponse, *exceptions.BaseErrorResponse) {
	var client ClientTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "client-type-code/" + code
	err := utils.CallAPI("GET", url, nil, &client)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve client type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "client type service is temporarily unavailable"
		}

		return client, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting client type by ID"),
		}
	}
	return client, nil
}
