package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type GetWarehouseGroupResponseAPI struct {
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
}

func GetWarehouseGroupById(id int) (GetWarehouseGroupResponseAPI, *exceptions.BaseErrorResponse) {
	var warehouseGroup GetWarehouseGroupResponseAPI
	url := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-group/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &warehouseGroup)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve warehouse group due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "warehouse group service is temporarily unavailable"
		}

		return warehouseGroup, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting warehouse group by ID"),
		}
	}
	return warehouseGroup, nil
}
