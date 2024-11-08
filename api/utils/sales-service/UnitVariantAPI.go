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
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consume external variant api",
			Err:        errors.New("error consume external variant api"),
		}
	}
	return response, nil
}
