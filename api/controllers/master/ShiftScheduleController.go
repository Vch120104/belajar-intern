package mastercontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"errors"

	// "after-sales/api/middlewares"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ShiftScheduleController interface {
	GetAllShiftSchedule(writer http.ResponseWriter, request *http.Request)
	// GetAllShiftScheduleIsActive(writer http.ResponseWriter, request *http.Request)
	// GetShiftScheduleByCode(writer http.ResponseWriter, request *http.Request)
	SaveShiftSchedule(writer http.ResponseWriter, request *http.Request)
	ChangeStatusShiftSchedule(writer http.ResponseWriter, request *http.Request)
	GetShiftScheduleById(writer http.ResponseWriter, request *http.Request)
	GetShiftScheduleDropdown(writer http.ResponseWriter, request *http.Request)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/shift-schedule/ [get]
func (r *ShiftScheduleControllerImpl) GetAllShiftSchedule(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"is_active":         queryValues.Get("is_active"),
		"company_id":        queryValues.Get("company_id"),
		"shift_schedule_id": queryValues.Get("shift_schedule_id"),
		"shift_code":        queryValues.Get("shift_code"),
		"effective_date":    queryValues.Get("effective_date"),
		"shift_group":       queryValues.Get("shift_group"),
		"start_time":        queryValues.Get("start_time"),
		"end_time":          queryValues.Get("end_time"),
		"rest_start_time":   queryValues.Get("rest_start_time"),
		"rest_end_time":     queryValues.Get("rest_end_time"),
		"monday":            queryValues.Get("monday"),
		"tuesday":           queryValues.Get("tuesday"),
		"wednesday":         queryValues.Get("wednesday"),
		"thursday":          queryValues.Get("thursday"),
		"friday":            queryValues.Get("friday"),
		"saturday":          queryValues.Get("saturday"),
		"sunday":            queryValues.Get("sunday"),
		"manpower":          queryValues.Get("manpower"),
		"manpower_booking":  queryValues.Get("manpower_booking"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.ShiftScheduleService.GetAllShiftSchedule(filterCondition, pagination)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// // @Summary Get All Shift Schedule drop down
// // @Description REST API Shift Schedule
// // @Accept json
// // @Produce json
// // @Tags Master : Shift Schedule
// // @Success 200 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// // @Router /aftersales-service/api/aftersales/shift-schedule/drop-down [get]
// func (r *ShiftScheduleControllerImpl) GetAllShiftScheduleIsActive(writer http.ResponseWriter, request *http.Request) {

// 	result := r.ShiftScheduleService.GetAllShiftScheduleIsActive()

// 	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
// }

// @Summary Get Shift Schedule By Id
// @Description REST API Shift Schedule
// @Accept json
// @Produce json
// @Tags Master : Shift Schedule
// @Param shift_schedule_id path string true "shift_schedule_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/shift-schedule/{shift_schedule_id} [get]
func (r *ShiftScheduleControllerImpl) GetShiftScheduleById(writer http.ResponseWriter, request *http.Request) {

	ShiftScheduleId, errA := strconv.Atoi(chi.URLParam(request, "shift_schedule_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ShiftScheduleService.GetShiftScheduleById(ShiftScheduleId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Shift Schedule
// @Description REST API Shift Schedule
// @Accept json
// @Produce json
// @Tags Master : Shift Schedule
// @param reqBody body masterpayloads.ShiftScheduleResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/shift-schedule/ [post]
func (r *ShiftScheduleControllerImpl) SaveShiftSchedule(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.ShiftScheduleResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create, err := r.ShiftScheduleService.SaveShiftSchedule(formRequest)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/shift-schedule/{shift_schedule_id} [patch]
func (r *ShiftScheduleControllerImpl) ChangeStatusShiftSchedule(writer http.ResponseWriter, request *http.Request) {

	ShiftScheduleId, errA := strconv.Atoi(chi.URLParam(request, "shift_schedule_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ShiftScheduleService.ChangeStatusShiftSchedule(int(ShiftScheduleId))
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *ShiftScheduleControllerImpl) GetShiftScheduleDropdown(writer http.ResponseWriter, request *http.Request) {
	result, err := r.ShiftScheduleService.GetShiftScheduleDropDown()

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}
