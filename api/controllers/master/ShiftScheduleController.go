package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ShiftScheduleController interface {
	GetAllShiftSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// GetAllShiftScheduleIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// GetShiftScheduleByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveShiftSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusShiftSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
type ShiftScheduleControllerImpl struct {
	ShiftScheduleService masterservice.ShiftScheduleService
}

func NewShiftScheduleController(ShiftScheduleService masterservice.ShiftScheduleService) ShiftScheduleController {
	return &ShiftScheduleControllerImpl{
		ShiftScheduleService: ShiftScheduleService,
	}
}

// @Summary Get All Shift Schedule
// @Description REST API Shift Schedule
// @Accept json
// @Produce json
// @Tags Master : Shift Schedule
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param shift_schedule_code query string false "shift_schedule_code"
// @Param shift_schedule_description query string false "shift_schedule_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/shift-schedule [get]
func (r *ShiftScheduleControllerImpl) GetAllShiftSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryParams := map[string]string{
		"shift_schedule_code":        params.ByName("shift_schedule_code"),
		"shift_schedule_description": params.ByName("shift_schedule_description"),
		"is_active":                   params.ByName("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(params, "limit"),
		Page:   utils.NewGetQueryInt(params, "page"),
		SortOf: params.ByName("sort_of"),
		SortBy: params.ByName("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.ShiftScheduleService.GetAllShiftSchedule(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// // @Summary Get All Shift Schedule drop down
// // @Description REST API Shift Schedule
// // @Accept json
// // @Produce json
// // @Tags Master : Shift Schedule
// // @Success 200 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.Error
// // @Router /aftersales-service/api/aftersales/shift-schedule/drop-down [get]
// func (r *ShiftScheduleControllerImpl) GetAllShiftScheduleIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

// 	result := r.ShiftScheduleService.GetAllShiftScheduleIsActive()

// 	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
// }

// // @Summary Get Shift Schedule By Code
// // @Description REST API Shift Schedule
// // @Accept json
// // @Produce json
// // @Tags Master : Shift Schedule
// // @Param shift_schedule_code path string true "shift_schedule_code"
// // @Success 200 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.Error
// // @Router /aftersales-service/api/aftersales/shift-schedule/by-code/{shift_schedule_code} [get]
// func (r *ShiftScheduleControllerImpl) GetShiftScheduleByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

// 	ShiftScheduleCode := params.ByName("shift_schedule_code")

// 	result := r.ShiftScheduleService.GetShiftScheduleByCode(ShiftScheduleCode)

// 	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
// }

// @Summary Save Shift Schedule
// @Description REST API Shift Schedule
// @Accept json
// @Produce json
// @Tags Master : Shift Schedule
// @param reqBody body masterpayloads.ShiftScheduleResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/shift-schedule [post]
func (r *ShiftScheduleControllerImpl) SaveShiftSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masterpayloads.ShiftScheduleResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.ShiftScheduleService.SaveShiftSchedule(formRequest)

	if formRequest.ShiftScheduleId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Shift Schedule
// @Description REST API Shift Schedule
// @Accept json
// @Produce json
// @Tags Master : Shift Schedule
// @param shift_schedule_id path int true "shift_schedule_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/shift-schedule/{shift_schedule_id} [patch]
func (r *ShiftScheduleControllerImpl) ChangeStatusShiftSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	ShiftScheduleId, _ := strconv.Atoi(params.ByName("shift_schedule_id"))

	response := r.ShiftScheduleService.ChangeStatusShiftSchedule(int(ShiftScheduleId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
