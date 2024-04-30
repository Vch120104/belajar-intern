package test

import (
	"after-sales/api/config"
	mastercontroller "after-sales/api/controllers/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	"fmt"
	"reflect"
	"runtime"

	"after-sales/api/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIncentiveMasterService struct {
	mock.Mock
}

// Mock the GetAllIncentiveMaster method
func (m *MockIncentiveMasterService) GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptionsss_test.BaseErrorResponse)
}

// Mock the GetAllIncentiveById method
func (m *MockIncentiveMasterService) GetIncentiveMasterById(id int) (masterpayloads.IncentiveMasterResponse, *exceptionsss_test.BaseErrorResponse) {
	args := m.Called(id)
	return args.Get(0).(masterpayloads.IncentiveMasterResponse), args.Get(1).(*exceptionsss_test.BaseErrorResponse)
}

// Mock the ChangeStatusIncentiveMaster method
func (m *MockIncentiveMasterService) ChangeStatusIncentiveMaster(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	args := m.Called(Id)
	return args.Bool(0), args.Get(1).(*exceptionsss_test.BaseErrorResponse)
}

// Mock the SaveIncentiveMaster method
func (m *MockIncentiveMasterService) SaveIncentiveMaster(payload masterpayloads.IncentiveMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	args := m.Called(payload)
	return args.Bool(0), args.Get(1).(*exceptionsss_test.BaseErrorResponse)
}

func TestGetAllIncentiveMaster_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8000/v1/incentive", nil)
	rr := httptest.NewRecorder()

	// Simulate a successful response from the service
	responseData := []map[string]interface{}{
		{"key": "value"},
	}
	mockService := new(MockIncentiveMasterService)
	mockService.On("GetAllIncentiveMaster", mock.Anything, mock.Anything).
		Return(responseData, len(responseData), len(responseData), (*exceptionsss_test.BaseErrorResponse)(nil)) // Return nil for BaseErrorResponse

	controller := mastercontroller.NewIncentiveMasterController(mockService)
	controller.GetAllIncentiveMaster(rr, req)

	// Check if the HTTP status code is as expected
	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode, "Status code should be 200")

	// Check if the response body contains the expected data
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err, "Error should be nil when unmarshalling response")

	// Print the response for debugging
	fmt.Println("Response:", response)

	// Check if the data key exists in the response
	assert.NotNil(t, response["data"], "Response should contain data")

	// Check if the data in the response has the expected type
	responseDataFromResponse, ok := response["data"].([]interface{})
	assert.True(t, ok, "Data in response should be []interface{}")

	// Check if the length of the data in the response matches the expected length
	assert.Equal(t, len(responseData), len(responseDataFromResponse), "Length of response data should match")
}

func TestGetIncentiveMasterById(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	IncentiveMasterRepository := masterrepositoryimpl.StartIncentiveMasterRepositoryImpl()
	IncentiveMasterService := masterserviceimpl.StartIncentiveMasterService(IncentiveMasterRepository, db, rdb)

	get, err := IncentiveMasterService.GetIncentiveMasterById(2) //change value test here
	if err != nil {
		t.Errorf("Error: %v", err) // Ubah t.Fatalf menjadi t.Errorf
		return                     // Kembalikan agar tes berhenti di sini jika terjadi kesalahan
	}

	assert.NotZero(t, get.IncentiveLevelId, "Expected non-zero result for IncentiveLevelId")
}

func BenchmarkGetIncentiveMasterById(b *testing.B) {
	fmt.Println("goos:", runtime.GOOS)
	fmt.Println("goarch:", runtime.GOARCH)
	fmt.Println("pkg:", reflect.TypeOf(b).PkgPath())
	fmt.Println("cpu:", cpuInfo())

	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	IncentiveMasterRepository := masterrepositoryimpl.StartIncentiveMasterRepositoryImpl()
	IncentiveMasterService := masterserviceimpl.StartIncentiveMasterService(IncentiveMasterRepository, db, rdb)

	b.ResetTimer()

	b.Run("GetIncentiveMasterById", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := IncentiveMasterService.GetIncentiveMasterById(3)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	})
}

func cpuInfo() string {
	info, _ := cpu.Info()
	if len(info) > 0 {
		return info[0].ModelName
	}
	return "Unknown"
}
