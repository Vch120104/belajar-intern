package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WorkOrderTypeResponse struct {
	WorkOrderTypeId   int    `json:"work_order_type_id"`
	WorkOrderTypeCode string `json:"work_order_type_code"`
	WorkOrderTypeName string `json:"work_order_type_name"`
}

func GetWorkOrderTypeByID(id int) (WorkOrderTypeResponse, *exceptions.BaseErrorResponse) {
	var workOrderType WorkOrderTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &workOrderType)

	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Work order type service is temporarily unavailable"
		}

		return workOrderType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting work order type by ID"),
		}
	}

	return workOrderType, nil
}
