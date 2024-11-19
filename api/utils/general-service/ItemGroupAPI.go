package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type GetAllItemGroupParams struct {
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
	IsActive      string `json:"is_active"`
	ItemGroupCode string `json:"item_group_code"`
	ItemGroupName string `json:"item_group_name"`
	SortBy        string `json:"sort_by"`
	SortOf        string `json:"sort_of"`
}

type ItemGroupResponse struct {
	IsActive      bool   `json:"is_active"`
	ItemGroupId   int    `json:"item_group_id"`
	ItemGroupCode string `json:"item_group_code"`
	ItemGroupName string `json:"item_group_name"`
}

type GetAllItemGroupResponse struct {
	StatusCode int                 `json:"status_code"`
	Message    string              `json:"message"`
	Data       []ItemGroupResponse `json:"data"`
}

func GetAllItemGroup(params GetAllItemGroupParams) (GetAllItemGroupResponse, *exceptions.BaseErrorResponse) {
	var getItemGroup GetAllItemGroupResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "item-groups"

	queryParams := fmt.Sprintf("page=%d&limit=%d", params.Page, params.Limit)

	v := reflect.ValueOf(params)
	typeOfParams := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		if strVal, ok := value.(string); ok && strVal != "" {
			key := typeOfParams.Field(i).Tag.Get("json")
			value := strings.ReplaceAll(strVal, " ", "%20")
			queryParams += "&" + key + "=" + value
		}
	}

	url := baseURL + "?" + queryParams

	err := utils.GetArray(url, nil, &getItemGroup)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve item group due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "item group service is temporarily unavailable"
		}

		return getItemGroup, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting item group"),
		}
	}
	return getItemGroup, nil
}
