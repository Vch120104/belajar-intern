package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type ShiftScheduleRepositoryImpl struct {
}

func StartShiftScheduleRepositoryImpl() masterrepository.ShiftScheduleRepository {
	return &ShiftScheduleRepositoryImpl{}
}

func (*ShiftScheduleRepositoryImpl) GetAllShiftSchedule(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterentities.ShiftSchedule{}
	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if len(entities) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

// func (*ShiftScheduleRepositoryImpl) GetAllShiftScheduleIsActive(tx *gorm.DB) ([]masterpayloads.ShiftScheduleResponse, error) {
// 	var ShiftSchedules []masterentities.ShiftSchedule
// 	response := []masterpayloads.ShiftScheduleResponse{}

// 	err := tx.Model(&ShiftSchedules).Where("is_active = 'true'").Scan(&response).Error

// 	if err != nil {
// 		return response, err
// 	}

// 	return response, nil
// }

func (*ShiftScheduleRepositoryImpl) GetShiftScheduleById(tx *gorm.DB, Id int) (masterpayloads.ShiftScheduleResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.ShiftSchedule{}
	response := masterpayloads.ShiftScheduleResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.ShiftSchedule{
			ShiftScheduleId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

// func (*ShiftScheduleRepositoryImpl) GetShiftScheduleByCode(tx *gorm.DB, Code string) (masterpayloads.ShiftScheduleResponse, error) {
// 	entities := masterentities.ShiftSchedule{}
// 	response := masterpayloads.ShiftScheduleResponse{}

// 	rows, err := tx.Model(&entities).
// 		Where(masterentities.ShiftSchedule{
// 			ShiftScheduleCode: Code,
// 		}).
// 		First(&response).
// 		Rows()

// 	if err != nil {
// 		return response, err
// 	}

// 	defer rows.Close()

// 	return response, nil
// }

func (*ShiftScheduleRepositoryImpl) SaveShiftSchedule(tx *gorm.DB, req masterpayloads.ShiftScheduleResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.ShiftSchedule{
		IsActive:        req.IsActive,
		ShiftScheduleId: req.ShiftScheduleId,
		ShiftCode:       req.ShiftCode,
		CompanyId:       req.CompanyId,
		ShiftGroup:      req.ShiftGroup,
		EffectiveDate:   req.EffectiveDate,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		RestStartTime:   req.RestStartTime,
		RestEndTime:     req.RestEndTime,
		Monday:          req.Monday,
		Tuesday:         req.Tuesday,
		Wednesday:       req.Wednesday,
		Thursday:        req.Thursday,
		Friday:          req.Friday,
		Saturday:        req.Saturday,
		Sunday:          req.Sunday,
		Manpower:        req.Manpower,
		ManpowerBooking: req.ManpowerBooking,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	return true, nil
}

func (*ShiftScheduleRepositoryImpl) ChangeStatusShiftSchedule(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterentities.ShiftSchedule

	result := tx.Model(&entities).
		Where("shift_schedule_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
