package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ShiftScheduleRepository interface {
	GetAllShiftSchedule(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetShiftScheduleById(*gorm.DB, int) (masterpayloads.ShiftScheduleResponse, error)
	SaveShiftSchedule(*gorm.DB, masterpayloads.ShiftScheduleResponse) (bool, error)
	ChangeStatusShiftSchedule(*gorm.DB, int) (bool, error)
	// GetShiftScheduleByCode(*gorm.DB, string) (masterpayloads.ShiftScheduleResponse, error)
	// GetAllShiftScheduleIsActive(*gorm.DB) ([]masterpayloads.ShiftScheduleResponse, error)
}
