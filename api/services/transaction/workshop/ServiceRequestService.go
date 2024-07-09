package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type ServiceRequestService interface {
	GenerateDocumentNumberServiceRequest(ServiceRequestId int) (string, *exceptions.BaseErrorResponse)

	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(id int) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse)
	New(request transactionworkshoppayloads.ServiceRequestSaveRequest) (bool, *exceptions.BaseErrorResponse)
	Save(id int, request transactionworkshoppayloads.ServiceRequestSaveRequest) (bool, *exceptions.BaseErrorResponse)
	Submit(id int) (bool, string, *exceptions.BaseErrorResponse)
	Void(id int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(id int) (bool, *exceptions.BaseErrorResponse)

	GetAllServiceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetServiceDetailById(idsys int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse)
	AddServiceDetail(idsys int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateServiceDetail(idsys int, idservice int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (bool, *exceptions.BaseErrorResponse)
	DeleteServiceDetail(idsys int, idservice int) (bool, *exceptions.BaseErrorResponse)
	DeleteServiceDetailMultiId(idsys int, idservice []int) (bool, *exceptions.BaseErrorResponse)
}
