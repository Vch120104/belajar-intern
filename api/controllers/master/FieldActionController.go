package mastercontroller

import (
	// "after-sales/api/helper"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"strconv"

	// masterpayloads "after-sales/api/payloads/master"
	helper_test "after-sales/api/helper_testt"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	"net/http"

	"github.com/go-chi/chi/v5"
	// "strconv"
)

type FieldActionController interface {
	GetAllFieldAction(writer http.ResponseWriter, request *http.Request)
	GetFieldActionHeaderById(writer http.ResponseWriter, request *http.Request)
	GetAllFieldActionVehicleDetailById(writer http.ResponseWriter, request *http.Request)
	GetFieldActionVehicleDetailById(writer http.ResponseWriter, request *http.Request)
	GetAllFieldActionVehicleItemDetailById(writer http.ResponseWriter, request *http.Request)
	GetFieldActionVehicleItemDetailById(writer http.ResponseWriter, request *http.Request)
	PostFieldActionVehicleItemDetail(writer http.ResponseWriter, request *http.Request)
	PostFieldActionVehicleDetail(writer http.ResponseWriter, request *http.Request)
	PostMultipleVehicleDetail(writer http.ResponseWriter, request *http.Request)
	PostVehicleItemIntoAllVehicleDetail(writer http.ResponseWriter, request *http.Request)
	ChangeStatusFieldAction(writer http.ResponseWriter, request *http.Request)
	ChangeStatusFieldActionVehicle(writer http.ResponseWriter, request *http.Request)
	ChangeStatusFieldActionVehicleItem(writer http.ResponseWriter, request *http.Request)

	SaveFieldAction(writer http.ResponseWriter, request *http.Request)
}
type FieldActionControllerImpl struct {
	FieldActionService masterservice.FieldActionService
}

func NewFieldActionController(FieldActionService masterservice.FieldActionService) FieldActionController {
	return &FieldActionControllerImpl{
		FieldActionService: FieldActionService,
	}
}

// @Summary Get All Field Action
// @Description REST API Field Action
// @Accept json
// @Produce json
// @Tags Master : Field Action
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param field_action_system_number query string false "field_action_system_number"
// @Param field_action_document_number query string false "field_action_document_number"
// @Param approval_status_id query string false "approval_status_id"
// @Param brand_id query string false "brand_id"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/field-action [get]
func (r *FieldActionControllerImpl) GetAllFieldAction(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"field_action_system_number":   queryValues.Get("field_action_system_number"),
		"field_action_document_number": queryValues.Get("field_action_document_number"),
		"brand_id":                     queryValues.Get("brand_id"),
		"approval_status_id":           queryValues.Get("approval_status_id"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, totalPages, totalRows, err := r.FieldActionService.GetAllFieldAction(filterCondition, pagination)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", 200, pagination.Limit, pagination.Page, int64(totalRows), totalPages)
}

// @Summary Save Field Action
// @Description REST API Field Action
// @Accept json
// @Produce json
// @Tags Master : Field Action
// @Param reqBody body masterpayloads.FieldActionResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/field-action [post]
func (r *FieldActionControllerImpl) SaveFieldAction(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.FieldActionRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.FieldActionService.SaveFieldAction(formRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.FieldActionSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Get Field Action By Id
// @Description REST API Field Action
// @Accept json
// @Produce json
// @Tags Master : Field Action
// @Param field_action_system_number path int true "field_action_system_number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/field-action/by-id/{field_action_system_number} [get]
func (r *FieldActionControllerImpl) GetFieldActionHeaderById(writer http.ResponseWriter, request *http.Request) {
	FieldActionId, _ := strconv.Atoi(chi.URLParam(request, "field_action_system_number"))

	result, err := r.FieldActionService.GetFieldActionHeaderById(FieldActionId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Field Action Vehicle Detail By Id
// @Description REST API Field Action Vehicle Detail
// @Accept json
// @Produce json
// @Tags Master : Field Action
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_id query string false "item_id"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/field-action/vehicle/by-id/{field_action_system_number} [get]
func (r *FieldActionControllerImpl) GetAllFieldActionVehicleDetailById(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	FieldActionId, _ := strconv.Atoi(chi.URLParam(request, "field_action_system_number"))

	queryParams := map[string]string{
		"vehicle_id": queryValues.Get("vehicle_id"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.FieldActionService.GetAllFieldActionVehicleDetailById(FieldActionId, pagination, filterCondition)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Field Action Vehicle Detail By Id
// @Description REST API Field Action Vehicle Detail
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_detail_id path int true "item_substitute_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-substitute/detail/by-id/{item_substitute_id} [get]
func (r *FieldActionControllerImpl) GetFieldActionVehicleDetailById(writer http.ResponseWriter, request *http.Request) {
	FieldActionVehicleDetailId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_system_number"))

	result, err := r.FieldActionService.GetFieldActionVehicleDetailById(FieldActionVehicleDetailId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) GetAllFieldActionVehicleItemDetailById(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	FieldActionVehicleDetailId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_system_number"))

	// queryParams := map[string]string{
	// 	"vehicle_id": queryValues.Get("vehicle_id"),
	// }
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.FieldActionService.GetAllFieldActionVehicleItemDetailById(FieldActionVehicleDetailId, pagination)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *FieldActionControllerImpl) GetFieldActionVehicleItemDetailById(writer http.ResponseWriter, request *http.Request) {
	FieldActionVehicleItemDetailId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_item_system_number"))

	result, err := r.FieldActionService.GetFieldActionVehicleItemDetailById(FieldActionVehicleItemDetailId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) PostFieldActionVehicleItemDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.FieldActionItemDetailResponse
	FIeldActionVehicleDetailId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_system_number"))
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.FieldActionService.PostFieldActionVehicleItemDetail(FIeldActionVehicleDetailId, formRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.FieldActionEligibleVehicleItemSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) PostFieldActionVehicleDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.FieldActionDetailResponse
	FieldActionId, _ := strconv.Atoi(chi.URLParam(request, "field_action_system_number"))
	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.FieldActionService.PostFieldActionVehicleDetail(FieldActionId, formRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.FieldActionEligibleVehicleSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) PostMultipleVehicleDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	var formRequest masterpayloads.FieldActionDetailResponse

	FieldActionId, _ := strconv.Atoi(chi.URLParam(request, "field_action_system_number"))

	// CompanyIdStr := params.ByName("company_id")
	// CompanyId, _ := strconv.Atoi(CompanyIdStr)
	queryId := queryValues.Get("multi_id")

	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.FieldActionService.PostMultipleVehicleDetail(FieldActionId, queryId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.FieldActionEligibleVehicleSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) PostVehicleItemIntoAllVehicleDetail(writer http.ResponseWriter, request *http.Request) {
	// queryValues := request.URL.Query()
	var formRequest masterpayloads.FieldActionItemDetailResponse

	FieldActionHeaderId, _ := strconv.Atoi(chi.URLParam(request, "field_action_system_number"))

	helper.ReadFromRequestBody(request, &formRequest)
	var message string

	create, err := r.FieldActionService.PostVehicleItemIntoAllVehicleDetail(FieldActionHeaderId, formRequest)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.FieldActionEligibleVehicleItemSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) ChangeStatusFieldAction(writer http.ResponseWriter, request *http.Request) {

	FieldActionId, _ := strconv.Atoi(chi.URLParam(request, "field_action_system_number"))

	response, err := r.FieldActionService.ChangeStatusFieldAction(FieldActionId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) ChangeStatusFieldActionVehicle(writer http.ResponseWriter, request *http.Request) {

	FieldActionVehicleId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_system_number"))

	response, err := r.FieldActionService.ChangeStatusFieldActionVehicle(FieldActionVehicleId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) ChangeStatusFieldActionVehicleItem(writer http.ResponseWriter, request *http.Request) {
	FieldActionVehicleItemId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_item_system_number"))

	response, err := r.FieldActionService.ChangeStatusFieldActionVehicleItem(FieldActionVehicleItemId)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
