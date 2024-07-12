package masteroperationcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"

	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
)

type OperationSectionController interface {
	GetAllOperationSectionList(writer http.ResponseWriter, request *http.Request)
	GetOperationSectionByID(writer http.ResponseWriter, request *http.Request)
	GetSectionCodeByGroupId(writer http.ResponseWriter, request *http.Request)
	GetOperationSectionName(writer http.ResponseWriter, request *http.Request)
	SaveOperationSection(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOperationSection(writer http.ResponseWriter, request *http.Request)
}

type OperationSectionControllerImpl struct {
	operationsectionservice masteroperationservice.OperationSectionService
}

func NewOperationSectionController(operationSectionService masteroperationservice.OperationSectionService) OperationSectionController {
	return &OperationSectionControllerImpl{
		operationsectionservice: operationSectionService,
	}
}

// @Summary Get All Operation Section
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_section_code query string false "operation_section_code"
// @Param operation_section_description query string false "operation_section_description"
// @Param operation_group_code query string false "operation_group_code"
// @Param operation_group_description query string false "operation_group_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-section/ [get]
func (r *OperationSectionControllerImpl) GetAllOperationSectionList(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	queryParams := map[string]string{
		"mtr_operation_group.operation_group_code":            query.Get("operation_group_code"),
		"mtr_operation_group.operation_group_description":     query.Get("operation_group_description"),
		"mtr_operation_section.operation_section_code":        query.Get("operation_section_code"),
		"mtr_operation_section.operation_section_description": query.Get("operation_section_description"),
		"mtr_operation_section.is_active":                     query.Get("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(query, "limit"),
		Page:   utils.NewGetQueryInt(query, "page"),
		SortOf: query.Get("sort_of"),
		SortBy: query.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.operationsectionservice.GetAllOperationSectionList(filterCondition, pagination)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Operation Section By ID
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param operation_section_id path int true "operation_section_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-section/{operation_section_id} [get]
func (r *OperationSectionControllerImpl) GetOperationSectionByID(writer http.ResponseWriter, request *http.Request) {

	operationSectionId, _ := strconv.Atoi(chi.URLParam(request, "operation_section_id"))

	result, err := r.operationsectionservice.GetOperationSectionById(int(operationSectionId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Section Code By Group Id
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param operation_group_id query int true "operation_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-section/code-by-group-id [get]
func (r *OperationSectionControllerImpl) GetSectionCodeByGroupId(writer http.ResponseWriter, request *http.Request) {

	groupId, _ := strconv.Atoi(chi.URLParam(request, "operation_group_id"))

	result, err := r.operationsectionservice.GetSectionCodeByGroupId(groupId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Section Name
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param operation_group_id query int true "operation_group_id"
// @Param operation_section_code query string true "operation_section_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-section/operation-section-name [get]
func (r *OperationSectionControllerImpl) GetOperationSectionName(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	operationGroupId := utils.NewGetQueryInt(query, "operation_group_id")
	section_code := query.Get("operation_section_code")

	result, err := r.operationsectionservice.GetOperationSectionName(operationGroupId, section_code)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Section
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @param reqBody body masteroperationpayloads.OperationSectionRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-section/ [put]
func (r *OperationSectionControllerImpl) SaveOperationSection(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteroperationpayloads.OperationSectionRequest
	helper.ReadFromRequestBody(request, &formRequest)

	var message = ""

	create, err := r.operationsectionservice.SaveOperationSection(formRequest)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if formRequest.OperationSectionId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Operation Section
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @param operation_section_id path int true "operation_section_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-section/{operation_section_id} [patch]
func (r *OperationSectionControllerImpl) ChangeStatusOperationSection(writer http.ResponseWriter, request *http.Request) {

	operationSectionId, _ := strconv.Atoi(chi.URLParam(request, "operation_section_id"))

	response, err := r.operationsectionservice.ChangeStatusOperationSection(int(operationSectionId))

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
