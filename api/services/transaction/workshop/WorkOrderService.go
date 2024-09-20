package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type WorkOrderService interface {

	// support function
	NewStatus(filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse)
	AddStatus(request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateStatus(id int, request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse)
	DeleteStatus(id int) (bool, *exceptions.BaseErrorResponse)

	NewType(filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse)
	AddType(request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateType(id int, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse)
	DeleteType(id int) (bool, *exceptions.BaseErrorResponse)

	NewBill() ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse)
	AddBill(request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateBill(id int, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse)
	DeleteBill(id int) (bool, *exceptions.BaseErrorResponse)

	NewTrxType() ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse)
	AddTrxType(request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse)
	UpdateTrxType(id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse)
	DeleteTrxType(id int) (bool, *exceptions.BaseErrorResponse)

	NewTrxTypeSo() ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse)
	AddTrxTypeSo(request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse)
	UpdateTrxTypeSo(id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse)
	DeleteTrxTypeSo(id int) (bool, *exceptions.BaseErrorResponse)

	NewDropPoint() ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse)

	NewVehicleBrand() ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse)
	NewVehicleModel(brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse)

	// Lookup Function
	GenerateDocumentNumber(workOrderId int) (string, *exceptions.BaseErrorResponse)

	// normal function
	New(request transactionworkshoppayloads.WorkOrderNormalRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderResponseDetail, *exceptions.BaseErrorResponse)
	Save(request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	Submit(Id int) (bool, string, *exceptions.BaseErrorResponse)
	Void(workOrderId int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(Id int) (bool, *exceptions.BaseErrorResponse)

	// Service Request
	GetAllRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetRequestById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse)
	UpdateRequest(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse)
	AddRequest(int, transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse)
	DeleteRequest(int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteRequestMultiId(idwosn int, idwos []int) (bool, *exceptions.BaseErrorResponse)

	// Service Vehicle
	GetAllVehicleService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleServiceById(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse)
	UpdateVehicleService(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse)
	AddVehicleService(int, transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse)
	DeleteVehicleService(int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteVehicleServiceMultiId(idwosn int, idwos []int) (bool, *exceptions.BaseErrorResponse)

	// detail work order
	GetAllDetailWorkOrder(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetDetailByIdWorkOrder(idwosn int, idwos int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse)
	UpdateDetailWorkOrder(idwosn int, idwos int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	AddDetailWorkOrder(int, transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	DeleteDetailWorkOrder(int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteDetailWorkOrderMultiId(idwosn int, idwos []int) (bool, *exceptions.BaseErrorResponse)

	// booking function
	NewBooking(request transactionworkshoppayloads.WorkOrderBookingRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	GetAllBooking(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetBookingById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderBookingResponse, *exceptions.BaseErrorResponse)
	SaveBooking(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse)

	// affiliate function
	NewAffiliated(workOrderId int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllAffiliated(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAffiliatedById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderAffiliateResponse, *exceptions.BaseErrorResponse)
	SaveAffiliated(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse)

	DeleteCampaign(workOrderId int) (transactionworkshoppayloads.DeleteCampaignPayload, *exceptions.BaseErrorResponse)
	ChangeBillTo(workOrderId int, request transactionworkshoppayloads.ChangeBillToRequest) (bool, *exceptions.BaseErrorResponse)
	ChangePhoneNo(workOrderId int, request transactionworkshoppayloads.ChangePhoneNoRequest) (bool, *exceptions.BaseErrorResponse)
	ConfirmPrice(workOrderId int, idwos []int) (bool, *exceptions.BaseErrorResponse)
}
