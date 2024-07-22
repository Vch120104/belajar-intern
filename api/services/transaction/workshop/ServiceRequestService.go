package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type ServiceRequestService interface {
	GenerateDocumentNumberServiceRequest(ServiceRequestId int) (string, *exceptions.BaseErrorResponse)
	NewStatus(filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterStatus, *exceptions.BaseErrorResponse)

	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse)
	New(request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
	Save(id int, request transactionworkshoppayloads.ServiceRequestSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
	Submit(id int) (bool, string, *exceptions.BaseErrorResponse)
	Void(id int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(id int) (bool, *exceptions.BaseErrorResponse)

	GetAllServiceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetServiceDetailById(idsys int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse)
	AddServiceDetail(idsys int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse)
	UpdateServiceDetail(idsys int, idservice int, request transactionworkshoppayloads.ServiceDetailUpdateRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse)
	DeleteServiceDetail(idsys int, idservice int) (bool, *exceptions.BaseErrorResponse)
	DeleteServiceDetailMultiId(idsys int, idservice []int) (bool, *exceptions.BaseErrorResponse)
}
