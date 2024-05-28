package test

import (
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPurchasePriceService struct {
	mock.Mock
}

func (m *MockPurchasePriceService) GetAllPurchasePrice(criteria []utils.FilterCondition, paginate pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(criteria, paginate)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockPurchasePriceService) SavePurchasePrice(payload masteritempayloads.PurchasePriceRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(payload)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockPurchasePriceService) GetPurchasePriceById(id int) (masteritempayloads.PurchasePriceRequest, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Get(0).(masteritempayloads.PurchasePriceRequest), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockPurchasePriceService) GetAllPurchasePriceDetail(criteria []utils.FilterCondition, paginate pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(criteria, paginate)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockPurchasePriceService) AddPurchasePrice(payload masteritempayloads.PurchasePriceDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(payload)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockPurchasePriceService) DeletePurchasePrice(id int) *exceptions.BaseErrorResponse {
	args := m.Called(id)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockPurchasePriceService) ChangeStatusPurchasePrice(purchasePriceId int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	args := m.Called(purchasePriceId)
	return args.Get(0).(masteritementities.PurchasePrice), args.Get(1).(*exceptions.BaseErrorResponse)
}

func TestSavePurchasePrice_Success(t *testing.T) {
	payload := masteritempayloads.PurchasePriceRequest{
		// Adjust payload with appropriate data
		PurchasePriceId: 1,
		SupplierId:      1,
		CurrencyId:      1,
		IsActive:        true,
	}

	mockService := new(MockPurchasePriceService)
	mockService.On("SavePurchasePrice", payload).Return(true, (*exceptions.BaseErrorResponse)(nil))

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8000/v1/purchase-price", bytes.NewReader(payloadBytes))
	rr := httptest.NewRecorder()

	controller := masteritemcontroller.NewPurchasePriceController(mockService)
	controller.SavePurchasePrice(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)

	expectedStatusCode := http.StatusOK
	expectedMessage := "Update Data Successfully!"

	if payload.PurchasePriceId == 0 {
		expectedStatusCode = http.StatusCreated
		expectedMessage = "Create Data Successfully!"
	}

	assert.Equal(t, expectedStatusCode, statusCode, "Status code should match expected")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.Equal(t, expectedMessage, response["message"], "Message should match expected")
	assert.Equal(t, expectedStatusCode, int(response["status_code"].(float64)), "Status code should match expected")
}

func TestSavePurchasePrice_Failure(t *testing.T) {
	payload := masteritempayloads.PurchasePriceRequest{
		// Adjust payload with appropriate data
		PurchasePriceId: 1,
		SupplierId:      1,
		CurrencyId:      1,
		IsActive:        true,
	}

	mockService := new(MockPurchasePriceService)
	mockService.On("SavePurchasePrice", payload).
		Return(false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("some error"),
		})

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8000/v1/purchase-price", bytes.NewReader(payloadBytes))
	rr := httptest.NewRecorder()

	controller := masteritemcontroller.NewPurchasePriceController(mockService)
	controller.SavePurchasePrice(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusBadRequest, statusCode, "Status code should be 400")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.NotNil(t, response["error"], "Response should contain error")
	assert.Equal(t, "some error", response["error"], "Error should be 'some error'")
}

func TestGetAllPurchasePrice_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8000/v1/purchase-price", nil)
	rr := httptest.NewRecorder()

	responseData := []map[string]interface{}{
		{"key": "value"}, // Adjust with expected data
	}
	mockService := new(MockPurchasePriceService)
	mockService.On("GetAllPurchasePrice", mock.Anything, mock.Anything).
		Return(responseData, len(responseData), len(responseData), (*exceptions.BaseErrorResponse)(nil))

	controller := masteritemcontroller.NewPurchasePriceController(mockService)
	controller.GetAllPurchasePrice(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode, "Status code should be 200")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.NotNil(t, response["data"], "Response should contain data")

	responseDataFromResponse, ok := response["data"].([]interface{})
	assert.True(t, ok, "Data in response should be []interface{}")

	assert.Equal(t, len(responseData), len(responseDataFromResponse), "Length of response data should match")
}

func TestGetPurchasePriceById(t *testing.T) {
	id := 1 // Change to the appropriate ID
	req, _ := http.NewRequest("GET", "http://localhost:8000/v1/purchase-price/by-id/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()

	responseData := masteritempayloads.PurchasePriceDetailResponse{
		// Adjust with expected data
	}
	mockService := new(MockPurchasePriceService)
	mockService.On("GetPurchasePriceById", id).
		Return(responseData, (*exceptions.BaseErrorResponse)(nil))

	controller := masteritemcontroller.NewPurchasePriceController(mockService)
	controller.GetPurchasePriceById(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode, "Status code should be 200")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.NotNil(t, response["data"], "Response should contain data")

	responseDataFromResponse, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Data in response should be map[string]interface{}")

	assert.Equal(t, responseData.PurchasePriceId, int(responseDataFromResponse["purchase_price_id"].(float64)), "Purchase price ID should match")
	// Add other assertions as needed
}

func TestAddPurchasePrice_Success(t *testing.T) {
	payload := masteritempayloads.PurchasePriceDetailRequest{
		// Adjust with appropriate payload
		PurchasePriceDetailId: 1,
		PurchasePriceId:       1,
		ItemId:                1,
		PurchasePrice:         1,
		IsActive:              true,
	}

	mockService := new(MockPurchasePriceService)
	mockService.On("AddPurchasePrice", payload).
		Return((*exceptions.BaseErrorResponse)(nil))

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8000/v1/purchase-price/detail", bytes.NewReader(payloadBytes))
	rr := httptest.NewRecorder()

	controller := masteritemcontroller.NewPurchasePriceController(mockService)
	controller.AddPurchasePrice(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode, "Status code should be 200")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.Equal(t, "Create Data Successfully!", response["message"], "Message should match expected")
}

func TestDeletePurchasePrice_Success(t *testing.T) {
	id := 1 // Change to the appropriate ID
	req, _ := http.NewRequest("DELETE", "http://localhost:8000/v1/purchase-price/all/detail/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()

	mockService := new(MockPurchasePriceService)
	mockService.On("DeletePurchasePrice", id).
		Return((*exceptions.BaseErrorResponse)(nil))

	controller := masteritemcontroller.NewPurchasePriceController(mockService)
	controller.DeletePurchasePrice(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode, "Status code should be 200")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.Equal(t, "Purchase Price deleted successfully", response["message"], "Message should match expected")
}

///////////////////////////////////////////////////////////////////////////////////////////////////
