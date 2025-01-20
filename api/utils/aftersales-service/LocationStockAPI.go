package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// GetAvailableItemLocationStock amlocationstock select option 1
func GetAvailableItemLocationStock(payload masterwarehousepayloads.GetAvailableQuantityPayload) (masterwarehousepayloads.GetQuantityAvailablePayload, *exceptions.BaseErrorResponse) {
	formattedTime := payload.PeriodDate.UTC().Format("2006-01-02T15:04:05.000Z")
	url := fmt.Sprintf("%slocation-stock/available_quantity?company_id=%s&warehouse_id=%s&location_id=%s&warehouse_group_id=%s&item_id=%s&uom_id=%s&period_date=%s",
		//config.EnvConfigs.AfterSalesServiceUrl,
		config.EnvConfigs.AfterSalesServiceUrl,
		strconv.Itoa(payload.CompanyId),
		strconv.Itoa(payload.WarehouseId),
		strconv.Itoa(payload.LocationId),
		strconv.Itoa(payload.WarehouseGroupId),
		strconv.Itoa(payload.ItemId),
		strconv.Itoa(payload.UomTypeId),
		formattedTime,
	)
	result := masterwarehousepayloads.GetQuantityAvailablePayload{}
	err := utils.CallAPI("GET", url, nil, &result)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve Location Stock due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Location Stock service is temporarily unavailable"
		}

		return result, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting Location Stock source conversion"),
		}
	}
	return result, nil
	//http://10.1.7.9:8000/v1/location-stock/available_quantity?company_id=473&warehouse_id=7&location_id=1&warehouse_group_id=1&item_id=293773&uom_id=1&period_date=2024-06-28T00%3A00%3A00.000Z
}
