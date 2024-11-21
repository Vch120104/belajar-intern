package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type UserDetailsResponse struct {
	IsActive         bool   `json:"is_active"`
	UserEmployeeId   int    `json:"user_employee_id"`
	UserId           int    `json:"user_id"`
	EmployeeName     string `json:"employee_name"`
	EmployeeNickname string `json:"employee_nickname"`
	ProfitCenterId   int    `json:"profit_center_id"`
}

func GetUserDetailsByID(id int) (UserDetailsResponse, *exceptions.BaseErrorResponse) {
	var userDetails UserDetailsResponse
	url := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &userDetails)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve user details due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "user details service is temporarily unavailable"
		}

		return userDetails, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting user details by ID"),
		}
	}
	return userDetails, nil
}
