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
		return unitBrand, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching unit brand by code",
			Err:        errors.New("error consuming external API while fetching unit brand by code"),
		}
	}
	return unitBrand, nil
}

func GetUnitBrandById(id int) (UnitBrandResponse, *exceptions.BaseErrorResponse) {
	var unitBrand UnitBrandResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &unitBrand)
	if err != nil {
		return unitBrand, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching unit brand by ID",
			Err:        errors.New("error consuming external API while fetching unit brand by ID"),
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
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching unit brands by multiple IDs",
			Err:        errors.New("error consuming external API while fetching unit brands by multiple IDs"),
		}
	}
	return nil
}
