package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type GetWarehouseResponseAPI struct {
	WarehouseId   int    `json:"warehouse_id"`
	WarehouseCode string `json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
}

func GetWarehouseById(id int) (GetWarehouseResponseAPI, *exceptions.BaseErrorResponse) {
	var warehouse GetWarehouseResponseAPI
	url := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-master/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &warehouse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve item due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "item service is temporarily unavailable"
		}

		return warehouse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting item by ID"),
		}
	}
	return warehouse, nil
}
