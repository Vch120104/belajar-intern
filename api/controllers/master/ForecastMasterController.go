package mastercontroller

import (

	// "after-sales/api/middlewares"

	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	masterservice "after-sales/api/services/master"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ForecastMasterController interface {
	GetForecastMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusForecastMaster(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
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

	ForecastMasterId, _ := strconv.Atoi(params.ByName("forecast_master_id"))

	response := r.ForecastMasterService.ChangeStatusForecastMaster(int(ForecastMasterId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
