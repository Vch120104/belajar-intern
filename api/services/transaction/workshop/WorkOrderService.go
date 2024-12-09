package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type WorkOrderService interface {

	// Lookup Function
	GenerateDocumentNumber(workOrderId int) (string, *exceptions.BaseErrorResponse)

	// normal function
	New(request transactionworkshoppayloads.WorkOrderNormalRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderResponseDetail, *exceptions.BaseErrorResponse)
	Save(request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	Submit(Id int) (bool, string, *exceptions.BaseErrorResponse)
	Void(workOrderId int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(Id int) (bool, *exceptions.BaseErrorResponse)

	// Service Request
	GetAllRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetRequestById(workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderServiceResponse, *exceptions.BaseErrorResponse)
	UpdateRequest(workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse)
	AddRequest(workorderID int, requests transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse)
	AddRequestMultiId(workorderID int, requests []transactionworkshoppayloads.WorkOrderServiceRequest) ([]transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse)
	DeleteRequest(workorderID int, detailID int) (bool, *exceptions.BaseErrorResponse)
	DeleteRequestMultiId(workorderID int, detailID []int) (bool, *exceptions.BaseErrorResponse)

	// Service Vehicle
	GetAllVehicleService(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleServiceById(workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderServiceVehicleResponse, *exceptions.BaseErrorResponse)
	UpdateVehicleService(workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse)
	AddVehicleService(int, transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse)
	DeleteVehicleService(int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteVehicleServiceMultiId(workorderID int, detailID []int) (bool, *exceptions.BaseErrorResponse)

	// detail work order
	GetAllDetailWorkOrder(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDetailByIdWorkOrder(workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderDetailResponse, *exceptions.BaseErrorResponse)
	UpdateDetailWorkOrder(workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	AddDetailWorkOrder(int, transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	DeleteDetailWorkOrder(int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteDetailWorkOrderMultiId(workorderID int, detailID []int) (bool, *exceptions.BaseErrorResponse)

	// booking function
	NewBooking(request transactionworkshoppayloads.WorkOrderBookingRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	GetAllBooking(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBookingById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderBookingResponse, *exceptions.BaseErrorResponse)
	SaveBooking(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse)

	// affiliate function
	NewAffiliated(workOrderId int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllAffiliated(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAffiliatedById(workOrderId int, id int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderAffiliateResponse, *exceptions.BaseErrorResponse)
	SaveAffiliated(workOrderId int, id int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse)

	// support function
	DeleteCampaign(workOrderId int) (transactionworkshoppayloads.DeleteCampaignPayload, *exceptions.BaseErrorResponse)
	AddContractService(workOrderId int, request transactionworkshoppayloads.WorkOrderContractServiceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	AddGeneralRepairPackage(workOrderId int, request transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	AddFieldAction(workOrderId int, request transactionworkshoppayloads.WorkOrderFieldActionRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	ChangeBillTo(workOrderId int, request transactionworkshoppayloads.ChangeBillToRequest) (transactionworkshoppayloads.ChangeBillToResponse, *exceptions.BaseErrorResponse)
	ChangePhoneNo(workOrderId int, request transactionworkshoppayloads.ChangePhoneNoRequest) (*transactionworkshoppayloads.ChangePhoneNoResponse, *exceptions.BaseErrorResponse)
	ConfirmPrice(workOrderId int, idwos []int, request transactionworkshoppayloads.WorkOrderConfirmPriceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)

	GetServiceRequestByWO(workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetClaimByWO(workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetClaimItemByWO(workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetWOByBillCode(workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetDetailWOByClaimBillCode(workOrderId int, transactionTypeId int, atpmClaimNumber string, pages pagination.Pagination) ([]transactionworkshoppayloads.GetClaimResponsePayload, *exceptions.BaseErrorResponse)
	GetDetailWOByBillCode(workOrderId int, transactionTypeId int, pages pagination.Pagination) ([]transactionworkshoppayloads.GetClaimResponsePayload, *exceptions.BaseErrorResponse)
	GetDetailWOByATPMBillCode(workOrderId int, transactionTypeId int, pages pagination.Pagination) ([]transactionworkshoppayloads.GetClaimResponsePayload, *exceptions.BaseErrorResponse)
	GetSupplyByWO(workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
