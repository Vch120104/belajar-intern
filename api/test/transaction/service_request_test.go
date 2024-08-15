package test

import (
	"after-sales/api/config"
	transactionworkshopcontroller "after-sales/api/controllers/transactions/workshop"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads"
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
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockServiceRequestService struct {
	mock.Mock
}

// Mock the methods
func (m *MockServiceRequestService) New(request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Get(0).(transactionworkshopentities.ServiceRequest), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) NewStatus(filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterStatus, *exceptions.BaseErrorResponse) {
	args := m.Called(filter)
	return args.Get(0).([]transactionworkshopentities.ServiceRequestMasterStatus), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) Save(id int, request transactionworkshoppayloads.ServiceRequestSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Get(0).(transactionworkshopentities.ServiceRequest), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) Submit(id int) (bool, string, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.String(1), args.Get(2).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(id, pages)
	return args.Get(0).(transactionworkshoppayloads.ServiceRequestResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) Void(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) CloseOrder(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) AddServiceDetail(idsys int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {
	args := m.Called(idsys, request)
	return args.Get(0).(transactionworkshopentities.ServiceRequestDetail), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) UpdateServiceDetail(idsys int, idservice int, request transactionworkshoppayloads.ServiceDetailUpdateRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {
	args := m.Called(idsys, idservice, request)
	return args.Get(0).(transactionworkshopentities.ServiceRequestDetail), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) DeleteServiceDetail(idsys int, idservice int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idsys, idservice)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) DeleteServiceDetailMultiId(idsys int, idservice []int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idsys, idservice)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) GenerateDocumentNumberServiceRequest(ServiceRequestId int) (string, *exceptions.BaseErrorResponse) {
	args := m.Called(ServiceRequestId)
	return args.String(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) GetAllServiceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockServiceRequestService) GetServiceDetailById(idsys int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(idsys)
	return args.Get(0).(transactionworkshoppayloads.ServiceDetailResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

func setup() (*gorm.DB, *redis.Client, transactionworkshopservice.ServiceRequestService) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	ServiceRequestRepository := transactionworkshoprepositoryimpl.OpenServiceRequestRepositoryImpl()
	ServiceRequestService := transactionworkshopserviceimpl.OpenServiceRequestServiceImpl(ServiceRequestRepository, db, rdb)
	return db, rdb, ServiceRequestService
}

func TestSaveServiceRequest_Success(t *testing.T) {
	db, rdb, ServiceRequestService := setup()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	// Get the current time
	currentTime := time.Now().UTC()
	// Set serviceDate to a future date
	serviceDate := currentTime.Add(24 * time.Hour) // 1 day in the future
	serviceRequestDate := currentTime

	payload := transactionworkshoppayloads.ServiceRequestSaveRequest{
		BookingSystemNumber:      0,
		BrandId:                  23,
		CompanyId:                21,
		DealerRepresentativeId:   3,
		EstimationSystemNumber:   0,
		ModelId:                  1,
		ProfitCenterId:           1,
		ReferenceDocSystemNumber: 1,
		ReferenceJobType:         "Maintenance",
		ReferenceTypeId:          6,
		ReplyId:                  7,
		ServiceCompanyId:         8,
		ServiceDate:              serviceDate,
		ServiceProfitCenterId:    2,
		ServiceRemark:            "Routine check-up",
		ServiceRequestBy:         "Administrator",
		ServiceRequestDate:       serviceRequestDate,
		ServiceRequestStatusId:   1,
		ServiceTypeId:            0,
		VariantId:                1,
		VehicleId:                4,
		WorkOrderSystemNumber:    1,
	}

	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/service-request", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller := transactionworkshopcontroller.NewServiceRequestController(ServiceRequestService)
	controller.New(rr, req)

	// Debugging information
	t.Logf("Response Status Code: %d", rr.Code)
	t.Logf("Response Body: %s", rr.Body.String())

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status Created")

	var response payloads.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
		return
	}

	if response.Message != "Create Data Successfully" {
		t.Errorf("Expected 'Create Data Successfully', got %v", response.Message)
	}

	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, response.StatusCode)
	}

	if response.Data == nil {
		t.Error("Expected response data to be non-nil, but got nil")
	}
}

func TestGetAllServiceRequest_Success(t *testing.T) {
	db, rdb, ServiceRequestService := setup()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	req, _ := http.NewRequest("GET", "/v1/service-request", nil)
	rr := httptest.NewRecorder()

	controller := transactionworkshopcontroller.NewServiceRequestController(ServiceRequestService)
	controller.GetAll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
		return
	}

	fmt.Println("Response:", response)

	assert.Equal(t, "Get Data Successfully", response["message"], "Expected success message")
}

func TestGetServiceRequestById_Success(t *testing.T) {
	db, rdb, ServiceRequestService := setup()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	req, _ := http.NewRequest("GET", "/v1/service-request/11", nil)
	rr := httptest.NewRecorder()

	controller := transactionworkshopcontroller.NewServiceRequestController(ServiceRequestService)
	controller.GetById(rr, req)

	pagination := pagination.Pagination{
		Limit: 10,
		Page:  0,
	}

	result, _ := ServiceRequestService.GetById(1, pagination)

	fmt.Println(result)

}

func BenchmarkGetServiceRequestById(b *testing.B) {
	db, rdb, ServiceRequestService := setup()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	paginate := pagination.Pagination{
		Limit:  10,
		Page:   1,
		SortOf: "",
		SortBy: "",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ServiceRequestService.GetById(11, paginate)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}
