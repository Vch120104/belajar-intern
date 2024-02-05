package mastercontroller

import (

	// "after-sales/api/middlewares"

	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ForecastMasterController interface {
	GetForecastMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/forecast-master/{forecast_master_id} [get]
func (r *ForecastMasterControllerImpl) GetForecastMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	ForecastMasterId, _ := strconv.Atoi(params.ByName("forecast_master_id"))

	result := r.ForecastMasterService.GetForecastMasterById(int(ForecastMasterId))

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Forecast Master
// @Description REST API Forecast Master
// @Accept json
// @Produce json
// @Tags Master : Forecast Master
// @param reqBody body masterpayloads.ForecastMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/forecast-master [post]
func (r *ForecastMasterControllerImpl) SaveForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masterpayloads.ForecastMasterResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.ForecastMasterService.SaveForecastMaster(formRequest)

	if formRequest.ForecastMasterId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Forecast Master
// @Description REST API Forecast Master
// @Accept json
// @Produce json
// @Tags Master : Forecast Master
// @param forecast_master_id path int true "forecast_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/forecast-master/{forecast_master_id} [patch]
func (r *ForecastMasterControllerImpl) ChangeStatusForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	forecast_master_id, _ := strconv.Atoi(params.ByName("forecast_master_id"))

	response := r.ForecastMasterService.ChangeStatusForecastMaster(int(forecast_master_id))

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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/forecast-master [get]
func (r *ForecastMasterControllerImpl) GetAllForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryParams := map[string]string{
		"supplier_name": params.ByName("supplier_name"),
		"mtr_moving_code.moving_code_description":           params.ByName("moving_code_description"),
		"order_type_name":                                   params.ByName("order_type_name"),
		"mtr_forecast_master.forecast_master_lead_time":     params.ByName("forecast_master_lead_time"),
		"mtr_forecast_master.forecast_master_safety_factor": params.ByName("forecast_master_safety_factor"),
		"mtr_forecast_master.forecast_master_order_cycle":   params.ByName("forecast_master_order_cycle"),
		"mtr_forecast_master.is_active":                     params.ByName("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(params, "limit"),
		Page:   utils.NewGetQueryInt(params, "page"),
		SortOf: params.ByName("sort_of"),
		SortBy: params.ByName("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows := r.ForecastMasterService.GetAllForecastMaster(criteria, paginate)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}
