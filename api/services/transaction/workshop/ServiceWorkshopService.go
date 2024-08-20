package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type ServiceWorkshopService interface {
	GetAllByTechnicianWO(idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ServiceWorkshopDetailResponse, *exceptions.BaseErrorResponse)
	StartService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
	PendingService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
	TransferService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
	StopService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
}
