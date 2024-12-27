package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type LineType0Response struct {
	Description      string  `json:"description"`
	FRT              float64 `json:"frt"`
	ModelCode        string  `json:"model_code"`
	PackageCode      string  `json:"package_code"`
	PackageID        int     `json:"package_id"`
	PackageName      string  `json:"package_name"`
	Price            int     `json:"price"`
	ProfitCenter     int     `json:"profit_center"`
	ProfitCenterName string  `json:"profit_center_name"`
}

type LineType1Response struct {
	FrtHour                     int     `json:"frt_hour"`
	OperationCode               string  `json:"operation_code"`
	OperationEntriesCode        *string `json:"operation_entries_code"`        // Nullable
	OperationEntriesDescription *string `json:"operation_entries_description"` // Nullable
	OperationID                 int     `json:"operation_id"`
	OperationKeyCode            *string `json:"operation_key_code"`        // Nullable
	OperationKeyDescription     *string `json:"operation_key_description"` // Nullable
	OperationName               string  `json:"operation_name"`
}

type LineType2To9Response struct {
	AvailableQty   int    `json:"available_qty"`
	ItemCode       string `json:"item_code"`
	ItemID         int    `json:"item_id"`
	ItemLevel1     int    `json:"item_level_1"`
	ItemLevel1Code string `json:"item_level_1_code"`
	ItemLevel2     int    `json:"item_level_2"`
	ItemLevel2Code string `json:"item_level_2_code"`
	ItemLevel3     int    `json:"item_level_3"`
	ItemLevel3Code string `json:"item_level_3_code"`
	ItemLevel4     int    `json:"item_level_4"`
	ItemLevel4Code string `json:"item_level_4_code"`
	ItemName       string `json:"item_name"`
}

type ApiResponse struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
}

// GetOperationItemById fetches the operation item details based on LineTypeId and OperationItemId.
func GetOperationItemById(LineTypeStr string, OperationItemId int) (interface{}, *exceptions.BaseErrorResponse) {
	lineType, err := strconv.Atoi(LineTypeStr)
	if err != nil || lineType < 0 || lineType > 9 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid LineTypeId",
			Err:        errors.New("invalid LineTypeId value"),
		}
	}

	url := fmt.Sprintf("%slookup/item-opr-code/%s/by-id/%d", config.EnvConfigs.AfterSalesServiceUrl, LineTypeStr, OperationItemId)
	fmt.Println("Requesting URL:", url)

	var body []byte
	err = utils.CallAPI("GET", url, nil, &body)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to retrieve operation item due to an external service error",
			Err:        errors.New("error consuming external API"),
		}
	}

	fmt.Println("Raw Response:", string(body)) // Debugging

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error decoding response: %v", err),
			Err:        err,
		}
	}

	if apiResponse.StatusCode != http.StatusOK || apiResponse.Data == nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: apiResponse.StatusCode,
			Message:    apiResponse.Message,
			Err:        errors.New(apiResponse.Message),
		}
	}

	switch LineTypeStr {
	case "0":
		var response LineType0Response
		if err := json.Unmarshal(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response for LineType 0: %v", err),
				Err:        err,
			}
		}
		return response, nil

	case "1":
		var response LineType1Response
		if err := json.Unmarshal(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response for LineType 1: %v", err),
				Err:        err,
			}
		}
		return response, nil

	case "2", "3", "4", "5", "6", "7", "8", "9":
		var response LineType2To9Response
		if err := json.Unmarshal(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response for LineType 2-9: %v", err),
				Err:        err,
			}
		}
		return response, nil

	default:
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid LineTypeId provided",
			Err:        errors.New("invalid LineTypeId provided"),
		}
	}
}
