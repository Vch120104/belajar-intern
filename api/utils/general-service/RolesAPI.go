package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type RoleResponse struct {
	RoleId   int    `json:"role_id"`
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
}

func GetRoleById(id int) (RoleResponse, *exceptions.BaseErrorResponse) {
	var role RoleResponse
	url := config.EnvConfigs.GeneralServiceUrl + "role/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &role)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve role due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "role service is temporarily unavailable"
		}

		return role, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting role by ID"),
		}
	}
	return role, nil
}

func GetRoleByCode(code string) (RoleResponse, *exceptions.BaseErrorResponse) {
	var role RoleResponse
	url := config.EnvConfigs.GeneralServiceUrl + "role-by-code/" + code
	err := utils.CallAPI("GET", url, nil, &role)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve role due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "role service is temporarily unavailable"
		}

		return role, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting role by ID"),
		}
	}
	return role, nil
}
