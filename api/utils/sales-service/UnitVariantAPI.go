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
