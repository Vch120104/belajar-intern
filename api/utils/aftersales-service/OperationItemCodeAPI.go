package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Define response structures for different line types.
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
	OperationEntriesCode        *string `json:"operation_entries_code"`
	OperationEntriesDescription *string `json:"operation_entries_description"`
	OperationID                 int     `json:"operation_id"`
	OperationKeyCode            *string `json:"operation_key_code"`
	OperationKeyDescription     *string `json:"operation_key_description"`
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

// General structure for API responses.
type ApiResponse struct {
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}

// Custom API request handler for getting operation item by ID.
func GetOperationItemById(LineTypeStr string, OperationItemId int) (interface{}, *exceptions.BaseErrorResponse) {
	// Validate LineType
	lineType, err := strconv.Atoi(LineTypeStr)
	if err != nil || lineType < 0 || lineType > 9 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid LineType",
			Err:        fmt.Errorf("invalid LineType: %s", LineTypeStr),
		}
	}

	// URL for the request
	url := fmt.Sprintf("%slookup/item-opr-code/%s/by-id/%d", config.EnvConfigs.AfterSalesServiceUrl, LineTypeStr, OperationItemId)
	log.Printf("Requesting URL: %s", url)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create request",
			Err:        err,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to make request",
			Err:        err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    "Failed to get operation item",
			Err:        errors.New("failed to get operation item"),
		}
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to read response body",
			Err:        err,
		}
	}

	//log.Printf("Raw Response: %s", string(body))

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to decode API response",
			Err:        err,
		}
	}

	if apiResponse.StatusCode != http.StatusOK {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: apiResponse.StatusCode,
			Message:    apiResponse.Message,
			Err:        errors.New(apiResponse.Message),
		}
	}

	// Handle data according to LineType
	var responseData interface{}
	switch lineType {
	case 2, 3, 4, 5, 6, 7, 8, 9:
		var response LineType2To9Response
		if err := mapToStruct(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error unmarshaling 'data' into LineType2To9Response",
				Err:        err,
			}
		}
		responseData = response
	case 1:
		var response LineType1Response
		if err := mapToStruct(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error unmarshaling 'data' into LineType1Response",
				Err:        err,
			}
		}
		responseData = response
	case 0:
		var response LineType0Response
		if err := mapToStruct(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error unmarshaling 'data' into LineType0Response",
				Err:        err,
			}
		}
		responseData = response
	default:
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unknown line type in operation item response",
			Err:        fmt.Errorf("unexpected line type %d", lineType),
		}
	}

	return responseData, nil
}

// Helper function to map data into the correct struct
func mapToStruct(data map[string]interface{}, result interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	if err := json.Unmarshal(dataBytes, result); err != nil {
		return fmt.Errorf("error unmarshaling into struct: %v", err)
	}
	return nil
}

// aftersalesserviceapiutils/response_utils.go
func HandleLineTypeResponse(lineTypeCode string, operationItemResponse interface{}) (string, string, *exceptions.BaseErrorResponse) {
	var OperationItemCode, Description string

	switch lineTypeCode {
	case "0":
		if response, ok := operationItemResponse.(LineType0Response); ok {
			OperationItemCode = response.PackageCode
			Description = response.PackageName
		} else {
			return "", "", &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to decode LineType0 response",
				Err:        errors.New("failed to decode LineType0 response"),
			}
		}

	case "1":
		if response, ok := operationItemResponse.(LineType1Response); ok {
			OperationItemCode = response.OperationCode
			Description = response.OperationName
		} else {
			return "", "", &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to decode LineType1 response",
				Err:        errors.New("failed to decode LineType1 response"),
			}
		}

	default:
		if response, ok := operationItemResponse.(LineType2To9Response); ok {
			OperationItemCode = response.ItemCode
			Description = response.ItemName
		} else {
			return "", "", &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to decode LineType2-9 response",
				Err:        errors.New("failed to decode LineType2-9 response"),
			}
		}
	}

	return OperationItemCode, Description, nil
}
