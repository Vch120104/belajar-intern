package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type QuantityConversionUomResponse struct {
	SourceType         string  `json:"source_type"`
	ItemId             int     `json:"item_id"`
	Quantity           float64 `json:"quantity"`
	QuantityConversion float64 `json:"quantity_conversion"`
}

type UomResponse struct {
	UomId              int    `json:"uom_id"`
	UomCode            string `json:"uom_code"`
	UomDescription     string `json:"uom_description"`
	UomTypeId          int    `json:"uom_type_id"`
	UomTypeDescription string `json:"uom_type_description"`
}

//var UomBaseUrl string = config.EnvConfigs.AfterSalesServiceUrl + "unit-measurement/"

func GetQuantityConversion(SourceType string, itemId int, quantity float64) (QuantityConversionUomResponse, *exceptions.BaseErrorResponse) {
	Url := config.EnvConfigs.AfterSalesServiceUrl + "unit-measurement/" + fmt.Sprintf("get_quantity_conversion?source_type=%s&item_id=%s&quantity=%s", SourceType, strconv.Itoa(itemId), strconv.FormatFloat(quantity, 'f', 6, 64))
	response := QuantityConversionUomResponse{}
	err := utils.CallAPI("POST", Url, nil, &response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve UOM item due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "UOM service is temporarily unavailable"
		}

		return response, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting UOM source conversion"),
		}
	}
	return response, nil
}

func GetUomById(id int) (UomResponse, *exceptions.BaseErrorResponse) {
	var getUom UomResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getUom)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve uom due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "uom service is temporarily unavailable"
		}

		return getUom, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting uom by ID"),
		}
	}
	return getUom, nil
}
