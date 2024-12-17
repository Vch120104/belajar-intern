package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
	"strconv"
)

type ItemClassResponse struct {
	ItemClassId   int    `json:"item_class_id"`
	ItemClassCode string `json:"item_class_code"`
	ItemClassName string `json:"item_class_name"`
}

func GetItemClassById(Id int) (ItemClassResponse, *exceptions.BaseErrorResponse) {
	var ItemClass ItemClassResponse
	ItemClassURL := config.EnvConfigs.AfterSalesServiceUrl + "item-class/" + strconv.Itoa(Id)
	if err := utils.CallAPI("GET", ItemClassURL, nil, &ItemClass); err != nil {
		return ItemClass, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Class data from external service",
			Err:        err,
		}
	}
	return ItemClass, nil
}

func GetItemClassByCode(code string) (ItemClassResponse, *exceptions.BaseErrorResponse) {
	var ItemClass ItemClassResponse
	ItemClassURL := config.EnvConfigs.AfterSalesServiceUrl + "item-class/by-code/" + code
	if err := utils.CallAPI("GET", ItemClassURL, nil, &ItemClass); err != nil {
		return ItemClass, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Class data from external service",
			Err:        err,
		}
	}
	return ItemClass, nil
}
