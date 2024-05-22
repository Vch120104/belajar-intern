package masteroperationcontroller

import (
	exceptions "after-sales/api/exceptions"
	helper_test "after-sales/api/helper_testt"
	"after-sales/api/validation"

	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OperationGroupController interface {
	GetAllOperationGroup(writer http.ResponseWriter, request *http.Request)
	GetAllOperationGroupIsActive(writer http.ResponseWriter, request *http.Request)
	GetOperationGroupByCode(writer http.ResponseWriter, request *http.Request)
	SaveOperationGroup(writer http.ResponseWriter, request *http.Request)
	ChangeStatusOperationGroup(writer http.ResponseWriter, request *http.Request)
}
type OperationGroupControllerImpl struct {
	OperationGroupService masteroperationservice.OperationGroupService
}

func NewOperationGroupController(operationGroupService masteroperationservice.OperationGroupService) OperationGroupController {
	return &OperationGroupControllerImpl{
		OperationGroupService: operationGroupService,
	}
}

// @Summary Get All Operation Group
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_group_code query string false "operation_group_code"
// @Param operation_group_description query string false "operation_group_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-group/ [get]
func (r *OperationGroupControllerImpl) GetAllOperationGroup(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"operation_group_code":        queryValues.Get("operation_group_code"),
		"operation_group_description": queryValues.Get("operation_group_description"),
		"is_active":                   queryValues.Get("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.OperationGroupService.GetAllOperationGroup(filterCondition, pagination)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Operation Group drop down
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-group/drop-down [get]
func (r *OperationGroupControllerImpl) GetAllOperationGroupIsActive(writer http.ResponseWriter, request *http.Request) {

	result, err := r.OperationGroupService.GetAllOperationGroupIsActive()
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Group By Code
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @Param operation_group_code path string true "operation_group_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-group/by-code/{operation_group_code} [get]
func (r *OperationGroupControllerImpl) GetOperationGroupByCode(writer http.ResponseWriter, request *http.Request) {

	operationGroupCode := chi.URLParam(request, "operation_group_code")

	result, err := r.OperationGroupService.GetOperationGroupByCode(operationGroupCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Group
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @param reqBody body masteroperationpayloads.OperationGroupResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-group/ [post]
func (r *OperationGroupControllerImpl) SaveOperationGroup(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteroperationpayloads.OperationGroupResponse
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

	create, err := r.OperationGroupService.SaveOperationGroup(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.OperationGroupId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Operation Group
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @param operation_group_id path int true "operation_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/operation-group/{operation_group_id} [patch]
func (r *OperationGroupControllerImpl) ChangeStatusOperationGroup(writer http.ResponseWriter, request *http.Request) {

	operationGroupId, _ := strconv.Atoi(chi.URLParam(request, "operation_group_id"))

	response, err := r.OperationGroupService.ChangeStatusOperationGroup(int(operationGroupId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
