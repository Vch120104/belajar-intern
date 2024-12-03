package transactionbodyshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ServiceBodyshopService interface {
	GetAllByTechnicianWOBodyshop(idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	StartService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
	PendingService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
	TransferService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
	StopService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse)
}
