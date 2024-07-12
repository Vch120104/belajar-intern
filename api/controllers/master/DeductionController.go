package mastercontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DeductionController interface {
	GetAllDeductionList(writer http.ResponseWriter, request *http.Request)
	GetByIdDeductionDetail(writer http.ResponseWriter, request *http.Request)
	GetDeductionById(writer http.ResponseWriter, request *http.Request)
	GetAllDeductionDetail(writer http.ResponseWriter, request *http.Request)
	SaveDeductionList(writer http.ResponseWriter, request *http.Request)
	SaveDeductionDetail(writer http.ResponseWriter, request *http.Request)
	ChangeStatusDeduction(writer http.ResponseWriter, request *http.Request)
}

type DeductionControllerImpl struct {
	DeductionService masterservice.DeductionService
}

func NewDeductionController(deductionService masterservice.DeductionService) DeductionController {
	return &DeductionControllerImpl{
		DeductionService: deductionService,
	}
}

// @Summary Get All Deduction
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param deduction_code query string false "deduction_code"
// @Param deduction_name query string false "deduction_name"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/ [get]
func (r *DeductionControllerImpl) GetAllDeductionList(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"is_active":      queryValues.Get("is_active"),
		"deduction_name": queryValues.Get("deduction_name"),
		"effective_date": queryValues.Get("effective_date"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.DeductionService.GetAllDeduction(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Deduction Detail By Id
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @Param deduction_detail_id path int true "deduction_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/detail/by-id/{deduction_detail_id} [get]
func (r *DeductionControllerImpl) GetByIdDeductionDetail(writer http.ResponseWriter, request *http.Request) {
	DeductionDetailIdstr, _ := strconv.Atoi(chi.URLParam(request, "id"))

	result, err := r.DeductionService.GetByIdDeductionDetail(DeductionDetailIdstr)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *DeductionControllerImpl) GetDeductionById(writer http.ResponseWriter, request *http.Request) {
	DeductionListId, _ := strconv.Atoi(chi.URLParam(request, "id"))

	result, err := r.DeductionService.GetDeductionById(DeductionListId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *DeductionControllerImpl) GetAllDeductionDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	DeductionDetailId, _ := strconv.Atoi(chi.URLParam(request, "id"))

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.DeductionService.GetAllDeductionDetail(DeductionDetailId, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Save Deduction List
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @param reqBody body masterpayloads.DeductionListResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/ [post]
func (r *DeductionControllerImpl) SaveDeductionList(writer http.ResponseWriter, request *http.Request) {
	DeductionRequest := masterpayloads.DeductionListResponse{}
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &DeductionRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, DeductionRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	create, err := r.DeductionService.PostDeductionList(DeductionRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	if DeductionRequest.DeductionListId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}
	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}

// @Summary Save Deduction Detail
// @Description REST API Deduction Detail
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @param reqBody body masterpayloads.DeductionDetailResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/detail [post]
func (r *DeductionControllerImpl) SaveDeductionDetail(writer http.ResponseWriter, request *http.Request) {
	DeductionDetailRequest := masterpayloads.DeductionDetailResponse{}
	var message string

	err := jsonchecker.ReadFromRequestBody(request, &DeductionDetailRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, DeductionDetailRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	create, err := r.DeductionService.PostDeductionDetail(DeductionDetailRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	if DeductionDetailRequest.DeductionDetailId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}
	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}

// @Summary Change Status Deduction
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @param deduction_list_id path int true "deduction_list_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/{deduction_list_id} [patch]
func (r *DeductionControllerImpl) ChangeStatusDeduction(writer http.ResponseWriter, request *http.Request) {
	DeductionId, _ := strconv.Atoi(chi.URLParam(request, "id"))

	response, err := r.DeductionService.ChangeStatusDeduction(DeductionId)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
