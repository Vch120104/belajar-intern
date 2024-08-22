package transactionjpcbrepository

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type TechnicianAttendanceRepository interface {
	GetAllTechnicianAttendance(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveTechnicianAttendance(tx *gorm.DB, req transactionjpcbpayloads.TechnicianAttendanceSaveRequest) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse)
	ChangeStatusTechnicianAttendance(tx *gorm.DB, technicianAttendanceId int) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse)
}
