package transactionjpcbrepositoryimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"net/http"
	"strconv"
	"strings"
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
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil
}

func (t *TechnicianAttendanceImpl) GetAddLineTechnician(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var serviceDate string
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "service_date") {
			serviceDate = filter.ColumnValue
			break
		}
	}

	entities := transactionjpcbentities.TechnicianAttendance{}
	userIds := []int{}
	err := tx.Model(&entities).Where("service_date = ?", serviceDate).Pluck("user_id", &userIds).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching technician attendance data",
			Err:        err,
		}
	}

	params := generalserviceapiutils.UserDetailParams{
		Page:        pages.Page,
		Limit:       10000000,
		UserIdNotIn: utils.IntSliceToString(userIds),
		RoleName:    "Technician",
	}
	employeeResponse, _ := generalserviceapiutils.GetAllUserDetail(params)

	result, totalPages, totalRows := pagination.NewDataFramePaginate(employeeResponse, &pages)
	pages.Rows = result
	pages.TotalPages = totalPages
	pages.TotalRows = int64(totalRows)

	return pages, nil
}

func (t *TechnicianAttendanceImpl) SaveTechnicianAttendance(tx *gorm.DB, req transactionjpcbpayloads.TechnicianAttendanceSaveRequest) ([]transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse) {
	serviceDate := req.ServiceDate.Truncate(24 * time.Hour)
	userIdStr := strings.Split(req.UserIds, ",")

	entities := []transactionjpcbentities.TechnicianAttendance{}
	for _, userIdStr := range userIdStr {
		userId, convErr := strconv.Atoi(userIdStr)
		if convErr != nil {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error converting user_id into integer",
				Err:        convErr,
			}
		}

		_, employeeErr := generalserviceapiutils.GetEmployeeMasterById(userId)
		if employeeErr != nil {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching employee data",
				Err:        employeeErr.Err,
			}
		}

		entity := transactionjpcbentities.TechnicianAttendance{
			CompanyId:   req.CompanyId,
			ServiceDate: serviceDate,
			UserId:      userId,
			Attendance:  true,
		}
		entities = append(entities, entity)
	}

	err := tx.Create(&entities).Error

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
