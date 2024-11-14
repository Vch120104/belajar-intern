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
		return workOrderStatus, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order status",
			Err:        errors.New("error consuming external API while getting work order status by ID"),
		}
	}
	return workOrderStatus, nil
}
