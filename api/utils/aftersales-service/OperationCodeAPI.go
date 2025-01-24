package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type OperationCodeResponse struct {
	OperationId   int    `json:"operation_id"`
	OperationCode string `json:"operation_code"`
	OperationName string `json:"operation_name"`
}

func GetOperationById(id int) (OperationCodeResponse, *exceptions.BaseErrorResponse) {
	var getOperation OperationCodeResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "operation-model-mapping/by-id/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getOperation)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve operation due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "operation service is temporarily unavailable"
		}

		return getOperation, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting operation by ID"),
		}
	}
	return getOperation, nil
}
