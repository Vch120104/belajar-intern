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
	"errors"
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
	UpdateDeductionDetail(writer http.ResponseWriter, request *http.Request)
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
		"deduction_id":   queryValues.Get("deduction_id"),
		"is_active":      queryValues.Get("is_active"),
		"deduction_name": queryValues.Get("deduction_name"),
		"deduction_code": queryValues.Get("deduction_code"),
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
	DeductionDetailIdstr, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.DeductionService.GetByIdDeductionDetail(DeductionDetailIdstr)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Deduction By Id
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @Param id path int true "id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/by-header-id/{id} [get]
func (r *DeductionControllerImpl) GetDeductionById(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	DeductionListId, errA := strconv.Atoi(chi.URLParam(request, "id"))
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.DeductionService.GetDeductionById(DeductionListId, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Deduction Detail
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @Param deduction_id path int true "deduction_id"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/{deduction_id} [get]
func (r *DeductionControllerImpl) GetAllDeductionDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	DeductionDetailId, _ := strconv.Atoi(chi.URLParam(request, "deduction_id"))

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
	if DeductionRequest.DeductionId == 0 {
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
	DeductionId, _ := strconv.Atoi(chi.URLParam(request, "deduction_id"))
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
	create, err := r.DeductionService.PostDeductionDetail(DeductionDetailRequest, DeductionId)
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
	DeductionId, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.DeductionService.ChangeStatusDeduction(DeductionId)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Update Deduction Detail
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @param id path int true "id"
// @param reqBody body masterpayloads.DeductionDetailUpdate true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/deduction/detail/{id} [put]
func (r *DeductionControllerImpl) UpdateDeductionDetail(writer http.ResponseWriter, request *http.Request) {
	DeductionDetailRequest := masterpayloads.DeductionDetailUpdate{}
	DeductionId, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
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
	res, err := r.DeductionService.UpdateDeductionDetail(DeductionId, DeductionDetailRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Update Data Successfully!", http.StatusOK)
}
