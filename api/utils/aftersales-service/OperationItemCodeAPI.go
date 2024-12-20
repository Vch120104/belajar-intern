package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	OperationEntriesCode        *string `json:"operation_entries_code"`
	OperationEntriesDescription *string `json:"operation_entries_description"`
	OperationID                 int     `json:"operation_id"`
	OperationKeyCode            *string `json:"operation_key_code"`
	OperationKeyDescription     *string `json:"operation_key_description"`
	OperationName               string  `json:"operation_name"`
}

type LineType2To9Response struct {
	AvailableQty   int     `json:"available_qty"`
	ItemCode       string  `json:"item_code"`
	ItemID         int     `json:"item_id"`
	ItemLevel1     int     `json:"item_level_1"`
	ItemLevel1Code string  `json:"item_level_1_code"`
	ItemLevel2     *int    `json:"item_level_2"`
	ItemLevel2Code *string `json:"item_level_2_code"`
	ItemLevel3     *int    `json:"item_level_3"`
	ItemLevel3Code *string `json:"item_level_3_code"`
	ItemLevel4     *int    `json:"item_level_4"`
	ItemLevel4Code *string `json:"item_level_4_code"`
	ItemName       string  `json:"item_name"`
}

type ApiResponse struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
}

// GetOperationItemById fetches the operation item details based on LineTypeId and OperationItemId.
func GetOperationItemById(LineTypeId int, OperationItemId int) (interface{}, *exceptions.BaseErrorResponse) {
	// Validate LineTypeId
	if LineTypeId <= 0 || LineTypeId > 9 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid LineTypeId",
			Err:        errors.New("invalid LineTypeId value"),
		}
	}

	// Construct the URL for the request
	url := fmt.Sprintf("%slookup/item-opr-code/%d/by-id/%d", config.EnvConfigs.AfterSalesServiceUrl, LineTypeId, OperationItemId)
	fmt.Println("Requesting URL:", url)

	// Make the API request using CallAPI utility
	var body []byte
	if err := utils.CallAPI("GET", url, nil, &body); err != nil {
		// Handle errors in the API request
		status := http.StatusBadGateway
		message := "Failed to retrieve operation item due to an external service error"
		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Operation item service is temporarily unavailable"
		}

		return nil, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting operation item by ID"),
		}
	}

	// Convert the response body to ApiResponse
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error decoding response: %v", err),
			Err:        err,
		}
	}

	body, _ = io.ReadAll(bytes.NewBuffer(apiResponse.Data))
	fmt.Println("Response Body:", string(body))
	// Check if the response is successful
	if apiResponse.StatusCode != http.StatusOK || apiResponse.Data == nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: apiResponse.StatusCode,
			Message:    fmt.Sprintf("Failed to retrieve operation item: %s", apiResponse.Message),
			Err:        errors.New(apiResponse.Message),
		}
	}

	// Select the appropriate struct based on LineTypeId
	switch LineTypeId {
	case 1:
		var response LineType0Response
		if err := json.Unmarshal(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response for lineType 0: %v", err),
				Err:        err,
			}
		}
		return response, nil

	case 2:
		var response LineType1Response
		if err := json.Unmarshal(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response for lineType 1: %v", err),
				Err:        err,
			}
		}
		return response, nil

	case 3, 4, 5, 6, 7, 8, 9:
		var response LineType2To9Response
		if err := json.Unmarshal(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response for lineType 2-9: %v", err),
				Err:        err,
			}
		}
		return response, nil

	default:
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid lineTypeId provided",
			Err:        errors.New("invalid lineTypeId provided"),
		}
	}
}

// ValidateOperationItemId validates if the given OperationItemId matches the expected value for the given LineTypeId.
func ValidateOperationItemId(lineTypeId, operationItemId int) (*http.Response, *exceptions.BaseErrorResponse) {
	url := fmt.Sprintf("%slookup/item-opr-code/%d/by-id/%d", config.EnvConfigs.AfterSalesServiceUrl, lineTypeId, operationItemId)
	fmt.Println("Requesting URL:", url)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error calling external service: %v", err),
			Err:        err,
		}
	}
	defer resp.Body.Close()

	// Log response body for debugging
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if resp.StatusCode != http.StatusOK {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Invalid combination of LineType & OperationItemId from external service",
			Err:        errors.New("invalid combination of LineType & OperationItemId from external service"),
		}
	}

	// Validate based on LineTypeId
	switch lineTypeId {
	case 1:
		var response LineType0Response
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response: %v", err),
				Err:        err,
			}
		}
		if response.PackageID != operationItemId {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "OperationItemId is invalid for LineType 0",
				Err:        errors.New("OperationItemId is invalid for LineType 0"),
			}
		}

	case 2:
		var response LineType1Response
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response: %v", err),
				Err:        err,
			}
		}
		if response.OperationID != operationItemId {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "OperationItemId is invalid for LineType 1",
				Err:        errors.New("OperationItemId is invalid for LineType 1"),
			}
		}

	case 3, 4, 5, 6, 7, 8, 9:
		var response LineType2To9Response
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Error decoding response: %v", err),
				Err:        err,
			}
		}
		if response.ItemID != operationItemId {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "OperationItemId is invalid for LineType 2-9",
				Err:        errors.New("OperationItemId is invalid for LineType 2-9"),
			}
		}

	default:
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid LineType provided",
			Err:        errors.New("invalid LineType provided"),
		}
	}

	return resp, nil
}
