package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type UserDetailParams struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	UserEmployeeId int    `json:"user_employee_id"`
	UserId         int    `json:"user_id"`
	EmployeeName   string `json:"employee_name"`
	SortBy         string `json:"sort_by"`
	SortOf         string `json:"sort_of"`
}
type UserDetailsResponse struct {
	IsActive         bool   `json:"is_active"`
	UserEmployeeId   int    `json:"user_employee_id"`
	UserId           int    `json:"user_id"`
	EmployeeName     string `json:"employee_name"`
	EmployeeNickname string `json:"employee_nickname"`
	ProfitCenterId   int    `json:"profit_center_id"`
}

type UserCompanyAccessResponse struct {
	CompanyId int `json:"company_id"`
	UserId    int `json:"user_id"`
	RoleId    int `json:"role_id"`
}

func GetAllUserDetail(params UserDetailParams) ([]UserDetailsResponse, *exceptions.BaseErrorResponse) {
	var getUserDetailMaster []UserDetailsResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "user-detail-list"

	queryParams := fmt.Sprintf("page=%d&limit=%d", params.Page, params.Limit)

	v := reflect.ValueOf(params)
	typeOfParams := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		if strVal, ok := value.(string); ok && strVal != "" {
			key := typeOfParams.Field(i).Tag.Get("json")
			value := strings.ReplaceAll(strVal, " ", "%20")
			queryParams += "&" + key + "=" + value
		}
	}

	url := baseURL + "?" + queryParams

	err := utils.CallAPI("GET", url, nil, &getUserDetailMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve user-detail master due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "user-detail master service is temporarily unavailable"
		}

		return getUserDetailMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting user-detail master by ID"),
		}
	}

	return getUserDetailMaster, nil
}

func GetUserCompanyAccessById(id int) (UserCompanyAccessResponse, *exceptions.BaseErrorResponse) {
	var companyAccess UserCompanyAccessResponse
	url := config.EnvConfigs.GeneralServiceUrl + "user-company-access/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &companyAccess)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve user company access due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "user company access service is temporarily unavailable"
		}

		return companyAccess, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting user company access by ID"),
		}
	}
	return companyAccess, nil
}

func GetUserDetailsByID(id int) (UserDetailsResponse, *exceptions.BaseErrorResponse) {
	var userDetails UserDetailsResponse
	url := config.EnvConfigs.GeneralServiceUrl + "user-detail?user_id=" + strconv.Itoa(id)

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
