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

type SupplierMasterVatResponse struct {
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

type SupplierMasterResponse struct {
	SupplierId                int                       `json:"supplier_id"`
	SupplierCode              string                    `json:"supplier_code"`
	SupplierName              string                    `json:"supplier_name"`
	SupplierTypeId            int                       `json:"supplier_type_id"`
	SupplierMasterVatResponse SupplierMasterVatResponse `json:"vat_supplier"`
	ClientTypeId              int                       `json:"client_type_id"`
}

type SupplierMasterParams struct {
	Page                      int    `json:"page"`
	Limit                     int    `json:"limit"`
	SupplierCode              string `json:"supplier_code"`
	SupplierName              string `json:"supplier_name"`
	ClientTypeCode            string `json:"client_type_code"`
	ClientTypeDescription     string `json:"client_type_description"`
	CompanyId                 string `json:"company_id"`
	Address_1                 string `json:"address_1"`
	Address_2                 string `json:"address_2"`
	Address_3                 string `json:"address_3"`
	IsActive                  string `json:"is_active"`
	SupplierStatusDescription string `json:"supplier_status_description"`
	SortBy                    string `json:"sort_by"`
	SortOf                    string `json:"sort_of"`
}

type SupplierMasterGetAllResponse struct {
	StatusCode int                      `json:"status_code"`
	Message    string                   `json:"message"`
	Page       int                      `json:"page"`
	PageLimit  int                      `json:"page_limit"`
	TotalRows  int                      `json:"total_rows"`
	TotalPages int                      `json:"total_pages"`
	Data       []SupplierMasterResponse `json:"data"`
}

func GetAllSupplierMaster(params SupplierMasterParams) (SupplierMasterGetAllResponse, *exceptions.BaseErrorResponse) {
	var getSupplierMaster SupplierMasterGetAllResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "supplier"

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

	err := utils.GetArray(url, nil, &getSupplierMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve supplier master due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "supplier master service is temporarily unavailable"
		}

		return getSupplierMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting supplier master by ID"),
		}
	}

	return getSupplierMaster, nil
}

func GetSupplierMasterByCode(code string) (SupplierMasterResponse, *exceptions.BaseErrorResponse) {
	var getSupplierMaster SupplierMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "supplier-by-code/" + code

	err := utils.CallAPI("GET", url, nil, &getSupplierMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve supplier master due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "supplier master service is temporarily unavailable"
		}

		return getSupplierMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting supplier master by ID"),
		}
	}
	return getSupplierMaster, nil
}

func GetSupplierMasterByID(id int) (SupplierMasterResponse, *exceptions.BaseErrorResponse) {
	var getSupplierMaster SupplierMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getSupplierMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve supplier master due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "supplier master service is temporarily unavailable"
		}

		return getSupplierMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting supplier master by ID"),
		}
	}
	return getSupplierMaster, nil
}

func GetSupplierMasterByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {
	ids = utils.RemoveDuplicateIds(ids)
	var nonZeroIds []string

	for _, id := range ids {
		if id != 0 {
			nonZeroIds = append(nonZeroIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(nonZeroIds, ",") + "]"
	url := config.EnvConfigs.GeneralServiceUrl + "supplier-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve supplier master due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "supplier master service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting supplier master by ID"),
		}
	}
	return nil
}
