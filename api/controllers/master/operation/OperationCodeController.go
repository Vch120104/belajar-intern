package masteroperationcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OperationCodeController interface {
	GetAllOperationCode(writer http.ResponseWriter, request *http.Request)
	GetByIdOperationCode(writer http.ResponseWriter, request *http.Request)
	GetByCodeOperationCode(writer http.ResponseWriter, request *http.Request)
	SaveOperationCode(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOperationCode(writer http.ResponseWriter, request *http.Request)
	UpdateOperationCode(writer http.ResponseWriter, request *http.Request)
}

type OperationCodeControllerImpl struct {
	operationCodeService masteroperationservice.OperationCodeService
}

func NewOperationCodeController(operationCodeservice masteroperationservice.OperationCodeService) OperationCodeController {
	return &OperationCodeControllerImpl{
		operationCodeService: operationCodeservice,
	}
}

// @Summary Get All OPeration Code
// @Description REST API Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_code query string false "operation_code"
// @Param operation_name query string false "operation_name"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-code/ [get]
func (r *OperationCodeControllerImpl) GetAllOperationCode(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":      queryValues.Get("is_active"),
		"operation_code": queryValues.Get("operation_code"),
		"operation_name": queryValues.Get("operation_name"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.operationCodeService.GetAllOperationCode(filterCondition, pagination)

	if err != nil {
		payloads.NewHandleSuccessPagination(writer, []interface{}{}, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Operation Code By Id
// @Description REST API Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @Param operation_id path int true "operation_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-code/by-id/{operation_id} [get]
func (r *OperationCodeControllerImpl) GetByIdOperationCode(writer http.ResponseWriter, request *http.Request) {
	OperationIdStr, errA := strconv.Atoi(chi.URLParam(request, "operation_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.operationCodeService.GetOperationCodeById(int(OperationIdStr))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

func (r *OperationCodeControllerImpl) GetByCodeOperationCode(writer http.ResponseWriter, request *http.Request) {
	OperationCodeStr := chi.URLParam(request, "operation_code")

	result, err := r.operationCodeService.GetOperationCodeByCode(OperationCodeStr)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Code
// @Description REST API Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @param reqBody body masteroperationpayloads.OperationCodeSave true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-code/ [post]
func (r *OperationCodeControllerImpl) SaveOperationCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationCodeSave
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.operationCodeService.SaveOperationCode(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// @Summary Change Status Patch Operation Code
// @Description REST API Patch Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @param operation_id path int true "operation_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-code/{operation_id} [patch]
func (r *OperationCodeControllerImpl) ChangeStatusOperationCode(writer http.ResponseWriter, request *http.Request) {

	OperationId, errA := strconv.Atoi(chi.URLParam(request, "operation_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.operationCodeService.ChangeStatusOperationCode(OperationId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *OperationCodeControllerImpl) UpdateOperationCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationCodeUpdate

	OperationCodeId, errA := strconv.Atoi(chi.URLParam(request, "operation_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	update, err := r.operationCodeService.UpdateItemCode(OperationCodeId, formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Create Data Successfully!", http.StatusOK)
}
