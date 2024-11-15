package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type BrandDocResponse struct {
	BrandId           int    `json:"brand_id"`
	BrandCode         string `json:"brand_code"`
	BrandName         string `json:"brand_name"`
	BrandAbbreviation string `json:"brand_abbreveation"`
}

func GetBrandGenerateDoc(id int) (BrandDocResponse, *exceptions.BaseErrorResponse) {
	var brandDoc BrandDocResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, &brandDoc, nil)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve brand generate doc due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand generate doc service is temporarily unavailable"
		}

		return brandDoc, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting brand generate doc by ID"),
		}
	}
	return brandDoc, nil
}
