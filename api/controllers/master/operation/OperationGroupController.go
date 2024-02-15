package masteroperationcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type OperationGroupController interface {
	GetAllOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllOperationGroupIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetOperationGroupByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-group [get]
func (r *OperationGroupControllerImpl) GetAllOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

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

	result := r.OperationGroupService.GetAllOperationGroup(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Operation Group drop down
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-group/drop-down [get]
func (r *OperationGroupControllerImpl) GetAllOperationGroupIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	result := r.OperationGroupService.GetAllOperationGroupIsActive()

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Group By Code
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @Param operation_group_code path string true "operation_group_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-group/by-code/{operation_group_code} [get]
func (r *OperationGroupControllerImpl) GetOperationGroupByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	operationGroupCode := params.ByName("operation_group_code")

	result := r.OperationGroupService.GetOperationGroupByCode(operationGroupCode)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Group
// @Description REST API Operation Group
// @Accept json
// @Produce json
// @Tags Master : Operation Group
// @param reqBody body masteroperationpayloads.OperationGroupResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-group [post]
func (r *OperationGroupControllerImpl) SaveOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masteroperationpayloads.OperationGroupResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.OperationGroupService.SaveOperationGroup(formRequest)

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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-group/{operation_group_id} [patch]
func (r *OperationGroupControllerImpl) ChangeStatusOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	operationGroupId, _ := strconv.Atoi(params.ByName("operation_group_id"))

	response := r.OperationGroupService.ChangeStatusOperationGroup(int(operationGroupId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
