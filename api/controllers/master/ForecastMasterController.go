package mastercontroller

import (

	// "after-sales/api/middlewares"

	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ForecastMasterController interface {
	GetForecastMasterById(writer http.ResponseWriter, request *http.Request)
	SaveForecastMaster(writer http.ResponseWriter, request *http.Request)
	ChangeStatusForecastMaster(writer http.ResponseWriter, request *http.Request)
	GetAllForecastMaster(writer http.ResponseWriter, request *http.Request)
	UpdateForecastMaster(writer http.ResponseWriter, request *http.Request)
}
type ForecastMasterControllerImpl struct {
	ForecastMasterService masterservice.ForecastMasterService
}

func NewForecastMasterController(forecastMasterService masterservice.ForecastMasterService) ForecastMasterController {
	return &ForecastMasterControllerImpl{
		ForecastMasterService: forecastMasterService,
	}
}

// @Summary Get Forecast Master By Id
// @Description REST API Forecast Master
// @Accept json
// @Produce json
// @Tags Master : Forecast Master
// @Param forecast_master_id path int true "forecast_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/forecast-master/{forecast_master_id} [get]
func (r *ForecastMasterControllerImpl) GetForecastMasterById(writer http.ResponseWriter, request *http.Request) {

	ForecastMasterId, errA := strconv.Atoi(chi.URLParam(request, "forecast_master_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ForecastMasterService.GetForecastMasterById(int(ForecastMasterId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Forecast Master
// @Description REST API Forecast Master
// @Accept json
// @Produce json
// @Tags Master : Forecast Master
// @param reqBody body masterpayloads.ForecastMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/forecast-master/ [post]
func (r *ForecastMasterControllerImpl) SaveForecastMaster(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.ForecastMasterResponse
	helper.ReadFromRequestBody(request, &formRequest)
	result, err := r.ForecastMasterService.SaveForecastMaster(formRequest)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Create Data Successfully!", http.StatusOK)
}

// @Summary Change Status Forecast Master
// @Description REST API Forecast Master
// @Accept json
// @Produce json
// @Tags Master : Forecast Master
// @param forecast_master_id path int true "forecast_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/forecast-master/{forecast_master_id} [patch]
func (r *ForecastMasterControllerImpl) ChangeStatusForecastMaster(writer http.ResponseWriter, request *http.Request) {

	forecast_master_id, errA := strconv.Atoi(chi.URLParam(request, "forecast_master_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ForecastMasterService.ChangeStatusForecastMaster(int(forecast_master_id))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get All Forecast Master
// @Description REST API Forecast Master
// @Accept json
// @Produce json
// @Tags Master : Forecast Master
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param supplier_name query string false "supplier_name"
// @Param moving_code_description query string false "moving_code_description"
// @Param order_type_name query string false "order_type_name"
// @Param forecast_master_lead_time query float64 false "forecast_master_lead_time"
// @Param forecast_master_safety_factor query float64 false "forecast_master_safety_factor"
// @Param forecast_master_order_cycle query float64 false "forecast_master_order_cycle"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/forecast-master/ [get]

func (r *ForecastMasterControllerImpl) GetAllForecastMaster(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"supplier_name": queryValues.Get("supplier_name"),
		"mtr_moving_code.moving_code_description":           queryValues.Get("moving_code_description"),
		"order_type_name":                                   queryValues.Get("order_type_name"),
		"mtr_forecast_master.forecast_master_lead_time":     queryValues.Get("forecast_master_lead_time"),
		"mtr_forecast_master.forecast_master_safety_factor": queryValues.Get("forecast_master_safety_factor"),
		"mtr_forecast_master.forecast_master_order_cycle":   queryValues.Get("forecast_master_order_cycle"),
		"mtr_forecast_master.is_active":                     queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}
	print(queryParams)

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.ForecastMasterService.GetAllForecastMaster(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *ForecastMasterControllerImpl) UpdateForecastMaster(writer http.ResponseWriter, request *http.Request) {
	forecast_master_id, errA := strconv.Atoi(chi.URLParam(request, "forecast_master_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	var formRequest masterpayloads.ForecastMasterResponse
	helper.ReadFromRequestBody(request, &formRequest)
	result, err := r.ForecastMasterService.UpdateForecastMaster(formRequest, forecast_master_id)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}
