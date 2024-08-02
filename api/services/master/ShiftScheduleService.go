package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ShiftScheduleService interface {
	GetAllShiftSchedule(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	// GetAllShiftScheduleIsActive() []masterpayloads.ShiftScheduleResponse
	GetShiftScheduleById(int) (masterpayloads.ShiftScheduleResponse, *exceptions.BaseErrorResponse)
	// GetShiftScheduleByCode(string) masterpayloads.ShiftScheduleResponse
	ChangeStatusShiftSchedule(int) (bool, *exceptions.BaseErrorResponse)
	SaveShiftSchedule(masterpayloads.ShiftScheduleResponse) (bool, *exceptions.BaseErrorResponse)
}
