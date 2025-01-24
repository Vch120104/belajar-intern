package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type ItemResponse struct {
	ItemId     int    `json:"item_id"`
	ItemCode   string `json:"item_code"`
	ItemName   string `json:"item_name"`
	UomStockId int    `json:"unit_of_measurement_stock_id"`
}

func GetItemId(id int) (ItemResponse, *exceptions.BaseErrorResponse) {
	var getItem ItemResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getItem)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve item due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "item service is temporarily unavailable"
		}

		return getItem, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting item by ID"),
		}
	}
	return getItem, nil
}
