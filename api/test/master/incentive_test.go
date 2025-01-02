package test

import (
	"after-sales/api/config"
	mastercontroller "after-sales/api/controllers/master"
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	"bytes"
	"fmt"

	"after-sales/api/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIncentiveMasterService struct {
	mock.Mock
}

// Mock the GetAllIncentiveMaster method
func (m *MockIncentiveMasterService) GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).(pagination.Pagination), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Mock the GetAllIncentiveById method
func (m *MockIncentiveMasterService) GetIncentiveMasterById(id int) (masterpayloads.IncentiveMasterResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Get(0).(masterpayloads.IncentiveMasterResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Mock the ChangeStatusIncentiveMaster method
func (m *MockIncentiveMasterService) ChangeStatusIncentiveMaster(id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Get(0).(masterentities.IncentiveMaster), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockIncentiveMasterService) SaveIncentiveMaster(payload masterpayloads.IncentiveMasterRequest) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse) {
	args := m.Called(payload)
	return args.Get(0).(masterentities.IncentiveMaster), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockIncentiveMasterService) UpdateIncentiveMaster(req masterpayloads.IncentiveMasterRequest, id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse) {
	args := m.Called(req, id)
	return args.Get(0).(masterentities.IncentiveMaster), args.Get(1).(*exceptions.BaseErrorResponse)
}
func TestSaveIncentiveMaster_Success(t *testing.T) {
	payload := masterpayloads.IncentiveMasterRequest{
		IncentiveLevelId:      1,
		IncentiveLevelCode:    6,
		JobPositionId:         1,
		IncentiveLevelPercent: 10,
	}

	mockService := new(MockIncentiveMasterService)
	mockService.On("SaveIncentiveMaster", payload).Return(true, (*exceptions.BaseErrorResponse)(nil))

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8000/v1/incentive", bytes.NewReader(payloadBytes))
	rr := httptest.NewRecorder()

	controller := mastercontroller.NewIncentiveMasterController(mockService)
	controller.SaveIncentiveMaster(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)

	// Check whether this is a new data creation or an update of existing data
	expectedStatusCode := http.StatusOK
	expectedMessage := "Update Data Successfully!"

	if payload.IncentiveLevelId == 0 {
		expectedStatusCode = http.StatusCreated
		expectedMessage = "Create Data Successfully!"
	}

	assert.Equal(t, expectedStatusCode, statusCode, "Status code should match expected")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	fmt.Println("Response:", response)

	// Check message and status code
	assert.Equal(t, expectedMessage, response["message"], "Message should match expected")
	assert.Equal(t, expectedStatusCode, int(response["status_code"].(float64)), "Status code should match expected")
}

func TestSaveIncentiveMaster_Failure(t *testing.T) {
	payload := masterpayloads.IncentiveMasterRequest{
		IncentiveLevelId:      0, // ganti value disini
		IncentiveLevelCode:    1,
		JobPositionId:         1,
		IncentiveLevelPercent: 10,
	}

	// Simulate a failure response from the service
	mockService := new(MockIncentiveMasterService)
	mockService.On("SaveIncentiveMaster", payload).
		Return(false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest, // Mengganti status kode menjadi 400
			Err:        fmt.Errorf("some error"),
		})

	// Create a new HTTP request
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8000/v1/incentive", bytes.NewReader(payloadBytes))
	rr := httptest.NewRecorder()

	// Call the controller
	controller := mastercontroller.NewIncentiveMasterController(mockService)
	controller.SaveIncentiveMaster(rr, req)

	// Check if the HTTP status code is as expected
	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusBadRequest, statusCode, "Status code should be 400") // Mengubah ekspektasi menjadi 400

	// Check if the response body contains the expected data
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	// Print the response for debugging
	fmt.Println("Response:", response)

	assert.NotNil(t, response["error"], "Response should contain error")
	assert.Equal(t, "some error", response["error"], "Error should be 'some error'")
}

func TestGetIncentiveMasterById(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	IncentiveMasterRepository := masterrepositoryimpl.StartIncentiveMasterRepositoryImpl()
	IncentiveMasterService := masterserviceimpl.StartIncentiveMasterService(IncentiveMasterRepository, db, rdb)

	get, err := IncentiveMasterService.GetIncentiveMasterById(1) // ganti value disini
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	assert.NotZero(t, get.IncentiveLevelId, "Expected non-zero result for IncentiveLevelId")
}

func BenchmarkGetIncentiveMasterById(b *testing.B) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	IncentiveMasterRepository := masterrepositoryimpl.StartIncentiveMasterRepositoryImpl()
	IncentiveMasterService := masterserviceimpl.StartIncentiveMasterService(IncentiveMasterRepository, db, rdb)

	b.ResetTimer()

	b.Run("GetIncentiveMasterById", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := IncentiveMasterService.GetIncentiveMasterById(1)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	})
}
