package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
	"strconv"
)

type LineTypeResponse struct {
	LineTypeId   int    `json:"line_type_id"`
	LineTypeCode string `json:"line_type_code"`
	LineTypeName string `json:"line_type_name"`
}

func GetLineTypeById(id int) (LineTypeResponse, *exceptions.BaseErrorResponse) {
	var line LineTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &line)
	if err != nil {
		return line, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get line type by id",
			Err:        err,
		}
	}
	return line, nil
}
