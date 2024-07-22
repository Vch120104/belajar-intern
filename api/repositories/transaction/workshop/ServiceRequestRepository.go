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

	GetAll(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(*gorm.DB, int) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse)
	New(*gorm.DB, transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
	Save(*gorm.DB, int, transactionworkshoppayloads.ServiceRequestSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse)
	Submit(*gorm.DB, int) (bool, string, *exceptions.BaseErrorResponse)
	Void(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)

	GetAllServiceDetail(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetServiceDetailById(*gorm.DB, int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse)
	AddServiceDetail(*gorm.DB, int, transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse)
	UpdateServiceDetail(*gorm.DB, int, int, transactionworkshoppayloads.ServiceDetailUpdateRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse)
	DeleteServiceDetail(*gorm.DB, int, int) (bool, *exceptions.BaseErrorResponse)
	DeleteServiceDetailMultiId(*gorm.DB, int, []int) (bool, *exceptions.BaseErrorResponse)
}
