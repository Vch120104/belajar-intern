package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ShiftScheduleRepository interface {
	GetAllShiftSchedule(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetShiftScheduleById(*gorm.DB, int) (masterpayloads.ShiftScheduleResponse, *exceptionsss_test.BaseErrorResponse)
	SaveShiftSchedule(*gorm.DB, masterpayloads.ShiftScheduleResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusShiftSchedule(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	// GetShiftScheduleByCode(*gorm.DB, string) (masterpayloads.ShiftScheduleResponse, error)
	// GetAllShiftScheduleIsActive(*gorm.DB) ([]masterpayloads.ShiftScheduleResponse, error)
}
