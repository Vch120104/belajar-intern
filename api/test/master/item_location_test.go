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
	"github.com/xuri/excelize/v2"
)

type MockItemLocationService struct {
	mock.Mock
}

func (m *MockItemLocationService) GetAllItemLocationDetail(criteria []utils.FilterCondition, paginate pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	args := m.Called(criteria, paginate)
	return args.Get(0).(pagination.Pagination), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) PopupItemLocation(criteria []utils.FilterCondition, paginate pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(criteria, paginate)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) AddItemLocation(itemLocID int, payload masteritempayloads.ItemLocationDetailRequest) (masteritementities.ItemLocationDetail, *exceptions.BaseErrorResponse) {
	args := m.Called(itemLocID, payload)
	entity := args.Get(0).(masteritementities.ItemLocationDetail)
	var err *exceptions.BaseErrorResponse
	if args.Get(1) != nil {
		err = args.Get(1).(*exceptions.BaseErrorResponse)
	}
	return entity, err
}

func (m *MockItemLocationService) DeleteItemLocation(id int) *exceptions.BaseErrorResponse {
	args := m.Called(id)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) GetAllItemLoc(filtercondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	args := m.Called(filtercondition, pages)
	return args.Get(0).(pagination.Pagination), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) GetByIdItemLoc(id int) (masteritempayloads.ItemLocationGetByIdResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Get(0).(masteritempayloads.ItemLocationGetByIdResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) SaveItemLoc(req masteritempayloads.SaveItemlocation) (masteritementities.ItemLocation, *exceptions.BaseErrorResponse) {
	args := m.Called(req)
	return args.Get(0).(masteritementities.ItemLocation), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) DeleteItemLoc(ids []int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(ids)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	args := m.Called()
	return args.Get(0).(*excelize.File), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) UploadPreviewFile(rows [][]string) ([]masteritempayloads.UploadItemLocationResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(rows)
	return args.Get(0).([]masteritempayloads.UploadItemLocationResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockItemLocationService) UploadProcessFile(uploadPreview []masteritempayloads.UploadItemLocationResponse) ([]masteritementities.ItemLocation, *exceptions.BaseErrorResponse) {
	args := m.Called(uploadPreview)
	return args.Get(0).([]masteritementities.ItemLocation), args.Get(1).(*exceptions.BaseErrorResponse)
}

func TestPopupItemLocation_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8000/v1/item-location/popup-location", nil)
	rr := httptest.NewRecorder()

	responseData := []map[string]interface{}{
		{"key": "value"}, // Sesuaikan dengan data yang diharapkan
	}
	mockService := new(MockItemLocationService)
	mockService.On("PopupItemLocation", mock.Anything, mock.Anything).
		Return(responseData, len(responseData), len(responseData), (*exceptions.BaseErrorResponse)(nil))

	controller := masteritemcontroller.NewItemLocationController(mockService)
	controller.PopupItemLocation(rr, req)

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

func TestAddItemLocation_Success(t *testing.T) {
	itemLocID := 1 // Ganti dengan ID yang sesuai
	payload := masteritempayloads.ItemLocationDetailRequest{
		// Sesuaikan dengan payload yang sesuai
		ItemLocationDetailId: 1,
		ItemLocationId:       1,
		ItemId:               1,
		ItemLocationSourceId: 1,
	}

	mockService := new(MockItemLocationService)
	mockService.On("AddItemLocation", itemLocID, payload).
		Return((*exceptions.BaseErrorResponse)(nil))

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8000/v1/item-location/"+strconv.Itoa(itemLocID)+"/detail", bytes.NewReader(payloadBytes))
	rr := httptest.NewRecorder()

	controller := masteritemcontroller.NewItemLocationController(mockService)
	controller.AddItemLocation(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusCreated, statusCode, "Status code should be 201")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.Equal(t, "Item location added successfully", response["message"], "Message should match expected")
}

func TestDeleteItemLocation_Success(t *testing.T) {
	id := 1 // Ganti dengan ID yang sesuai
	req, _ := http.NewRequest("DELETE", "http://localhost:8000/v1/item-location/detail/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()

	mockService := new(MockItemLocationService)
	mockService.On("DeleteItemLocation", id).
		Return((*exceptions.BaseErrorResponse)(nil))

	controller := masteritemcontroller.NewItemLocationController(mockService)
	controller.DeleteItemLocation(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode, "Status code should be 200")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	assert.Equal(t, "Item location deleted successfully", response["message"], "Message should match expected")
}

func TestGenerateTemplateFile_Success(t *testing.T) {
	f := excelize.NewFile()
	req, _ := http.NewRequest("GET", "http://localhost:8000/v1/item-location/download-template", nil)
	rr := httptest.NewRecorder()

	mockService := new(MockItemLocationService)
	mockService.On("GenerateTemplateFile").
		Return(f, (*exceptions.BaseErrorResponse)(nil))

	controller := masteritemcontroller.NewItemLocationController(mockService)
	controller.DownloadTemplate(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode, "Status code should be 200")
}
