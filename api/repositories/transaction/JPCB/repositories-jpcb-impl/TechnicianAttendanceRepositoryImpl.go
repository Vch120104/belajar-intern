package transactionjpcbrepositoryimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type TechnicianAttendanceImpl struct {
}

func StartTechnicianAttendanceRepositoryImpl() transactionjpcbrepository.TechnicianAttendanceRepository {
	return &TechnicianAttendanceImpl{}
}

func (t *TechnicianAttendanceImpl) GetAllTechnicianAttendance(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.TechnicianAttendance{}
	responses := []transactionjpcbpayloads.TechnicianAttendanceGetAllResponse{}

	baseModelQuery := tx.Model(&entities)
	whereQuery := utils.ApplyFilterExact(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil
}

func (t *TechnicianAttendanceImpl) SaveTechnicianAttendance(tx *gorm.DB, req transactionjpcbpayloads.TechnicianAttendanceSaveRequest) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse) {
	serviceDate := req.ServiceDate.Truncate(24 * time.Hour)
	entities := transactionjpcbentities.TechnicianAttendance{
		CompanyId:   req.CompanyId,
		ServiceDate: serviceDate,
		UserId:      req.UserId,
		Attendance:  req.Attendance,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return entities, nil
}

func (t *TechnicianAttendanceImpl) ChangeStatusTechnicianAttendance(tx *gorm.DB, technicianAttendanceId int) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.TechnicianAttendance{}

	err := tx.Model(&entities).Where(transactionjpcbentities.TechnicianAttendance{TechnicianAttendanceId: technicianAttendanceId}).First(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entities.Attendance = !entities.Attendance

	err = tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}
