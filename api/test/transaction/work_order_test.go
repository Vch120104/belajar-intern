package test

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockWorkOrderService struct {
	mock.Mock
}

func (m *MockWorkOrderService) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(tx, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewStatus(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	args := m.Called(tx)
	return args.Get(0).([]transactionworkshopentities.WorkOrderMasterStatus), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewType(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	args := m.Called(tx)
	return args.Get(0).([]transactionworkshopentities.WorkOrderMasterType), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewBill(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
	args := m.Called(tx)
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderBillable), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewDropPoint(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse) {
	args := m.Called(tx)
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderDropPoint), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewVehicleBrand(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse) {
	args := m.Called(tx)
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderVehicleBrand), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) NewVehicleModel(tx *gorm.DB, brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse) {
	args := m.Called(tx, brandId)
	return args.Get(0).([]transactionworkshoppayloads.WorkOrderVehicleModel), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) GetById(id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse) {
	args := m.Called(id)
	return args.Get(0).(transactionworkshoppayloads.WorkOrderRequest), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(tx, request)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) Submit(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	args := m.Called(tx, Id)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) Void(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	args := m.Called(tx, Id)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	args := m.Called(tx, Id)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) VehicleLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) CampaignLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	args := m.Called(filterCondition, pages)
	return args.Get(0).([]map[string]interface{}), args.Int(1), args.Int(2), args.Get(3).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddRequest(workOrderId int, request transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse {
	args := m.Called(workOrderId, request)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteRequest(workOrderId int, requestId int) *exceptions.BaseErrorResponse {
	args := m.Called(workOrderId, requestId)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) AddVehicleService(workOrderId int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse {
	args := m.Called(workOrderId, request)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteVehicleService(workOrderId int, vehicleServiceId int) *exceptions.BaseErrorResponse {
	args := m.Called(workOrderId, vehicleServiceId)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) UpdateWorkOrder(id int, payload transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	args := m.Called(id, payload)
	return args.Bool(0), args.Get(1).(*exceptions.BaseErrorResponse)
}

func (m *MockWorkOrderService) DeleteWorkOrder(id int) *exceptions.BaseErrorResponse {
	args := m.Called(id)
	return args.Get(0).(*exceptions.BaseErrorResponse)
}
