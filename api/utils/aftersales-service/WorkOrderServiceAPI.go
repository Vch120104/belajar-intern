package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WorkOrderServiceResponse struct {
	WorkOrderServiceId     int    `json:"work_order_service_id"`
	WorkOrderSystemNumber  int    `json:"work_order_system_number"`
	WorkOrderServiceRemark string `json:"work_order_service_remark"`
}

func GetWorkOrderServiceById(wosys int, id int) (WorkOrderServiceResponse, *exceptions.BaseErrorResponse) {
	var workOrderServiceResponse WorkOrderServiceResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "/normal/" + strconv.Itoa(wosys) + "/requestservice/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &workOrderServiceResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "work order service is temporarily unavailable"
		}

		return workOrderServiceResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting work order by ID"),
		}
	}
	return workOrderServiceResponse, nil
}
