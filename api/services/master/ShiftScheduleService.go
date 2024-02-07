package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ShiftScheduleService interface {
	GetAllShiftSchedule(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	// GetAllShiftScheduleIsActive() []masterpayloads.ShiftScheduleResponse
	GetShiftScheduleById(int) masterpayloads.ShiftScheduleResponse
	// GetShiftScheduleByCode(string) masterpayloads.ShiftScheduleResponse
	ChangeStatusShiftSchedule(int) bool
	SaveShiftSchedule(masterpayloads.ShiftScheduleResponse) bool
}
