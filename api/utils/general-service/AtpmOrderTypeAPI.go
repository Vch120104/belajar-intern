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
		return getAtpmOrderType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching atpm order type by id",
			Err:        errors.New("failed to retrieve atpm order type data from external API by id"),
		}
	}
	return getAtpmOrderType, nil
}
