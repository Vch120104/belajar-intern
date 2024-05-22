package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UnitOfMeasurementController interface {
	GetAllUnitOfMeasurement(writer http.ResponseWriter, request *http.Request)
	GetAllUnitOfMeasurementIsActive(writer http.ResponseWriter, request *http.Request)
	GetUnitOfMeasurementByCode(writer http.ResponseWriter, request *http.Request)
	SaveUnitOfMeasurement(writer http.ResponseWriter, request *http.Request)
	ChangeStatusUnitOfMeasurement(writer http.ResponseWriter, request *http.Request)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/unit-of-measurement/ [get]
func (r *UnitOfMeasurementControllerImpl) GetAllUnitOfMeasurement(writer http.ResponseWriter, request *http.Request) {
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

	result, err := r.unitofmeasurementservice.GetAllUnitOfMeasurement(filterCondition, pagination)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Unit Of Measurement drop down
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/unit-of-measurement/drop-down [get]
func (r *UnitOfMeasurementControllerImpl) GetAllUnitOfMeasurementIsActive(writer http.ResponseWriter, request *http.Request) {

	result, err := r.unitofmeasurementservice.GetAllUnitOfMeasurementIsActive()

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Unit Of Measurement By Code
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Param uom_code path string true "uom_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/unit-of-measurement/by-code/{uom_code} [get]
func (r *UnitOfMeasurementControllerImpl) GetUnitOfMeasurementByCode(writer http.ResponseWriter, request *http.Request) {

	uomCode := chi.URLParam(request, "uom_code")
	result, err := r.unitofmeasurementservice.GetUnitOfMeasurementByCode(uomCode)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Unit Of Measurement
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @param reqBody body masteritempayloads.UomResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/unit-of-measurement/ [post]
func (r *UnitOfMeasurementControllerImpl) SaveUnitOfMeasurement(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.UomResponse
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message = ""

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.unitofmeasurementservice.SaveUnitOfMeasurement(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/unit-of-measurement/{uom_id} [patch]
func (r *UnitOfMeasurementControllerImpl) ChangeStatusUnitOfMeasurement(writer http.ResponseWriter, request *http.Request) {

	uomId, _ := strconv.Atoi(chi.URLParam(request, "uom_id"))

	response, err := r.unitofmeasurementservice.ChangeStatusUnitOfMeasurement(int(uomId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
