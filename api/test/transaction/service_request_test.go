package test

import (
	"after-sales/api/config"
	transactionworkshopcontroller "after-sales/api/controllers/transactions/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"after-sales/api/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockServiceRequestService struct {
	mock.Mock
}

// Mock the SaveServiceRequest method
func (m *MockServiceRequestService) New(payload transactionworkshoppayloads.ServiceRequestSaveRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(payload)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Mock the GetAllServiceRequest method
func (m *MockServiceRequestService) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

// Mock the GetServiceRequestById method
func (m *MockServiceRequestService) GetById(id int) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Get(0).(transactionworkshoppayloads.ServiceRequestResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

var (
	dbservice             *gorm.DB
	rdbservice            *redis.Client
	ServiceRequestService transactionworkshopservice.ServiceRequestService
)

func TestMain(m *testing.M) {
	config.InitEnvConfigs(true, "")
	dbservice = config.InitDB()
	rdbservice = config.InitRedis()

	// Run the tests
	code := m.Run()

	// Close the database and Redis connections
	sqlDB, _ := dbservice.DB()
	sqlDB.Close()
	rdbservice.Close()

	// Exit with the appropriate code
	os.Exit(code)
}

func TestSaveServiceRequest_Success(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()

	ServiceRequestRepository := transactionworkshoprepositoryimpl.OpenServiceRequestRepositoryImpl()
	ServiceRequestService := transactionworkshopserviceimpl.OpenServiceRequestServiceImpl(ServiceRequestRepository, db, rdb)

	// Prepare a valid payload with an existing service request number
	payload := transactionworkshoppayloads.ServiceRequestSaveRequest{
		ServiceRequestSystemNumber: 3, // Adjust with a valid existing service request number
		CompanyId:                  1,
		ReplyId:                    3,
	}
	payloadBytes, _ := json.Marshal(payload)

	// Set the URL path with the correct service_request_system_number
	req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/service-request/%d", payload.ServiceRequestSystemNumber), bytes.NewReader(payloadBytes))
	rr := httptest.NewRecorder()

	controller := transactionworkshopcontroller.NewServiceRequestController(ServiceRequestService)
	controller.Save(rr, req)

	// Check HTTP status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	// Parse response JSON
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
		return
	}

	// Check for 'message' field in response
	message, ok := response["message"].(string)
	if !ok {
		t.Error("Expected 'message' field in response")
		return
	}

	// Assert the success message
	assert.Equal(t, "Save Data Successfully", message, "Expected success message")
}

func TestGetAllServiceRequest_Success(t *testing.T) {
	// Initialize dependencies
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	ServiceRequestRepository := transactionworkshoprepositoryimpl.OpenServiceRequestRepositoryImpl()
	ServiceRequestService := transactionworkshopserviceimpl.OpenServiceRequestServiceImpl(ServiceRequestRepository, db, rdb)

	req, _ := http.NewRequest("GET", "/v1/service-request", nil)
	rr := httptest.NewRecorder()

	// Create controller and invoke GetAll method
	controller := transactionworkshopcontroller.NewServiceRequestController(ServiceRequestService)
	controller.GetAll(rr, req)

	// Check HTTP status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	// Parse response JSON
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
		return
	}

	fmt.Println("Response:", response)

	// Check for expected message
	assert.Equal(t, "Get Data Successfully", response["message"], "Expected success message")
}

func TestGetServiceRequestById_Success(t *testing.T) {
	// Initialize dependencies
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	ServiceRequestRepository := transactionworkshoprepositoryimpl.OpenServiceRequestRepositoryImpl()
	ServiceRequestService := transactionworkshopserviceimpl.OpenServiceRequestServiceImpl(ServiceRequestRepository, db, rdb)

	// Prepare request
	req, _ := http.NewRequest("GET", "/v1/service-request/1", nil)
	rr := httptest.NewRecorder()

	// Create controller and invoke GetById method
	controller := transactionworkshopcontroller.NewServiceRequestController(ServiceRequestService)
	controller.GetById(rr, req)

	// Check HTTP status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	// Parse response JSON
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return
	}

	fmt.Println("Response:", response)

	// Check for expected 'data' field
	data, dataExists := response["data"].(map[string]interface{})
	if !dataExists {
		t.Error("Expected 'data' field to be present and of type map[string]interface{}")
		return
	}

	// Example assertion: Check if service_request_system_number is present and non-zero
	serviceRequestNumber, ok := data["service_request_system_number"].(float64) // Assuming it's a float64
	if !ok || serviceRequestNumber == 0 {
		t.Error("Expected non-zero result for service_request_system_number")
		return
	}

}

func BenchmarkGetServiceRequestById(b *testing.B) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()

	ServiceRequestRepository := transactionworkshoprepositoryimpl.OpenServiceRequestRepositoryImpl()
	ServiceRequestService := transactionworkshopserviceimpl.OpenServiceRequestServiceImpl(ServiceRequestRepository, db, rdb)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ServiceRequestService.GetById(1)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}
