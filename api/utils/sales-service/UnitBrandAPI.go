package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type UnitBrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
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

func GetUnitBrandByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {
	ids = utils.RemoveDuplicateIds(ids)
	if len(ids) == 0 {
		return nil
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
	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve brand due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting brand by multi ID"),
		}
	}
	return nil
}
