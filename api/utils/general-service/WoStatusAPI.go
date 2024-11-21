package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WorkOrderStatusResponse struct {
	WorkOrderStatusId   int    `json:"work_order_status_id"`
	WorkOrderStatusCode string `json:"work_order_status_code"`
	WorkOrderStatusName string `json:"work_order_status_description"`
}

func GetWorkOrderStatusByID(id int) (WorkOrderStatusResponse, *exceptions.BaseErrorResponse) {
	var workOrderStatus WorkOrderStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-status/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &workOrderStatus)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Workorder status service is temporarily unavailable"
		}

		return workOrderStatus, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting workorder status by ID"),
		}
	}
	return workOrderStatus, nil
}
