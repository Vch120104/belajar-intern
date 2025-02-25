package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	masteritemservice "after-sales/api/services/master/item"

	// "after-sales/api/middlewares"

	// "strconv"

	"github.com/go-chi/chi/v5"
)

type MarkupRateController interface {
	GetAllMarkupRate(writer http.ResponseWriter, request *http.Request)
	GetMarkupRateByID(writer http.ResponseWriter, request *http.Request)
	SaveMarkupRate(writer http.ResponseWriter, request *http.Request)
	ChangeStatusMarkupRate(writer http.ResponseWriter, request *http.Request)
	GetMarkupRateByMarkupMasterAndOrderType(writer http.ResponseWriter, request *http.Request)
}
type MarkupRateControllerImpl struct {
	MarkupRateService masteritemservice.MarkupRateService
}

func NewMarkupRateController(markupRateService masteritemservice.MarkupRateService) MarkupRateController {
	return &MarkupRateControllerImpl{
		MarkupRateService: markupRateService,
	}
}

// @Summary Get All Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Security AuthorizationKeyAuth
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param markup_master_code query string false "markup_master_code"
// @Param markup_master_description query string false "markup_master_description"
// @Param order_type_name query string false "order_type_name"
// @Param markup_rate query float64 false "markup_rate"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-rate/ [get]
func (r *MarkupRateControllerImpl) GetAllMarkupRate(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_markup_master.markup_code":        queryValues.Get("markup_code"),
		"mtr_markup_master.markup_description": queryValues.Get("markup_description"),
		"order_type_name":                      queryValues.Get("order_type_name"),
		"mtr_markup_rate.markup_rate":          queryValues.Get("markup_rate"),
		"mtr_markup_rate.is_active":            queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.MarkupRateService.GetAllMarkupRate(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

// @Summary Get Markup Rate By ID
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Security AuthorizationKeyAuth
// @Param markup_rate_id path int true "markup_rate_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-rate/{markup_rate_id} [get]
func (r *MarkupRateControllerImpl) GetMarkupRateByID(writer http.ResponseWriter, request *http.Request) {

	markupRateId, errA := strconv.Atoi(chi.URLParam(request, "markup_rate_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.MarkupRateService.GetMarkupRateById(markupRateId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Security AuthorizationKeyAuth
// @param reqBody body masteritempayloads.MarkupRateRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-rate/ [post]
func (r *MarkupRateControllerImpl) SaveMarkupRate(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.MarkupRateRequest
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

	create, err := r.MarkupRateService.SaveMarkupRate(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.MarkupRateId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Security AuthorizationKeyAuth
// @param markup_rate_id path int true "markup_rate_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-rate/{markup_rate_id} [patch]
func (r *MarkupRateControllerImpl) ChangeStatusMarkupRate(writer http.ResponseWriter, request *http.Request) {

	markupRateId, errA := strconv.Atoi(chi.URLParam(request, "markup_rate_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.MarkupRateService.ChangeStatusMarkupRate(int(markupRateId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get Markup Rate By Markup Master And Order Type
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Security AuthorizationKeyAuth
// @Param markup_master_id path int true "markup_master_id"
// @Param order_type_id path int true "order_type_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/markup-rate/markup-master/{markup_master_id}/order-type/{order_type_id} [get]
func (r *MarkupRateControllerImpl) GetMarkupRateByMarkupMasterAndOrderType(writer http.ResponseWriter, request *http.Request) {

	markupMasterId, errA := strconv.Atoi(chi.URLParam(request, "markup_master_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	orderTypeId, errA := strconv.Atoi(chi.URLParam(request, "order_type_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.MarkupRateService.GetMarkupRateByMarkupMasterAndOrderType(markupMasterId, orderTypeId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
