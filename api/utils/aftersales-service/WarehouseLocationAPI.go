package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type GetWarehouseLocationResponseAPI struct {
	WarehouseLocationId   int    `json:"warehouse_location_id"`
	WarehouseLocationCode string `json:"warehouse_location_code"`
	WarehouseLocationName string `json:"warehouse_location_name"`
}

func GetWarehouseLocationById(id int) (GetWarehouseLocationResponseAPI, *exceptions.BaseErrorResponse) {
	var warehouseLocation GetWarehouseLocationResponseAPI
	url := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-location/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &warehouseLocation)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve item due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "item service is temporarily unavailable"
		}

		return warehouseLocation, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting item by ID"),
		}
	}
	return warehouseLocation, nil
}
