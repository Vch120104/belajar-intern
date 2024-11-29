package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
	"strconv"
)

type ItemGroupResponse struct {
	ItemGroupId   int    `json:"item_group_id"`
	ItemGroupCode string `json:"item_group_code"`
	ItemGroupName string `json:"item_group_name"`
}

func GetItemGroupById(itemGroupId int) (ItemGroupResponse, *exceptions.BaseErrorResponse) {
	var ItemGroup ItemGroupResponse
	ItemGroupURL := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(itemGroupId)
	if err := utils.CallAPI("GET", ItemGroupURL, nil, &ItemGroup); err != nil {
		return ItemGroup, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Group data from external service",
			Err:        err,
		}
	}
	return ItemGroup, nil
}
