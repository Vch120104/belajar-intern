package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ShiftScheduleRepository interface {
	GetAllShiftSchedule(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetShiftScheduleById(*gorm.DB, int) (masterpayloads.ShiftScheduleResponse, *exceptions.BaseErrorResponse)
	SaveShiftSchedule(*gorm.DB, masterpayloads.ShiftScheduleResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusShiftSchedule(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	// GetShiftScheduleByCode(*gorm.DB, string) (masterpayloads.ShiftScheduleResponse, error)
	// GetAllShiftScheduleIsActive(*gorm.DB) ([]masterpayloads.ShiftScheduleResponse, error)
	GetShiftScheduleDropDown(tx *gorm.DB) ([]masterpayloads.ShiftScheduleDropDownResponse, *exceptions.BaseErrorResponse)
	UpdateShiftSchedule(tx *gorm.DB, Id int, request masterpayloads.ShiftScheduleUpdate) (masterentities.ShiftSchedule, *exceptions.BaseErrorResponse)
}
