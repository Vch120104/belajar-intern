package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type UnitModelResponse struct {
	ModelId              int    `json:"model_id"`
	ModelCode            string `json:"model_code"`
	ModelName            string `json:"model_description"`
	ModelCodeDescription string `json:"model_code_description"`
}

func GetUnitModelByCode(code string) (UnitModelResponse, *exceptions.BaseErrorResponse) {
	var getUnitModel UnitModelResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-model-by-code/" + code
	err := utils.CallAPI("GET", url, nil, &getUnitModel)
	if err != nil {
		return getUnitModel, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consume external model api",
			Err:        errors.New("error consume external model api"),
		}
	}
	return getUnitModel, nil
}

func GetUnitModelById(id int) (UnitModelResponse, *exceptions.BaseErrorResponse) {
	var getUnitModel UnitModelResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getUnitModel)
	if err != nil {
		return getUnitModel, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consume external model api",
			Err:        errors.New("error consume external model api"),
		}
	}
	return getUnitModel, nil
}
