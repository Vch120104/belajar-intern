package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
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
