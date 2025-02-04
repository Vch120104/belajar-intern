package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderRepository interface {
	New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, workorderID int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderResponseDetail, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, workorderID int) (bool, string, *exceptions.BaseErrorResponse)
	Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(tx *gorm.DB, workorderID int) (bool, *exceptions.BaseErrorResponse)

	GetAllRequest(tx *gorm.DB, filterCondition []utils.FilterCondition, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetRequestById(tx *gorm.DB, workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderServiceResponse, *exceptions.BaseErrorResponse)
	UpdateRequest(tx *gorm.DB, workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse)
	AddRequest(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse)
	AddRequestMultiId(tx *gorm.DB, workorderID int, requests []transactionworkshoppayloads.WorkOrderServiceRequest) ([]transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse)
	DeleteRequest(tx *gorm.DB, workorderID int, detailID int) (bool, *exceptions.BaseErrorResponse)
	DeleteRequestMultiId(tx *gorm.DB, workorderID int, detailID []int) (bool, *exceptions.BaseErrorResponse)

	GetAllVehicleService(tx *gorm.DB, filterCondition []utils.FilterCondition, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleServiceById(tx *gorm.DB, workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderServiceVehicleResponse, *exceptions.BaseErrorResponse)
	UpdateVehicleService(tx *gorm.DB, workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse)
	AddVehicleService(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse)
	DeleteVehicleService(tx *gorm.DB, workorderID int, detailID int) (bool, *exceptions.BaseErrorResponse)
	DeleteVehicleServiceMultiId(tx *gorm.DB, workorderID int, detailID []int) (bool, *exceptions.BaseErrorResponse)

	GetAllDetailWorkOrder(tx *gorm.DB, filterCondition []utils.FilterCondition, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDetailByIdWorkOrder(tx *gorm.DB, workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderDetailResponse, *exceptions.BaseErrorResponse)
	UpdateDetailWorkOrder(tx *gorm.DB, workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	AddDetailWorkOrder(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	DeleteDetailWorkOrder(tx *gorm.DB, workorderID int, detailID int) (bool, *exceptions.BaseErrorResponse)
	DeleteDetailWorkOrderMultiId(tx *gorm.DB, workorderID int, detailID []int) (bool, *exceptions.BaseErrorResponse)

	NewBooking(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderBookingRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse)
	GetAllBooking(tx *gorm.DB, filterCondition []utils.FilterCondition, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBookingById(tx *gorm.DB, workorderID int, detailID int, page pagination.Pagination) (transactionworkshoppayloads.WorkOrderBookingResponse, *exceptions.BaseErrorResponse)
	SaveBooking(tx *gorm.DB, int, workorderID int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse)

	NewAffiliated(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllAffiliated(tx *gorm.DB, filterCondition []utils.FilterCondition, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAffiliatedById(tx *gorm.DB, workorderID int, detailID int, page pagination.Pagination) (transactionworkshoppayloads.WorkOrderAffiliateResponse, *exceptions.BaseErrorResponse)
	SaveAffiliated(tx *gorm.DB, workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse)

	GenerateDocumentNumber(tx *gorm.DB, workorderID int) (string, *exceptions.BaseErrorResponse)
	CalculateWorkOrderTotal(tx *gorm.DB, workorderID int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	DeleteCampaign(tx *gorm.DB, workorderID int) (transactionworkshoppayloads.DeleteCampaignPayload, *exceptions.BaseErrorResponse)
	AddContractService(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderContractServiceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	AddGeneralRepairPackage(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	AddFieldAction(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderFieldActionRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)
	ChangeBillTo(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.ChangeBillToRequest) (transactionworkshoppayloads.ChangeBillToResponse, *exceptions.BaseErrorResponse)
	ChangePhoneNo(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.ChangePhoneNoRequest) (*transactionworkshoppayloads.ChangePhoneNoResponse, *exceptions.BaseErrorResponse)
	ConfirmPrice(tx *gorm.DB, workorderID int, detailID []int, request transactionworkshoppayloads.WorkOrderConfirmPriceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse)

	GetOperationItemById(LineTypeStr string, OperationItemId int) (interface{}, *exceptions.BaseErrorResponse)
	HandleLineTypeResponse(lineTypeCode string, operationItemResponse interface{}) (string, string, *exceptions.BaseErrorResponse)
}
