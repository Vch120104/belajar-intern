package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
	"strconv"
)

type DivisionResponseAPI struct {
	DivisionId   int    `json:"division_id"`
	DivisionCode string `json:"division_code"`
	DivisionName string `json:"division_name"`
}

func GetDivisionById(divisionId int) (DivisionResponseAPI, *exceptions.BaseErrorResponse) {
	var DivisionResponse DivisionResponseAPI

	DivisionURL := config.EnvConfigs.GeneralServiceUrl + "division/" + strconv.Itoa(divisionId)
	if err := utils.CallAPI("GET", DivisionURL, nil, &DivisionResponse); err != nil {
		return DivisionResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Division data from external service",
			Err:        err,
		}
	}
	return DivisionResponse, nil
}
