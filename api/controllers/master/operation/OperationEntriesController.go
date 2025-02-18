package masteroperationcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"strconv"

	// "after-sales/api/middlewares"

	masteroperationservice "after-sales/api/services/master/operation"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type OperationEntriesController interface {
	GetAllOperationEntries(writer http.ResponseWriter, request *http.Request)
	GetOperationEntriesByID(writer http.ResponseWriter, request *http.Request)
	GetOperationEntriesName(writer http.ResponseWriter, request *http.Request)
	SaveOperationEntries(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOperationEntries(writer http.ResponseWriter, request *http.Request)
}

type OperationEntriesControllerImpl struct {
	operationEntriesService masteroperationservice.OperationEntriesService
}

func NewOperationEntriesController(operationEntriesService masteroperationservice.OperationEntriesService) OperationEntriesController {
	return &OperationEntriesControllerImpl{
		operationEntriesService: operationEntriesService,
	}
}

// @Summary Get All Operation Entries
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_section_description query string false "operation_section_description"
// @Param operation_group_description query string false "operation_group_description"
// @Param operation_key_description query string false "operation_key_description"
// @Param operation_entries_code query string false "operation_entries_code"
// @Param operation_entries_desc query string false "operation_entries_desc"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-entries [get]
func (r *OperationEntriesControllerImpl) GetAllOperationEntries(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	queryParams := map[string]string{
		"mtr_operation_group.operation_group_description":     query.Get("operation_group_description"),
		"mtr_operation_section.operation_section_description": query.Get("operation_section_description"),
		"mtr_operation_entries.operation_entries_code":        query.Get("operation_entries_code"),
		"mtr_operation_entries.operation_entries_description": query.Get("operation_entries_description"),
		"mtr_operation_key.is_active":                         query.Get("is_active"),
		"mtr_operation_key.operation_key_description":         query.Get("operation_key_description"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(query, "limit"),
		Page:   utils.NewGetQueryInt(query, "page"),
		SortOf: query.Get("sort_of"),
		SortBy: query.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.operationEntriesService.GetAllOperationEntries(criteria, pagination)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if result.TotalPages == 0 {
		result.Rows = []masteroperationpayloads.OperationSectionListResponse{}
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Operation Entries By ID
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @Param operation_entries_id path int true "operation_entries_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-entries/{operation_entries_id} [get]
func (r *OperationEntriesControllerImpl) GetOperationEntriesByID(writer http.ResponseWriter, request *http.Request) {

	operationEntriesId, errA := strconv.Atoi(chi.URLParam(request, "operation_entries_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	result, err := r.operationEntriesService.GetOperationEntriesById(int(operationEntriesId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Entries Name
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @Param operation_group_id query int true "operation_group_id"
// @Param operation_section_id query int true "operation_section_id"
// @Param operation_key_id query int true "operation_key_id"
// @Param operation_entries_code query string true "operation_entries_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-entries/by-name [get]
func (r *OperationEntriesControllerImpl) GetOperationEntriesName(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()

	operationGroupId := utils.NewGetQueryInt(query, "operation_group_id")
	operationSectionId := utils.NewGetQueryInt(query, "operation_section_id")
	operationKeyId := utils.NewGetQueryInt(query, "operation_key_id")
	operationEntriesCode := query.Get("operation_entries_code")

	result, err := r.operationEntriesService.GetOperationEntriesName(masteroperationpayloads.OperationEntriesRequest{
		OperationGroupId:     operationGroupId,
		OperationSectionId:   operationSectionId,
		OperationKeyId:       operationKeyId,
		OperationEntriesCode: operationEntriesCode,
	})

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Entries
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @param reqBody body masteroperationpayloads.OperationEntriesResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-entries [post]
func (r *OperationEntriesControllerImpl) SaveOperationEntries(writer http.ResponseWriter, request *http.Request) {

	var requestForm masteroperationpayloads.OperationEntriesResponse
	var message = ""
	helper.ReadFromRequestBody(request, &requestForm)
	if validationErr := validation.ValidationForm(writer, request, &requestForm); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	create, err := r.operationEntriesService.SaveOperationEntries(requestForm)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	if requestForm.OperationEntriesId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Operation Entries
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @param operation_entries_id path int true "operation_entries_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-entries/{operation_entries_id} [patch]
func (r *OperationEntriesControllerImpl) ChangeStatusOperationEntries(writer http.ResponseWriter, request *http.Request) {
	operationEntriesId, errA := strconv.Atoi(chi.URLParam(request, "operation_entries_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.operationEntriesService.ChangeStatusOperationEntries(int(operationEntriesId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
