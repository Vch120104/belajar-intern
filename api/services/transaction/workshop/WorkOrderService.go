package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderService interface {

	// support function
	NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse)
	NewType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse)
	NewBill(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse)
	NewDropPoint(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse)
	NewVehicleBrand(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse)
	NewVehicleModel(tx *gorm.DB, brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse)

	// Lookup Function
	VehicleLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CampaignLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)

	// normal function
	New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse)
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest, workOrderId int) (bool, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) (bool, string, *exceptions.BaseErrorResponse)
	Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)

	// Service Request
	GetAllRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetRequestById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse)
	UpdateRequest(tx *gorm.DB, idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse
	AddRequest(int, transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse
	DeleteRequest(int, int) *exceptions.BaseErrorResponse

	// Service Vehicle
	GetAllVehicleService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleServiceById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse)
	UpdateVehicleService(tx *gorm.DB, idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse
	AddVehicleService(int, transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse
	DeleteVehicleService(int, int) *exceptions.BaseErrorResponse

	// detail work order
	GetAllDetailWorkOrder(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetDetailByIdWorkOrder(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse)
	UpdateDetailWorkOrder(tx *gorm.DB, idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderDetailRequest) *exceptions.BaseErrorResponse
	AddDetailWorkOrder(int, transactionworkshoppayloads.WorkOrderDetailRequest) *exceptions.BaseErrorResponse
	DeleteDetailWorkOrder(int, int) *exceptions.BaseErrorResponse

	// booking function
	NewBooking(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllBooking(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetBookingById(workOrderId int, id int) (transactionworkshoppayloads.WorkOrderBookingRequest, *exceptions.BaseErrorResponse)
	SaveBooking(tx *gorm.DB, workOrderId int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse)
	SubmitBooking(tx *gorm.DB, workOrderId int, Id int) (bool, *exceptions.BaseErrorResponse)
	VoidBooking(tx *gorm.DB, workOrderId int, Id int) (bool, *exceptions.BaseErrorResponse)
	CloseBooking(tx *gorm.DB, workOrderId int, Id int) (bool, *exceptions.BaseErrorResponse)
}
