package salesserviceapiutils

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

type UnitBrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

type UnitBrandParams struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
	SortBy    string `json:"sort_by"`
	SortOf    string `json:"sort_of"`
}

func GetAllUnitBrand(params UnitBrandParams) ([]UnitBrandResponse, *exceptions.BaseErrorResponse) {
	var getUnitBrand []UnitBrandResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "unit-brand-list"

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

	err := utils.CallAPI("GET", url, nil, &getUnitBrand)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit brand due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit brand service is temporarily unavailable"
		}

		return getUnitBrand, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit brand by ID"),
		}
	}

	return getUnitBrand, nil
}

func GetUnitBrandByCode(code string) (UnitBrandResponse, *exceptions.BaseErrorResponse) {
	var unitBrand UnitBrandResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-brand-by-code/" + code
	err := utils.CallAPI("GET", url, nil, &unitBrand)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve brand due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand service is temporarily unavailable"
		}

		return unitBrand, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting brand by code"),
		}
	}
	return unitBrand, nil
}

func GetUnitBrandById(id int) (UnitBrandResponse, *exceptions.BaseErrorResponse) {
	var unitBrand UnitBrandResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &unitBrand)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve brand due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand service is temporarily unavailable"
		}

		return unitBrand, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting brand by ID"),
		}
	}
	return unitBrand, nil
}

func GetUnitBrandByMultiId(ids []int) ([]UnitBrandResponse, *exceptions.BaseErrorResponse) {
	var unitBrand []UnitBrandResponse

	ids = utils.RemoveDuplicateIds(ids)
	if len(ids) == 0 {
		return unitBrand, nil
	}

	var strIds string
	for _, id := range ids {
		if id != 0 {
			strIds += strconv.Itoa(id) + ","
		}
	}
	if strIds != "" {
		strIds = "[" + strIds[:len(strIds)-1] + "]"
	} else {
		strIds = "[]"
	}

	url := config.EnvConfigs.SalesServiceUrl + "unit-brand-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, &unitBrand)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve brand due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand service is temporarily unavailable"
		}

		return unitBrand, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting brand by multi ID"),
		}
	}
	return unitBrand, nil
}
