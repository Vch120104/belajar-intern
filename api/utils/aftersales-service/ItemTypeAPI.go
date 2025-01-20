package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
	"strconv"
)

type ItemTypeResponse struct {
	ItemTypeId   int    `json:"item_type_id"`
	ItemTypeCode string `json:"item_type_code"`
	ItemTypeName string `json:"item_type_name"`
}

func GetItemTypeById(Id int) (ItemTypeResponse, *exceptions.BaseErrorResponse) {
	var ItemType ItemTypeResponse
	ItemTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "item-type/" + strconv.Itoa(Id)
	if err := utils.CallAPI("GET", ItemTypeURL, nil, &ItemType); err != nil {
		return ItemType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Type data from external service",
			Err:        err,
		}
	}
	return ItemType, nil
}

func GetItemTypeByCode(code string) (ItemTypeResponse, *exceptions.BaseErrorResponse) {
	var ItemType ItemTypeResponse
	ItemTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "item-type/code/" + code
	if err := utils.CallAPI("GET", ItemTypeURL, nil, &ItemType); err != nil {
		return ItemType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Type data from external service",
			Err:        err,
		}
	}
	return ItemType, nil
}
