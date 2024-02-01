package mastercontroller

import (

	// "after-sales/api/middlewares"

	"after-sales/api/payloads"
	masterservice "after-sales/api/services/master"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ForecastMasterController interface {
	GetForecastMasterById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
type ForecastMasterControllerImpl struct {
	ForecastMasterService masterservice.ForecastMasterService
}

func NewForecastMasterController(forecastMasterService masterservice.ForecastMasterService) ForecastMasterController {
	return &ForecastMasterControllerImpl{
		ForecastMasterService: forecastMasterService,
	}
}

// @Summary Get Operation Group By Id
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
