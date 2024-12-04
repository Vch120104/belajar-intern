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

type CustomerMasterVatResponse struct {
	NpwpNumber         string `json:"npwp_no"`
	NpwpDate           string `json:"npwp_date"`
	PkpType            bool   `json:"pkp_type"`
	PkpNumber          string `json:"pkp_no"`
	PkpDate            string `json:"pkp_date"`
	TaxTransactionId   int    `json:"tax_transaction_id"`
	Name               string `json:"name"`
	AddressStreet1     string `json:"address_street_1"`
	AddressStreet2     string `json:"address_street_2"`
	AddressStreet3     string `json:"address_street_3"`
	VillageId          int    `json:"village_id"`
	TaxServiceOfficeId int    `json:"tax_service_office_id"`
}

type CustomerMasterDetailResponse struct {
	IsActive                  bool                      `json:"is_active"`
	CustomerId                int                       `json:"customer_id"`
	CustomerCode              string                    `json:"customer_code"`
	CustomerName              string                    `json:"customer_name"`
	CustomerTitlePrefix       string                    `json:"customer_title_prefix"`
	CustomerTitleSuffix       string                    `json:"customer_title_suffix"`
	ClientTypeId              int                       `json:"client_type_id"`
	IdType                    int                       `json:"id_type"`
	IdNumber                  string                    `json:"id_number"`
	AddressId                 int                       `json:"id_address_id"`
	CustomerMasterVatResponse CustomerMasterVatResponse `json:"vat_customer"`
	TaxCustomer               TaxCustomer               `json:"tax_customer"`
}

type TaxCustomer struct {
	NpwpNumber string `json:"npwp_no"`
}

type CustomerMasterResponse struct {
	CustomerId     int    `json:"customer_id"`
	CustomerCode   string `json:"customer_code"`
	CustomerName   string `json:"customer_name"`
	IdType         int    `json:"id_type"`
	IdNumber       string `json:"id_number"`
	IdAddressId    int    `json:"id_address_id"`
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageName    string `json:"village_name"`
	VillageZipCode string `json:"village_zip_code"`
	DistrictName   string `json:"district_name"`
	CityName       string `json:"city_name"`
	CityPhoneArea  string `json:"city_phone_area"`
	ProvinceName   string `json:"province_name"`
	CountryName    string `json:"country_name"`
}

type CustomerMasterByCodeResponse struct {
	CustomerId   int    `json:"customer_id"`
	CustomerCode string `json:"customer_code"`
	CustomerName string `json:"customer_name"`
	ClientTypeId int    `json:"client_type_id"`
}

type CustomerMasterParams struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	CustomerId     int    `json:"customer_id"`
	CustomerCode   string `json:"customer_code"`
	CustomerName   string `json:"customer_name"`
	IdType         int    `json:"id_type"`
	IdNumber       string `json:"id_number"`
	IdAddressId    int    `json:"id_address_id"`
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageName    string `json:"village_name"`
	VillageZipCode string `json:"village_zip_code"`
	DistrictName   string `json:"district_name"`
	CityName       string `json:"city_name"`
	CityPhoneArea  string `json:"city_phone_area"`
	ProvinceName   string `json:"province_name"`
	CountryName    string `json:"country_name"`
	SortBy         string `json:"sort_by"`
	SortOf         string `json:"sort_of"`
}

func GetAllCustomerMaster(params CustomerMasterParams) ([]CustomerMasterResponse, *exceptions.BaseErrorResponse) {
	var getCustomerMaster []CustomerMasterResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "customer-list"

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

	err := utils.CallAPI("GET", url, nil, &getCustomerMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve customer due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "customer service is temporarily unavailable"
		}

		return getCustomerMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting customer by ID"),
		}
	}

	return getCustomerMaster, nil
}

func GetCustomerMasterDetailById(id int) (CustomerMasterDetailResponse, *exceptions.BaseErrorResponse) {
	var getCustomerMaster CustomerMasterDetailResponse
	url := config.EnvConfigs.GeneralServiceUrl + "customer-detail/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getCustomerMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve customer due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "customer service is temporarily unavailable"
		}

		return getCustomerMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting customer by ID"),
		}
	}
	return getCustomerMaster, nil
}

func GetCustomerMasterByID(id int) (CustomerMasterResponse, *exceptions.BaseErrorResponse) {
	var getCustomerMaster CustomerMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "customer/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getCustomerMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve customer due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "customer service is temporarily unavailable"
		}

		return getCustomerMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting customer by ID"),
		}
	}
	return getCustomerMaster, nil
}

func GetCustomerMasterByCode(code string) (CustomerMasterByCodeResponse, *exceptions.BaseErrorResponse) {
	var getCustomerMaster CustomerMasterByCodeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "customer-code/" + code
	err := utils.CallAPI("GET", url, nil, &getCustomerMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve customer due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "customer service is temporarily unavailable"
		}

		return getCustomerMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting customer by ID"),
		}
	}
	return getCustomerMaster, nil
}

func GetCustomerMultiId(ids []int, response interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)
	validIds := make([]string, 0, len(ids))

	for _, id := range ids {
		if id != 0 {
			validIds = append(validIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(validIds, ",") + "]"
	url := config.EnvConfigs.GeneralServiceUrl + "customer-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve customer due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "customer service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting customer by ID"),
		}
	}
	return nil
}
