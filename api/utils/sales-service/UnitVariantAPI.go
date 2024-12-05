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

type UnitVariantResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantName        string `json:"variant_name"`
	VariantDescription string `json:"variant_description"`
}

type UnitVariantMultiIdResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantDescription string `json:"variant_description"`
}

type UnitVariantByBrandResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantDescription string `json:"variant_description"`
	ModelId            int    `json:"model_id"`
	ModelCode          string `json:"model_code"`
	ModelDescription   string `json:"model_description"`
	BrandId            int    `json:"brand_id"`
	BrandCode          string `json:"brand_code"`
	BrandName          string `json:"brand_name"`
}

type UnitVariantParams struct {
	Page               int    `json:"page"`
	Limit              int    `json:"limit"`
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantDescription string `json:"variant_description"`
	SortBy             string `json:"sort_by"`
	SortOf             string `json:"sort_of"`
}

func GetAllUnitVariant(params UnitVariantParams) ([]UnitVariantResponse, *exceptions.BaseErrorResponse) {
	var getUnitVariant []UnitVariantResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "unit-variant-list"

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

	err := utils.CallAPI("GET", url, nil, &getUnitVariant)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit variant due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit variant service is temporarily unavailable"
		}

		return getUnitVariant, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit variant by ID"),
		}
	}

	return getUnitVariant, nil
}

func GetUnitVariantById(id int) (UnitVariantResponse, *exceptions.BaseErrorResponse) {
	var response UnitVariantResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit variant due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit variant service is temporarily unavailable"
		}

		return response, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit variant by ID"),
		}
	}
	return response, nil
}

func GetUnitVariantByMultiId(ids []int) ([]UnitVariantMultiIdResponse, *exceptions.BaseErrorResponse) {
	var getUnitVariant []UnitVariantMultiIdResponse
	ids = utils.RemoveDuplicateIds(ids)
	if len(ids) == 0 {
		return getUnitVariant, nil
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

	url := config.EnvConfigs.SalesServiceUrl + "unit-variant-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, &getUnitVariant)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit variant due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit variant service is temporarily unavailable"
		}

		return getUnitVariant, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit variant by multi ID"),
		}
	}

	return getUnitVariant, nil
}

func GetUnitVariantByBrand(brandId int) ([]UnitVariantByBrandResponse, *exceptions.BaseErrorResponse) {
	var getUnitVariant []UnitVariantByBrandResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-variant-by-brand/" + strconv.Itoa(brandId)
	err := utils.CallAPI("GET", url, nil, &getUnitVariant)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit variant due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit variant service is temporarily unavailable"
		}

		return getUnitVariant, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit variant by brand"),
		}
	}

	return getUnitVariant, nil
}
