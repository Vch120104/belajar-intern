package mastercontroller

import (
	// "after-sales/api/helper"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"strconv"

	// masterpayloads "after-sales/api/payloads/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	"net/http"
	// "strconv"
	"github.com/go-chi/chi/v5"
	//"github.com/julienschmidt/httprouter"
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
// @Param approval_value query string false "approval_value"
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
		"approval_value":               queryValues.Get("approval_value"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.FieldActionService.GetAllFieldAction(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
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

	var formRequest masterpayloads.FieldActionResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.FieldActionService.SaveFieldAction(formRequest)

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
	FieldActionIdStr := chi.URLParam(request, "field_action_system_number")

	FieldActionId, _ := strconv.Atoi(FieldActionIdStr)

	result := r.FieldActionService.GetFieldActionHeaderById(FieldActionId)

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

	FieldActionIdStr := chi.URLParam(request, "field_action_system_number")
	FieldActionId, _ := strconv.Atoi(FieldActionIdStr)

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

	result := r.FieldActionService.GetAllFieldActionVehicleDetailById(FieldActionId, pagination, filterCondition)

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
	FieldActionVehicleDetailIdStr := chi.URLParam(request, "field_action_eligible_vehicle_system_number")

	FieldActionVehicleDetailId, _ := strconv.Atoi(FieldActionVehicleDetailIdStr)

	result := r.FieldActionService.GetFieldActionVehicleDetailById(FieldActionVehicleDetailId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) GetAllFieldActionVehicleItemDetailById(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	FieldActionVehicleDetailIdStr := chi.URLParam(request, "field_action_eligible_vehicle_system_number")
	FieldActionVehicleDetailId, _ := strconv.Atoi(FieldActionVehicleDetailIdStr)

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

	result := r.FieldActionService.GetAllFieldActionVehicleItemDetailById(FieldActionVehicleDetailId, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *FieldActionControllerImpl) GetFieldActionVehicleItemDetailById(writer http.ResponseWriter, request *http.Request) {
	FieldActionVehicleItemDetailIdStr := chi.URLParam(request, "field_action_eligible_vehicle_item_system_number")

	FieldActionVehicleItemDetailId, _ := strconv.Atoi(FieldActionVehicleItemDetailIdStr)

	result := r.FieldActionService.GetFieldActionVehicleItemDetailById(FieldActionVehicleItemDetailId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) PostFieldActionVehicleItemDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.FieldActionItemDetailResponse
	FIeldActionVehicleDetailIdStr := chi.URLParam(request, "field_action_eligible_vehicle_system_number")

	FIeldActionVehicleDetailId, _ := strconv.Atoi(FIeldActionVehicleDetailIdStr)
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.FieldActionService.PostFieldActionVehicleItemDetail(FIeldActionVehicleDetailId, formRequest)

	if formRequest.FieldActionEligibleVehicleItemSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) PostFieldActionVehicleDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masterpayloads.FieldActionDetailResponse
	FieldActionIdStr := chi.URLParam(request, "field_action_system_number")

	FieldActionId, _ := strconv.Atoi(FieldActionIdStr)
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.FieldActionService.PostFieldActionVehicleDetail(FieldActionId, formRequest)

	if formRequest.FieldActionSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) PostMultipleVehicleDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	var formRequest masterpayloads.FieldActionDetailResponse

	FieldActionIdStr := chi.URLParam(request, "field_action_system_number")
	FieldActionId, _ := strconv.Atoi(FieldActionIdStr)

	// CompanyIdStr := params.ByName("company_id")
	// CompanyId, _ := strconv.Atoi(CompanyIdStr)
	queryId := queryValues.Get("multi_id")

	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.FieldActionService.PostMultipleVehicleDetail(FieldActionId, queryId)

	if formRequest.FieldActionSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) PostVehicleItemIntoAllVehicleDetail(writer http.ResponseWriter, request *http.Request) {
	// queryValues := request.URL.Query()
	var formRequest masterpayloads.FieldActionItemDetailResponse

	FieldActionHeaderIdStr := chi.URLParam(request, "field_action_system_number")
	FieldActionHeaderId, _ := strconv.Atoi(FieldActionHeaderIdStr)

	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.FieldActionService.PostVehicleItemIntoAllVehicleDetail(FieldActionHeaderId, formRequest)

	if formRequest.FieldActionEligibleVehicleItemSystemNumber == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *FieldActionControllerImpl) ChangeStatusFieldAction(writer http.ResponseWriter, request *http.Request) {

	FieldActionId, _ := strconv.Atoi(chi.URLParam(request, "field_action_system_number"))

	response := r.FieldActionService.ChangeStatusFieldAction(FieldActionId)

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) ChangeStatusFieldActionVehicle(writer http.ResponseWriter, request *http.Request) {

	FieldActionVehicleId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_system_number"))

	response := r.FieldActionService.ChangeStatusFieldActionVehicle(FieldActionVehicleId)

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *FieldActionControllerImpl) ChangeStatusFieldActionVehicleItem(writer http.ResponseWriter, request *http.Request) {

	FieldActionVehicleItemId, _ := strconv.Atoi(chi.URLParam(request, "field_action_eligible_vehicle_item_system_number"))

	response := r.FieldActionService.ChangeStatusFieldActionVehicleItem(FieldActionVehicleItemId)

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
