package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type ProfitCenterResponse struct {
	ProfitCenterId   int    `json:"profit_center_id"`
	ProfitCenterCode string `json:"profit_center_code"`
	ProfitCenterName string `json:"profit_center_name"`
}

func GetProfitCenterByCode(code string) (ProfitCenterResponse, *exceptions.BaseErrorResponse) {
	var getProfitCenter ProfitCenterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "profit-center-by-code/" + code

	err := utils.CallAPI("GET", url, nil, &getProfitCenter)
	if err != nil {
		return getProfitCenter, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching profit center by code",
			Err:        errors.New("failed to retrieve profit center data from external API by code"),
		}
	}
	return getProfitCenter, nil
}

func GetProfitCenterById(id int) (ProfitCenterResponse, *exceptions.BaseErrorResponse) {
	var getProfitCenter ProfitCenterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "profit-center/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getProfitCenter)
	if err != nil {
		return getProfitCenter, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching profit center by ID",
			Err:        errors.New("failed to retrieve profit center data from external API by ID"),
		}
	}
	return getProfitCenter, nil
}

func GetProfitCenterByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)
	var nonZeroIds []string
	for _, id := range ids {
		if id != 0 {
			nonZeroIds = append(nonZeroIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(nonZeroIds, ",") + "]"
	url := config.EnvConfigs.GeneralServiceUrl + "profit-center-by-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching profit center by multiple IDs",
			Err:        errors.New("failed to retrieve profit center data from external API for multiple IDs"),
		}
	}
	return nil
}
