package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ShiftScheduleService interface {
	GetAllShiftSchedule(filterCondition []utils.FilterCondition, pages pagination.Pagination)(pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	// GetAllShiftScheduleIsActive() []masterpayloads.ShiftScheduleResponse
	GetShiftScheduleById(int) (masterpayloads.ShiftScheduleResponse,*exceptionsss_test.BaseErrorResponse)
	// GetShiftScheduleByCode(string) masterpayloads.ShiftScheduleResponse
	ChangeStatusShiftSchedule(int) (bool,*exceptionsss_test.BaseErrorResponse)
	SaveShiftSchedule(masterpayloads.ShiftScheduleResponse) (bool,*exceptionsss_test.BaseErrorResponse)
}
