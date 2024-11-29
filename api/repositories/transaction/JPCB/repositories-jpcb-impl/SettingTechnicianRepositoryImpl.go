package transactionjpcbrepositoryimpl

import (
	"after-sales/api/config"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type SettingTechnicianRepositoryImpl struct {
}

func StartSettingTechnicianRepositoryImpl() transactionjpcbrepository.SettingTechnicianRepository {
	return &SettingTechnicianRepositoryImpl{}
}

func (r *SettingTechnicianRepositoryImpl) GetAllSettingTechnician(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.SettingTechnician{}
	responses := []transactionjpcbpayloads.SettingTechnicianPayload{}

	for i, filter := range filterCondition {
		if filter.ColumnField == "effective_date" {
			filterCondition[i].ColumnValue = utils.SafeConvertDateStrFormat(filter.ColumnValue)
		}
	}

	baseModelQuery := tx.Model(&entities)
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	mapResponses := []transactionjpcbpayloads.SettingTechnicianGetAllResponse{}
	for _, result := range responses {
		responsePayloads := transactionjpcbpayloads.SettingTechnicianGetAllResponse{
			SettingTechnicianSystemNumber: result.SettingTechnicianSystemNumber,
			EffectiveDate:                 result.EffectiveDate,
			SettingId:                     result.SettingTechnicianSystemNumber,
		}
		mapResponses = append(mapResponses, responsePayloads)
	}

	pages.Rows = mapResponses

	return pages, nil
}

func (r *SettingTechnicianRepositoryImpl) GetAllSettingTechnicianDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.SettingTechnicianDetail{}
	responses := []transactionjpcbpayloads.SettingTechnicianGetAllDetailPayload{}

	baseModelQuery := tx.Model(&entities).
		Select(`
			trx_setting_technician_detail.setting_technician_detail_system_number,
			trx_setting_technician_detail.setting_technician_system_number,
			trx_setting_technician_detail.technician_number,
			trx_setting_technician_detail.technician_employee_number_id,
			trx_setting_technician_detail.group_express,
			trx_setting_technician_detail.shift_group_id,
			mss.shift_group,
			trx_setting_technician_detail.is_booking`).
		Joins("INNER JOIN mtr_shift_schedule mss ON mss.shift_schedule_id = trx_setting_technician_detail.shift_group_id")
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	userEmployeeId := []int{}
	for _, result := range responses {
		userEmployeeId = append(userEmployeeId, result.TechnicianEmployeeNumberId)
	}
	userEmployeeId = utils.RemoveDuplicateIds(userEmployeeId)

	employeeResponse := []transactionjpcbpayloads.SettingTechnicianEmployeeResponse{}
	for _, Id := range userEmployeeId {
		employeeURL := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(Id)
		employeePayloads := transactionjpcbpayloads.SettingTechnicianEmployeeResponse{}
		if err := utils.Get(employeeURL, &employeePayloads, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if employeePayloads.UserEmployeeId == 0 {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("employee data not found"),
			}
		}
		employeeResponse = append(employeeResponse, employeePayloads)
	}

	mapResponses := []transactionjpcbpayloads.SettingTechnicianGetAllDetailResponse{}
	for _, result := range responses {
		employeeName := ""
		for _, employee := range employeeResponse {
			if result.TechnicianEmployeeNumberId == employee.UserEmployeeId {
				employeeName = employee.EmployeeName
			}
		}
		responsePayloads := transactionjpcbpayloads.SettingTechnicianGetAllDetailResponse{
			SettingTechnicianDetailSystemNumber: result.SettingTechnicianDetailSystemNumber,
			SettingTechnicianSystemNumber:       result.SettingTechnicianSystemNumber,
			TechnicianNumber:                    result.TechnicianNumber,
			UserEmployeeId:                      result.TechnicianEmployeeNumberId,
			EmployeeName:                        employeeName,
			GroupExpress:                        result.GroupExpress,
			ShiftGroupId:                        result.ShiftGroupId,
			ShiftGroup:                          result.ShiftGroup,
			IsBooking:                           result.IsBooking,
		}
		mapResponses = append(mapResponses, responsePayloads)
	}

	pages.Rows = mapResponses

	return pages, nil
}

func (r *SettingTechnicianRepositoryImpl) GetSettingTechnicianById(tx *gorm.DB, settingTechnicianId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.SettingTechnician{}
	responses := transactionjpcbpayloads.SettingTechnicianPayload{}

	err := tx.Model(&entities).Where(transactionjpcbentities.SettingTechnician{SettingTechnicianSystemNumber: settingTechnicianId}).First(&responses).Error

	result := transactionjpcbpayloads.SettingTechnicianGetByIdResponse{}

	if err != nil {
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	result.SettingTechnicianSystemNumber = responses.SettingTechnicianSystemNumber
	result.CompanyId = responses.CompanyId
	result.SettingId = responses.SettingTechnicianSystemNumber
	result.EffectiveDate = responses.EffectiveDate

	return result, nil
}

func (r *SettingTechnicianRepositoryImpl) GetSettingTechnicianDetailById(tx *gorm.DB, settingTechnicianDetailId int) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.SettingTechnicianDetail{}
	responses := transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse{}

	err := tx.Model(&entities).Where(transactionjpcbentities.SettingTechnicianDetail{SettingTechnicianDetailSystemNumber: settingTechnicianDetailId}).First(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return responses, nil
}

func (r *SettingTechnicianRepositoryImpl) GetSettingTechnicianByCompanyDate(tx *gorm.DB, companyId int, effectiveDate time.Time) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.SettingTechnician{}
	response := transactionjpcbpayloads.SettingTechnicianGetByIdResponse{}

	effectiveDate = time.Date(effectiveDate.Year(), effectiveDate.Month(), effectiveDate.Day(), 0, 0, 0, 0, effectiveDate.Location())

	err := tx.Model(&entities).
		Where(transactionjpcbentities.SettingTechnician{
			CompanyId:     companyId,
			EffectiveDate: &effectiveDate,
		}).
		First(&entities).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response.SettingTechnicianSystemNumber = entities.SettingTechnicianSystemNumber
	response.CompanyId = entities.CompanyId
	response.SettingId = entities.SettingTechnicianSystemNumber
	response.EffectiveDate = *entities.EffectiveDate

	return response, nil
}

func (r *SettingTechnicianRepositoryImpl) SaveSettingTechnician(tx *gorm.DB, CompanyId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse) {
	currentTime := time.Now().Truncate(24 * time.Hour)
	entities := transactionjpcbentities.SettingTechnician{
		CompanyId:     CompanyId,
		EffectiveDate: &currentTime,
	}
	responses := transactionjpcbpayloads.SettingTechnicianGetByIdResponse{}

	data := transactionjpcbentities.SettingTechnician{}
	err := tx.Model(&data).Where(entities).FirstOrCreate(&data).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	responses.SettingTechnicianSystemNumber = data.SettingTechnicianSystemNumber
	responses.CompanyId = data.CompanyId
	responses.SettingId = data.SettingTechnicianSystemNumber
	responses.EffectiveDate = *data.EffectiveDate

	return responses, nil
}

func (r *SettingTechnicianRepositoryImpl) SaveSettingTechnicianDetail(tx *gorm.DB, req transactionjpcbpayloads.SettingTechnicianDetailSaveRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.SettingTechnicianDetail{
		SettingTechnicianSystemNumbers: req.SettingTechnicianSystemNumber,
		TechnicianEmployeeNumberId:     req.TechnicianEmployeeNumberId,
		ShiftGroupId:                   req.ShiftGroupId,
		TechnicianNumber:               req.TechnicianNumber,
		IsBooking:                      req.IsBooking,
		GroupExpress:                   req.GroupExpress,
	}
	responses := transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse{}

	employeeURL := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(entities.TechnicianEmployeeNumberId)
	employeePayloads := transactionjpcbpayloads.SettingTechnicianEmployeeResponse{}
	if err := utils.Get(employeeURL, &employeePayloads, nil); err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if employeePayloads.UserEmployeeId == 0 {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("employee data not found"),
		}
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	responses.SettingTechnicianDetailSystemNumber = entities.SettingTechnicianDetailSystemNumber
	responses.ShiftGroupId = entities.ShiftGroupId
	responses.TechnicianNumber = entities.TechnicianNumber
	responses.IsBooking = entities.IsBooking
	responses.GroupExpress = entities.GroupExpress

	return responses, nil
}

func (r *SettingTechnicianRepositoryImpl) UpdateSettingTechnicianDetail(tx *gorm.DB, settingTechnicianDetailId int, req transactionjpcbpayloads.SettingTechnicianDetailUpdateRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse) {
	entities := transactionjpcbentities.SettingTechnicianDetail{}
	responses := transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse{}

	whereQuery := transactionjpcbentities.SettingTechnicianDetail{
		SettingTechnicianDetailSystemNumber: settingTechnicianDetailId,
	}

	err := tx.Model(&entities).Where(whereQuery).First(&entities).Error
	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entities.ShiftGroupId = req.ShiftGroupId
	entities.TechnicianNumber = req.TechnicianNumber
	entities.IsBooking = req.IsBooking
	entities.GroupExpress = req.GroupExpress

	err = tx.Save(&entities).Error
	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	responses.SettingTechnicianDetailSystemNumber = entities.SettingTechnicianDetailSystemNumber
	responses.ShiftGroupId = entities.ShiftGroupId
	responses.IsBooking = entities.IsBooking
	responses.TechnicianNumber = entities.TechnicianNumber
	responses.GroupExpress = entities.GroupExpress

	return responses, nil
}
