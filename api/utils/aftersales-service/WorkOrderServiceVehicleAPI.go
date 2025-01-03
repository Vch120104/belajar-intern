package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type WorkOrderServiceVehicleResponse struct {
	WorkOrderServiceVehicleId int       `json:"work_order_service_vehicle_id"`
	WorkOrderSystemNumber     int       `json:"work_order_system_number"`
	WorkOrderVehicleDate      time.Time `json:"work_order_vehicle_date"`
	WorkOrderVehicleRemark    string    `json:"work_order_vehicle_remark"`
}

func GetWorkOrderServiceVehicleById(wosys int, id int) (WorkOrderServiceVehicleResponse, *exceptions.BaseErrorResponse) {
	var workOrderServiceVehicleResponse WorkOrderServiceVehicleResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "/normal/" + strconv.Itoa(wosys) + "/vehicleservice/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &workOrderServiceVehicleResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "work order service is temporarily unavailable"
		}

		return workOrderServiceVehicleResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting work order by ID"),
		}
	}
	return workOrderServiceVehicleResponse, nil
}
