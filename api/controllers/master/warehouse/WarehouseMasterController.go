package masterwarehousecontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"net/http"
	"strconv"
	"strings"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehouseservice "after-sales/api/services/master/warehouse"

	"github.com/go-chi/chi/v5"
)

type WarehouseMasterControllerImpl struct {
	WarehouseMasterService masterwarehouseservice.WarehouseMasterService
}

type WarehouseMasterController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetAllIsActive(writer http.ResponseWriter, request *http.Request)
	DropdownWarehouse(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	GetByCode(writer http.ResponseWriter, request *http.Request)
	GetWarehouseWithMultiId(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	ChangeStatus(writer http.ResponseWriter, request *http.Request)
	DropdownbyGroupId(writer http.ResponseWriter, request *http.Request)
	GetAuthorizeUser(writer http.ResponseWriter, request *http.Request)
	PostAuthorizeUser(writer http.ResponseWriter, request *http.Request)
	DeleteMultiIdAuthorizeUser(writer http.ResponseWriter, request *http.Request)
}

func NewWarehouseMasterController(WarehouseMasterService masterwarehouseservice.WarehouseMasterService) WarehouseMasterController {
	return &WarehouseMasterControllerImpl{
		WarehouseMasterService: WarehouseMasterService,
	}
}

// DropdownbyGroupId implements WarehouseMasterController.
func (r *WarehouseMasterControllerImpl) DropdownbyGroupId(writer http.ResponseWriter, request *http.Request) {

	warehouseGroupId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_group_id"))

	get, err := r.WarehouseMasterService.DropdownbyGroupId(warehouseGroupId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Warehouse Master
// @Description Get All Warehouse Master
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param is_active query bool false "Is Active"
// @Param warehouse_name query string false "Warehouse Name"
// @Param warehouse_code query string false "Warehouse Code"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/ [get]
func (r *WarehouseMasterControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	filter := map[string]string{
		"warehouse_name":       queryValues.Get("warehouse_name"),
		"warehouse_code":       queryValues.Get("warehouse_code"),
		"warehouse_group_name": queryValues.Get("warehouse_group_name"),
		"is_active":            queryValues.Get("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(filter)

	get, err := r.WarehouseMasterService.GetAll(filterCondition, pagination)

	if err != nil {
		helper.ReturnError(writer, request, err)

		return
	}

	payloads.NewHandleSuccessPagination(writer, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Get All Warehouse Master Is Active
// @Description Get All Warehouse Master Is Active
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/is-active [get]
func (r *WarehouseMasterControllerImpl) GetAllIsActive(writer http.ResponseWriter, request *http.Request) {

	get, err := r.WarehouseMasterService.GetAllIsActive()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Dropdown Warehouse
// @Description Get Dropdown Warehouse
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/drop-down [get]
func (r *WarehouseMasterControllerImpl) DropdownWarehouse(writer http.ResponseWriter, request *http.Request) {

	get, err := r.WarehouseMasterService.DropdownWarehouse()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Warehouse Master By Id
// @Description Get Warehouse Master By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Param warehouse_id path int true "warehouse_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/{warehouse_id} [get]
func (r *WarehouseMasterControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	warehouseId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_id"))

	get, err := r.WarehouseMasterService.GetById(warehouseId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Warehouse Master By Code
// @Description Get Warehouse Master By Code
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Param warehouse_code path string true "warehouse_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/by-code/{warehouse_code} [get]
func (r *WarehouseMasterControllerImpl) GetByCode(writer http.ResponseWriter, request *http.Request) {

	code := chi.URLParam(request, "warehouse_code")

	get, err := r.WarehouseMasterService.GetWarehouseMasterByCode(code)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(get), "Get Data Successfully!", http.StatusOK)

}

// @Summary Get Warehouse Master With MultiId
// @Description Get Warehouse Master
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Param warehouse_ids path string true "warehouse_ids"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/multi-id/{warehouse_ids} [get]
func (r *WarehouseMasterControllerImpl) GetWarehouseWithMultiId(writer http.ResponseWriter, request *http.Request) {

	warehouse_ids := chi.URLParam(request, "warehouse_ids")

	sliceOfString := strings.Split(warehouse_ids, ",")

	result, err := r.WarehouseMasterService.GetWarehouseWithMultiId(sliceOfString)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Save Warehouse Master
// @Description Save Warehouse Master
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @param reqBody body masterwarehousepayloads.GetWarehouseMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/ [post]
func (r *WarehouseMasterControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	var message string

	formRequest := masterwarehousepayloads.GetWarehouseMasterResponse{}
	helper.ReadFromRequestBody(request, &formRequest)

	save, err := r.WarehouseMasterService.Save(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.WarehouseId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, save, message, http.StatusOK)

}

// @Summary Change Warehouse Master Status By Id
// @Description Change Warehouse Master Status By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Param warehouse_id path int true "warehouse_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/{warehouse_id} [patch]
func (r *WarehouseMasterControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request) {

	warehouseId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_id"))

	change_status, err := r.WarehouseMasterService.ChangeStatus(warehouseId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, change_status, "Updated successfully", http.StatusOK)

}

func (r *WarehouseMasterControllerImpl) GetAuthorizeUser(writer http.ResponseWriter, request *http.Request){
	queryValues := request.URL.Query()
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	warehouseId,_ := strconv.Atoi(chi.URLParam(request,"warehouse_id"))
	result,err := r.WarehouseMasterService.GetAuthorizeUser(pagination,warehouseId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *WarehouseMasterControllerImpl) PostAuthorizeUser(writer http.ResponseWriter, request *http.Request){
	formRequest:= masterwarehousepayloads.WarehouseAuthorize{}
	helper.ReadFromRequestBody(request, &formRequest)
	save,err := r.WarehouseMasterService.PostAuthorizeUser(formRequest)
	if err != nil{
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, save, "data saved succesfully", http.StatusOK)
}

func (r *WarehouseMasterControllerImpl) DeleteMultiIdAuthorizeUser(writer http.ResponseWriter, request *http.Request){
	warehouseAuthorizeId := chi.URLParam(request,"warehouse_authorize_id")
	delete,err:= r.WarehouseMasterService.DeleteMultiIdAuthorizeUser(warehouseAuthorizeId)
	if err != nil{
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, delete, "data deleted succesfully", http.StatusOK)
}