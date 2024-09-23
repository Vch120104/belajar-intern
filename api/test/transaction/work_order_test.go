package test

import (
	"after-sales/api/config"
	transactionworkshopcontroller "after-sales/api/controllers/transactions/workshop"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"after-sales/api/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/zeebo/assert"
)

type MockWorkOrderService struct {
	mock.Mock
}

// Support Function
func (m *MockWorkOrderService) NewStatus(filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	args := m.Called(filter)
	return args.Get(0).([]transactionworkshopentities.WorkOrderMasterStatus), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddStatus(request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateStatus(id int, request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteStatus(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewType(filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	args := m.Called(filter)
	return args.Get(0).([]transactionworkshopentities.WorkOrderMasterType), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddType(request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateType(id int, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteType(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewLineType() ([]transactionworkshoppayloads.Linetype, *exceptions.BaseErrorResponse) {
	args := m.Called()
	return args.Get(0).([]transactionworkshoppayloads.Linetype), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddLineType(request transactionworkshoppayloads.Linetype) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateLineType(id int, request transactionworkshoppayloads.Linetype) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteLineType(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewBill() ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
	args := m.Called()
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderBillable), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddBill(request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateBill(id int, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteBill(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewTrxType() ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	args := m.Called()
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderTransactionType), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddTrxType(request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateTrxType(id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteTrxType(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewTrxTypeSo() ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	args := m.Called()
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderTransactionType), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddTrxTypeSo(request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateTrxTypeSo(id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteTrxTypeSo(id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewDropPoint() ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse) {
	args := m.Called()
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderDropPoint), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewVehicleBrand() ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse) {
	args := m.Called()
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderVehicleBrand), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewVehicleModel(brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse) {
	args := m.Called(brandId)
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderVehicleModel), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Lookup Function
func (m *MockWorkOrderService) GenerateDocumentNumber(workOrderId int) (string, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId)
	return args.String(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Normal Function
func (m *MockWorkOrderService) New(request transactionworkshoppayloads.WorkOrderNormalRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Get(0).(transactionworkshopentities.WorkOrder), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderResponseDetail, *exceptions.BaseErrorResponse) {
	args := m.Called(id, pages)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderResponseDetail), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) Save(request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {
	args := m.Called(request, workOrderId)
	return args.Get(0).(transactionworkshopentities.WorkOrder), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) Submit(Id int) (bool, string, *exceptions.BaseErrorResponse) {
	args := m.Called(Id)
	return args.Bool(0), args.String(1), args.Get(2).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) Void(workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) CloseOrder(Id int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(Id)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Service Request
func (m *MockWorkOrderService) GetAllRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetRequestById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderServiceRequest), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateRequest(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos, request)
	return args.Get(0).(transactionworkshopentities.WorkOrderRequestDescription), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddRequest(id int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Get(0).(transactionworkshopentities.WorkOrderRequestDescription), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteRequest(idwosn int, idwos int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteRequestMultiId(idwosn int, idwos []int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Service Vehicle
func (m *MockWorkOrderService) GetAllVehicleService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetVehicleServiceById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderServiceVehicleRequest), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateVehicleService(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos, request)
	return args.Get(0).(transactionworkshopentities.WorkOrderServiceVehicle), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddVehicleService(id int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Get(0).(transactionworkshopentities.WorkOrderServiceVehicle), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteVehicleService(idwosn int, idwos int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteVehicleServiceMultiId(idwosn int, idwos []int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Detail Work Order
func (m *MockWorkOrderService) GetAllDetailWorkOrder(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetDetailByIdWorkOrder(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderDetailRequest), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateDetailWorkOrder(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos, request)
	return args.Get(0).(transactionworkshopentities.WorkOrderDetail), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddDetailWorkOrder(id int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	args := m.Called(id, request)
	return args.Get(0).(transactionworkshopentities.WorkOrderDetail), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteDetailWorkOrder(idwosn int, idwos int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteDetailWorkOrderMultiId(idwosn int, idwos []int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(idwosn, idwos)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Booking Function
func (m *MockWorkOrderService) NewBooking(request transactionworkshoppayloads.WorkOrderBookingRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {
	args := m.Called(request)
	return args.Get(0).(transactionworkshopentities.WorkOrder), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetAllBooking(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetBookingById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderBookingResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, id, pages)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderBookingResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) SaveBooking(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Affiliate Function
func (m *MockWorkOrderService) NewAffiliated(workOrderId int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetAllAffiliated(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetAffiliatedById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderAffiliateResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, id, pages)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderAffiliateResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) SaveAffiliated(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, id, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) ChangeBillTo(workOrderId int, request transactionworkshoppayloads.ChangeBillToRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) ChangePhoneNo(workOrderId int, request transactionworkshoppayloads.ChangePhoneNoRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) ConfirmPrice(workOrderId int, idwos []int) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, idwos)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteCampaign(workOrderId int) (transactionworkshoppayloads.DeleteCampaignPayload, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId)
	return args.Get(0).(transactionworkshoppayloads.DeleteCampaignPayload), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddContractService(workOrderId int, request transactionworkshoppayloads.WorkOrderContractServiceRequest) (transactionworkshoppayloads.WorkOrderContractServiceResponse, *exceptions.BaseErrorResponse) {
	args := m.Called(workOrderId, request)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderContractServiceResponse), args.Get(1).(*exceptions.BaseErrorResponse)
}

// Get All Normal Work Order
func TestGetAllWorkOrder_Success(t *testing.T) {

	req, _ := http.NewRequest("GET", "http://localhost:8000/v1/work-order", nil)
	rr := httptest.NewRecorder()

	responseData := []map[string]interface{}{
		{
			"work_order_system_number":   float64(1),
			"work_order_document_number": "WSWO/N/09/24/00001",
			"formatted_work_order_date":  "2024-09-10 10:15:32",
			"work_order_type_id":         float64(1),
			"work_order_type_name":       "Normal",
			"status_id":                  float64(2),
			"status_name":                "New",
			"service_advisor_id":         float64(0),
			"brand_id":                   float64(23),
			"brand_name":                 "Nissan",
			"model_id":                   float64(56),
			"model_name":                 "GRAND LIVINA 1.5",
			"variant_id":                 float64(0),
			"service_site":               "",
			"vehicle_id":                 float64(4),
			"vehicle_code":               "B234ZTZT000000631",
			"vehicle_tnkb":               "B2502UOQ",
			"customer_id":                float64(0),
			"billto_customer_id":         float64(0),
			"repeated_job":               float64(0),
		},
	}

	mockService := new(MockWorkOrderService)
	mockService.On("GetAll", mock.Anything, mock.Anything).
		Return(responseData, len(responseData), len(responseData), (*exceptions.BaseErrorResponse)(nil)) // Return nil for BaseErrorResponse

	controller := transactionworkshopcontroller.NewWorkOrderController(mockService)
	controller.GetAll(rr, req)

	statusCode := rr.Code
	fmt.Println("Status code:", statusCode)
	assert.Equal(t, http.StatusOK, statusCode)

	fmt.Println("Response body:", rr.Body.String())

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, float64(200), response["status_code"])
	assert.Equal(t, "Get Data Successfully", response["message"])
	assert.Equal(t, float64(0), response["page"])
	assert.Equal(t, float64(0), response["page_limit"])
	assert.Equal(t, float64(len(responseData)), response["total_rows"])
	assert.Equal(t, float64(1), response["total_pages"])

	dataFromResponse, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, len(responseData), len(dataFromResponse))

	for i, item := range responseData {
		responseItemMap, ok2 := dataFromResponse[i].(map[string]interface{})
		assert.True(t, ok2)
		assert.Equal(t, item, responseItemMap)
	}

	mockService.AssertCalled(t, "GetAll", mock.Anything, mock.Anything)
}

// GetById Work Order
func TestGetWorkOrderById_Success(t *testing.T) {

	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	WorkorderRepository := transactionworkshoprepositoryimpl.OpenWorkOrderRepositoryImpl()
	WorkorderService := transactionworkshopserviceimpl.OpenWorkOrderServiceImpl(WorkorderRepository, db, rdb)

	pagination := pagination.Pagination{
		Page:       0,
		Limit:      10,
		TotalRows:  1,
		TotalPages: 1,
	}

	get, err := WorkorderService.GetById(1, pagination) // Handle both return values
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	assert.NotEqual(t, 0, get.WorkOrderSystemNumber)
}

func BenchmarkGetWorkOrderById(b *testing.B) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	WorkorderRepository := transactionworkshoprepositoryimpl.OpenWorkOrderRepositoryImpl()
	WorkorderService := transactionworkshopserviceimpl.OpenWorkOrderServiceImpl(WorkorderRepository, db, rdb)

	pagination := pagination.Pagination{
		Page:       0,
		Limit:      10,
		TotalRows:  1,
		TotalPages: 1,
	}

	b.ResetTimer()

	b.Run("GetById", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := WorkorderService.GetById(1, pagination)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	})
}
