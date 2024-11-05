package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type EmployeeMasterResponse struct {
	UserEmployeeId int    `json:"user_employee_id"`
	EmployeeName   string `json:"employee_name"`
	CostCenterId   int    `json:"cost_center_id"`
}

func GetEmployeeByID(id int) (EmployeeMasterResponse, *exceptions.BaseErrorResponse) {
	var getEmployee EmployeeMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getEmployee)
	if err != nil {
		return getEmployee, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching employee data by ID",
			Err:        errors.New("error consuming external API for employee data by ID"),
		}
	}
	return getEmployee, nil
}
