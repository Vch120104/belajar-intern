package transactionjpcbservice

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
)

type TechnicianAttendanceService interface {
	GetAllTechnicianAttendance(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAddLineTechnician(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveTechnicianAttendance(req transactionjpcbpayloads.TechnicianAttendanceSaveRequest) ([]transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse)
	ChangeStatusTechnicianAttendance(technicianAttendanceId int) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse)
}
