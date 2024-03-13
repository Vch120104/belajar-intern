package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UnitOfMeasurementController interface {
	GetAllUnitOfMeasurement(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllUnitOfMeasurementIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetUnitOfMeasurementByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveUnitOfMeasurement(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusUnitOfMeasurement(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type UnitOfMeasurementControllerImpl struct {
	unitofmeasurementservice masteritemservice.UnitOfMeasurementService
}

func NewUnitOfMeasurementController(UnitOfMeasurementService masteritemservice.UnitOfMeasurementService) UnitOfMeasurementController {
	return &UnitOfMeasurementControllerImpl{
		unitofmeasurementservice: UnitOfMeasurementService,
	}
}

// @Summary Get All Unit Of Measurement
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param uom_code query string false "uom_code"
// @Param uom_description query string false "uom_description"
// @Param uom_type_desc query string false "uom_type_desc"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement [get]
func (r *UnitOfMeasurementControllerImpl) GetAllUnitOfMeasurement(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_uom.is_active":          queryValues.Get("is_active"),
		"mtr_uom.uom_code":           queryValues.Get("uom_code"),
		"mtr_uom.uom_description":    queryValues.Get("uom_description"),
		"mtr_uom_type.uom_type_desc": queryValues.Get("uom_type_desc"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.unitofmeasurementservice.GetAllUnitOfMeasurement(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Unit Of Measurement drop down
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement-drop-down [get]
func (r *UnitOfMeasurementControllerImpl) GetAllUnitOfMeasurementIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	result := r.unitofmeasurementservice.GetAllUnitOfMeasurementIsActive()

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Unit Of Measurement By Code
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Param uom_code path string true "uom_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement-by-code/{uom_code} [get]
func (r *UnitOfMeasurementControllerImpl) GetUnitOfMeasurementByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	operationGroupCode := params.ByName("uom_code")
	result := r.unitofmeasurementservice.GetUnitOfMeasurementByCode(operationGroupCode)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Unit Of Measurement
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @param reqBody body masteritempayloads.UomResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement [post]
func (r *UnitOfMeasurementControllerImpl) SaveUnitOfMeasurement(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masteritempayloads.UomResponse
	var message = ""

	helper.ReadFromRequestBody(request, &formRequest)

	create := r.unitofmeasurementservice.SaveUnitOfMeasurement(formRequest)

	if formRequest.UomId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Unit Of Measurement
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @param uom_id path int true "uom_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement/{uom_id} [patch]
func (r *UnitOfMeasurementControllerImpl) ChangeStatusUnitOfMeasurement(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	uomId, _ := strconv.Atoi(params.ByName("uom_id"))

	response := r.unitofmeasurementservice.ChangeStatusUnitOfMeasurement(int(uomId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
