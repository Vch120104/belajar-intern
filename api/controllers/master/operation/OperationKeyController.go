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
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OperationKeyController interface {
	GetAllOperationKeyList(writer http.ResponseWriter, request *http.Request)
	GetOperationKeyByID(writer http.ResponseWriter, request *http.Request)
	GetOperationKeyName(writer http.ResponseWriter, request *http.Request)
	SaveOperationKey(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOperationKey(writer http.ResponseWriter, request *http.Request)
	GetOperationKeyDropdown(writer http.ResponseWriter, request *http.Request)
}

type OperationKeyControllerImpl struct {
	operationkeyservice masteroperationservice.OperationKeyService
}

func NewOperationKeyController(operationKeyService masteroperationservice.OperationKeyService) OperationKeyController {
	return &OperationKeyControllerImpl{
		operationkeyservice: operationKeyService,
	}
}

// @Summary Get All Operation Key
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Security AuthorizationKeyAuth
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_section_code query string false "operation_section_code"
// @Param operation_section_description query string false "operation_section_description"
// @Param operation_group_code query string false "operation_group_code"
// @Param operation_group_description query string false "operation_group_description"
// @Param operation_key_code query string false "operation_key_code"
// @Param operation_key_description query string false "operation_key_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-key/ [get]
func (r *OperationKeyControllerImpl) GetAllOperationKeyList(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	queryParams := map[string]string{
		"mtr_operation_group.operation_group_code":            query.Get("operation_group_code"),
		"mtr_operation_group.operation_group_description":     query.Get("operation_group_description"),
		"mtr_operation_section.operation_section_code":        query.Get("operation_section_code"),
		"mtr_operation_section.operation_section_description": query.Get("operation_section_description"),
		"mtr_operation_key.is_active":                         query.Get("is_active"),
		"mtr_operation_key.operation_key_code":                query.Get("operation_key_code"),
		"mtr_operation_key.operation_key_description":         query.Get("operation_key_description"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(query, "limit"),
		Page:   utils.NewGetQueryInt(query, "page"),
		SortOf: query.Get("sort_of"),
		SortBy: query.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.operationkeyservice.GetAllOperationKeyList(criteria, pagination)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if result.TotalPages == 0 {
		result.Rows = []masteroperationpayloads.OperationkeyListResponse{}
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Operation Key By ID
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Security AuthorizationKeyAuth
// @Param operation_key_id path int true "operation_key_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-key/{operation_key_id} [get]
func (r *OperationKeyControllerImpl) GetOperationKeyByID(writer http.ResponseWriter, request *http.Request) {
	operationKeyId, errA := strconv.Atoi(chi.URLParam(request, "operation_key_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	result, err := r.operationkeyservice.GetOperationKeyById(operationKeyId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Key Name
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Security AuthorizationKeyAuth
// @Param operation_group_id query int true "operation_group_id"
// @Param operation_section_id query int true "operation_section_id"
// @Param operation_key_code query string true "operation_key_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-key/name [get]
func (r *OperationKeyControllerImpl) GetOperationKeyName(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	operationGroupId := utils.NewGetQueryInt(query, "operation_group_id")
	operationSectionId := utils.NewGetQueryInt(query, "operation_section_id")
	keyCode := query.Get("operation_key_code")

	result, err := r.operationkeyservice.GetOperationKeyName(masteroperationpayloads.OperationKeyRequest{
		OperationGroupId:   int(operationGroupId),
		OperationSectionId: int(operationSectionId),
		OperationKeyCode:   keyCode,
	})

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Key
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Security AuthorizationKeyAuth
// @param reqBody body masteroperationpayloads.OperationKeyResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-key/ [post]
func (r *OperationKeyControllerImpl) SaveOperationKey(writer http.ResponseWriter, request *http.Request) {
	var requestForm masteroperationpayloads.OperationKeyResponse
	var message = ""

	err := jsonchecker.ReadFromRequestBody(request, &requestForm)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, requestForm)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.operationkeyservice.SaveOperationKey(requestForm)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if requestForm.OperationKeyId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Operation Key
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Security AuthorizationKeyAuth
// @param operation_key_id path int true "operation_key_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-key/{operation_key_id} [patch]
func (r *OperationKeyControllerImpl) ChangeStatusOperationKey(writer http.ResponseWriter, request *http.Request) {
	operationKeyId, errA := strconv.Atoi(chi.URLParam(request, "operation_key_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.operationkeyservice.ChangeStatusOperationKey(int(operationKeyId))

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Key Dropdown
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Security AuthorizationKeyAuth
// @Param operation_group_id path int true "operation_group_id"
// @Param operation_section_id path int true "operation_section_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-key/drop-down/{operation_group_id}/{operation_section_id} [get]
func (r *OperationKeyControllerImpl) GetOperationKeyDropdown(writer http.ResponseWriter, request *http.Request) {

	operationGroupId, err := strconv.Atoi(chi.URLParam(request, "operation_group_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Operation Group ID", http.StatusBadRequest)
		return
	}

	operationSectionId, err := strconv.Atoi(chi.URLParam(request, "operation_section_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Operation Section ID", http.StatusBadRequest)
		return
	}

	get, errResp := r.operationkeyservice.GetOperationKeyDropdown(operationGroupId, operationSectionId)
	if errResp != nil {
		helper.ReturnError(writer, request, errResp)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}
