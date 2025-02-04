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
	UserIdNotIn    string `json:"user_id_not_in"`
	EmployeeName   string `json:"employee_name"`
	RoleName       string `json:"role_name"`
	SortBy         string `json:"sort_by"`
	SortOf         string `json:"sort_of"`
}
type UserDetailsResponse struct {
	IsActive       bool   `json:"is_active"`
	UserEmployeeId int    `json:"user_detail_id"`
	UserId         int    `json:"user_id"`
	EmployeeName   string `json:"employee_name"`
	Username       string `json:"username"`
	CompanyId      int    `json:"company_id"`
	CompanyName    string `json:"company_name"`
	RoleId         int    `json:"role_id"`
	RoleName       string `json:"role_name"`
}

type UserCompanyAccessResponse struct {
	CompanyId int    `json:"company_id"`
	UserId    int    `json:"user_id"`
	RoleId    int    `json:"role_id"`
	IsActive  bool   `json:"is_active"`
	Username  string `json:"username"`
}

type EmployeeMasterResponse struct {
	UserEmployeeId int    `json:"user_employee_id"`
	EmployeeName   string `json:"employee_name"`
	CostCenterId   int    `json:"cost_center_id"`
}
type Address struct {
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageId      int    `json:"village_id"`
}
type CompanyAccessData struct {
	IsActive                 bool   `json:"is_active"`
	DealerRepresentativeId   string `json:"dealer_representative_id"`
	DealerRepresentativeCode string `json:"dealer_representative_code"`
	DealerRepresentativeName bool   `json:"dealer_representative_name"`
	WarehouseId              int    `json:"warehouse_id"`
}
type CompanyAccess struct {
	Page      int `json:"page"`
	PageLimit int `json:"page_limit"`
	Npages    int `json:"npages"`
	Nrows     int `json:"nrows"`
	Data      []CompanyAccessData
}
type CompanyOutlet struct {
	Page      int           `json:"page"`
	PageLimit int           `json:"page_limit"`
	Npages    int           `json:"npages"`
	Nrows     int           `json:"nrows"`
	Data      []interface{} `json:"data"`
}
type BackAccount struct {
	Page      int           `json:"page"`
	PageLimit int           `json:"page_limit"`
	Npages    int           `json:"npages"`
	Nrows     int           `json:"nrows"`
	Data      []interface{} `json:"data"`
}
type EmployeeMasterResponses struct {
	IsActive          bool          `json:"is_active"`
	UserEmployeeId    int           `json:"user_employee_id"`
	UserId            int           `json:"user_id"`
	EmployeeName      string        `json:"employee_name"`
	EmployeeNickname  string        `json:"employee_nickname"`
	IdTypeId          int           `json:"id_type_id"`
	IdNumber          string        `json:"id_number"`
	CompanyId         int           `json:"company_id"`
	JobTitleId        int           `json:"job_title_id"`
	JobPositionId     int           `json:"job_position_id"`
	DivisionId        int           `json:"division_id"`
	CostCenterId      int           `json:"cost_center_id"`
	ProfitCenterId    int           `json:"profit_center_id"`
	AddressId         int           `json:"address_id"`
	Address           Address       `json:"address"`
	OfficePhoneNumber interface{}   `json:"office_phone_number"`
	HomePhoneNumber   string        `json:"home_phone_number"`
	MobilePhone       string        `json:"mobile_phone"`
	EmailAddress      string        `json:"email_address"`
	StartDate         string        `json:"start_date"`
	TerminationDate   string        `json:"termination_date"`
	GenderId          int           `json:"gender_id"`
	DateOfBirth       string        `json:"date_of_birth"`
	CityOfBirthId     int           `json:"city_of_birth_id"`
	MaritalStatusId   int           `json:"marital_status_id"`
	NumberOfChildren  int           `json:"number_of_children"`
	CitizenshipId     int           `json:"citizenship_id"`
	LastEducationId   int           `json:"last_education_id"`
	LastEmployment    string        `json:"last_employment"`
	FactorX           float64       `json:"factor_x"`
	SkillLevelId      int           `json:"skill_level_id"`
	CompanyAccess     CompanyAccess `json:"company_access"`
	CompanyOutlet     CompanyOutlet `json:"company_outlet"`
	BankAccount       BackAccount   `json:"bank_account"`
}
type EmployeeCompanyAccessParams struct {
	Page                int    `json:"page"`
	Limit               int    `json:"limit"`
	UserCompanyAccessId int    `json:"id"`
	UserId              int    `json:"user_id"`
	IsActive            bool   `json:"is_active"`
	CompanyId           int    `json:"company_id"`
	SortBy              string `json:"sort_by"`
	SortOf              string `json:"sort_of"`
}

type EmployeeCompanyAccessResponse struct {
	IsActive  bool `json:"is_active"`
	Id        int  `json:"id"`
	CompanyId int  `json:"company_id"`
	RoleId    int  `json:"role_id"`
}

func GetEmployeeById(id int) (EmployeeMasterResponse, *exceptions.BaseErrorResponse) {
	var getEmployee EmployeeMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "user-detail?user_id=" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getEmployee)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve employee due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "employee service is temporarily unavailable"
		}

		return getEmployee, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting employee by ID"),
		}
	}
	return getEmployee, nil
}

func GetEmployeeMasterById(id int) (EmployeeMasterResponses, *exceptions.BaseErrorResponse) {
	var getEmployee EmployeeMasterResponses
	url := config.EnvConfigs.GeneralServiceUrl + "user-detail?user_id=" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getEmployee)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve employee due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "employee service is temporarily unavailable"
		}

		return getEmployee, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting employee by ID"),
		}
	}
	return getEmployee, nil
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

func GetAllUserCompanyAccess(params EmployeeCompanyAccessParams) ([]EmployeeCompanyAccessResponse, *exceptions.BaseErrorResponse) {
	var getEmployeeMaster []EmployeeCompanyAccessResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "user-company-access-list"

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

	err := utils.CallAPI("GET", url, nil, &getEmployeeMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve user-detail master due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "user-detail master service is temporarily unavailable"
		}

		return getEmployeeMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting user-detail master by ID"),
		}
	}

	return getEmployeeMaster, nil
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
