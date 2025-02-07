package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UnitOfMeasurementController interface {
	GetAllUnitOfMeasurement(writer http.ResponseWriter, request *http.Request)
	GetAllUnitOfMeasurementIsActive(writer http.ResponseWriter, request *http.Request)
	GetUnitOfMeasurementByCode(writer http.ResponseWriter, request *http.Request)
	GetUnitOfMeasurementById(writer http.ResponseWriter, request *http.Request)
	SaveUnitOfMeasurement(writer http.ResponseWriter, request *http.Request)
	ChangeStatusUnitOfMeasurement(writer http.ResponseWriter, request *http.Request)
	GetUnitOfMeasurementItem(writer http.ResponseWriter, request *http.Request)
	GetQuantityConversion(writer http.ResponseWriter, request *http.Request)
}

type UnitOfMeasurementControllerImpl struct {
	unitofmeasurementservice masteritemservice.UnitOfMeasurementService
}

func NewUnitOfMeasurementController(UnitOfMeasurementService masteritemservice.UnitOfMeasurementService) UnitOfMeasurementController {
	return &UnitOfMeasurementControllerImpl{
		unitofmeasurementservice: UnitOfMeasurementService,
	}
}

// @Summary Get Unit Of Measurement By Id
// @Description	REST API Unit Of Measurement
// @Accept json
// @Produce	json
// @Tags Master Item : Unit Of Measurement
// @Param uom_id path string true "uom_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/unit-of-measurement/{uom_id} [get]
func (r *UnitOfMeasurementControllerImpl) GetUnitOfMeasurementById(writer http.ResponseWriter, request *http.Request) {
	uomId, errA := strconv.Atoi(chi.URLParam(request, "uom_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.unitofmeasurementservice.GetUnitOfMeasurementById(int(uomId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// @Summary		Get All Unit Of Measurement
// @Description	REST API Unit Of Measurement
// @Accept			json
// @Produce		json
// @Tags Master Item : Unit Of Measurement
// @Param			page					query		string	true	"page"
// @Param			limit					query		string	true	"limit"
// @Param			is_active				query		string	false	"is_active"	Enums(true, false)
// @Param			uom_code				query		string	false	"uom_code"
// @Param			uom_description			query		string	false	"uom_description"
// @Param			uom_type_desc			query		string	false	"uom_type_desc"
// @Param			sort_by					query		string	false	"sort_by"
// @Param			sort_of					query		string	false	"sort_of"
// @Success		200						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/unit-of-measurement/ [get]
func (r *UnitOfMeasurementControllerImpl) GetAllUnitOfMeasurement(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_uom.is_active":                 queryValues.Get("is_active"),
		"mtr_uom.uom_code":                  queryValues.Get("uom_code"),
		"mtr_uom.uom_description":           queryValues.Get("uom_description"),
		"mtr_uom_type.uom_type_description": queryValues.Get("uom_type_description"),
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
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary		Get All Unit Of Measurement drop down
// @Description	REST API Unit Of Measurement
// @Accept			json
// @Produce		json
// @Tags Master Item : Unit Of Measurement
// @Success		200						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/unit-of-measurement/drop-down [get]
func (r *UnitOfMeasurementControllerImpl) GetAllUnitOfMeasurementIsActive(writer http.ResponseWriter, request *http.Request) {

	result, err := r.unitofmeasurementservice.GetAllUnitOfMeasurementIsActive()

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary		Get Unit Of Measurement By Code
// @Description	REST API Unit Of Measurement
// @Accept			json
// @Produce		json
// @Tags Master Item : Unit Of Measurement
// @Param			uom_code				path		string	true	"uom_code"
// @Success		200						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/unit-of-measurement/by-code/{uom_code} [get]
func (r *UnitOfMeasurementControllerImpl) GetUnitOfMeasurementByCode(writer http.ResponseWriter, request *http.Request) {

	uomCode := chi.URLParam(request, "uom_code")
	result, err := r.unitofmeasurementservice.GetUnitOfMeasurementByCode(uomCode)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary		Save Unit Of Measurement
// @Description	REST API Unit Of Measurement
// @Accept			json
// @Produce		json
// @Tags Master Item : Unit Of Measurement
// @param			reqBody					body		masteritempayloads.UomResponse	true	"Form Request"
// @Success		200						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/unit-of-measurement/ [post]
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
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.UomId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary		Change Status Unit Of Measurement
// @Description	REST API Unit Of Measurement
// @Accept			json
// @Produce		json
// @Tags Master Item : Unit Of Measurement
// @param			uom_id					path		int	true	"uom_id"
// @Success		200						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/unit-of-measurement/{uom_id} [patch]
func (r *UnitOfMeasurementControllerImpl) ChangeStatusUnitOfMeasurement(writer http.ResponseWriter, request *http.Request) {

	uomId, errA := strconv.Atoi(chi.URLParam(request, "uom_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.unitofmeasurementservice.ChangeStatusUnitOfMeasurement(int(uomId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary		Get Unit Of Measurement Item By Item Id
// @Description	REST API Unit Of Measurement Item
// @Accept			json
// @Produce		json
// @Tags Master Item : Unit Of Measurement
// @Param			item_id					path		string	true	"item_id"
// @Param			source_type				path		string	true	"source_type"
// @Success		200						{object}	masteritempayloads.UomItemResponses
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/unit-of-measurement/{item_id}/{source_type} [get]
func (r *UnitOfMeasurementControllerImpl) GetUnitOfMeasurementItem(writer http.ResponseWriter, request *http.Request) {
	ItemId, errA := strconv.Atoi(chi.URLParam(request, "item_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	SourceUomType := chi.URLParam(request, "source_type")
	payload := masteritempayloads.UomItemRequest{
		SourceType: SourceUomType,
		ItemId:     ItemId,
	}
	response, err := r.unitofmeasurementservice.GetUnitOfMeasurementItem(payload)
	if err != nil {
		helper.ReturnError(writer, request, err)
	}
	payloads.NewHandleSuccess(writer, response, "Get Data Success!", http.StatusOK)

}

// @Summary		Get Quantity Conversion
// @Description	REST API Unit Of Measurement
// @Accept			json
// @Produce		json
// @Tags Master Item : Unit Of Measurement
// @Param			source_type				query		string	true	"source_type"
// @Param			item_id					query		int		true	"item_id"
// @Param			quantity				query		float64	true	"quantity"
// @Success		200						{object}	payloads.Response
// @Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
// @Router			/v1/unit-of-measurement/get_quantity_conversion [get]
func (r *UnitOfMeasurementControllerImpl) GetQuantityConversion(writer http.ResponseWriter, request *http.Request) {
	//var formRequest masteritempayloads.UomItemRequest
	//groupServiceUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_name=" + groupName
	//helper.ReadFromRequestBody(request, &formRequest)
	//SourceType string  `json:"source_type"`
	//ItemId     int     `json:"item_id"`
	//Quantity   float64 `json:"quantity"`
	//<<localhostp8000>>unit-of-measurement/get_quantity_conversion?source_type=S&item_id=893891&quantity=1.0
	queryValues := request.URL.Query()
	formRequest := masteritempayloads.UomGetQuantityConversion{
		Quantity:   utils.NewGetQueryfloat(queryValues, "quantity"),
		SourceType: queryValues.Get("source_type"),
		ItemId:     utils.NewGetQueryInt(queryValues, "item_id"),
	}
	//var GoodsReceiveHeaderPayloads masteritempayloads.UomItemRequest
	//helper.ReadFromRequestBody(request, &GoodsReceiveHeaderPayloads)

	//if err != nil {
	//	exceptions.NewEntityException(writer, request, err)
	//	return
	//}

	err := validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.unitofmeasurementservice.GetQuantityConversion(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "get quantity conversion", http.StatusOK)
}
