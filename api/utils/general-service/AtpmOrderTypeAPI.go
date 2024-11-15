package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type AtpmOrderTypeResponse struct {
	IsActive                 bool   `json:"is_active"`
	AtpmOrderTypeId          int    `json:"atpm_order_type_id"`
	AtpmOrderTypeCode        string `json:"atpm_order_type_code"`
	AtpmOrderTypeDescription string `json:"atpm_order_type_description"`
}

func GetAtpmOrderTypeById(id int) (AtpmOrderTypeResponse, *exceptions.BaseErrorResponse) {
	var getAtpmOrderType AtpmOrderTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "atpm-order-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getAtpmOrderType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve atpm order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "atpm order type service is temporarily unavailable"
		}

		return getAtpmOrderType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting atpm order type by ID"),
		}
	}
	return getAtpmOrderType, nil
}
