package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ServiceRequestRepository interface {
	GenerateDocumentNumberServiceRequest(tx *gorm.DB, ServiceRequestId int) (string, *exceptions.BaseErrorResponse)
	NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterStatus, *exceptions.BaseErrorResponse)
	NewServiceType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterServiceType, *exceptions.BaseErrorResponse)

	GetAll(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(*gorm.DB, int, pagination.Pagination) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse)
	New(*gorm.DB, transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
	Save(*gorm.DB, int, transactionworkshoppayloads.ServiceRequestSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
	Submit(*gorm.DB, int) (bool, string, *exceptions.BaseErrorResponse)
	Void(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)

	GetAllServiceDetail(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetServiceDetailById(*gorm.DB, int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse)
	AddServiceDetail(*gorm.DB, int, transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse)
	UpdateServiceDetail(*gorm.DB, int, int, transactionworkshoppayloads.ServiceDetailUpdateRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse)
	DeleteServiceDetail(*gorm.DB, int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteServiceDetailMultiId(*gorm.DB, int, []int) (bool, *exceptions.BaseErrorResponse)
}
