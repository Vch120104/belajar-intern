package masterwarehousecontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"after-sales/api/validation"
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
	GetWarehouseMasterByCodeCompany(writer http.ResponseWriter, request *http.Request)
	GetWarehouseWithMultiId(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Update(writer http.ResponseWriter, request *http.Request)
	ChangeStatus(writer http.ResponseWriter, request *http.Request)
	DropdownbyGroupId(writer http.ResponseWriter, request *http.Request)
	GetAuthorizeUser(writer http.ResponseWriter, request *http.Request)
	PostAuthorizeUser(writer http.ResponseWriter, request *http.Request)
	DeleteMultiIdAuthorizeUser(writer http.ResponseWriter, request *http.Request)
	InTransitWarehouseCodeDropdown(writer http.ResponseWriter, request *http.Request)
}

func NewWarehouseMasterController(WarehouseMasterService masterwarehouseservice.WarehouseMasterService) WarehouseMasterController {
	return &WarehouseMasterControllerImpl{
		WarehouseMasterService: WarehouseMasterService,
	}
}

// DropdownbyGroupId implements WarehouseMasterController.
func (r *WarehouseMasterControllerImpl) DropdownbyGroupId(writer http.ResponseWriter, request *http.Request) {

	warehouseDropDownGroupId, err := strconv.Atoi(chi.URLParam(request, "warehouse_group_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Warehouse Group ID", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(chi.URLParam(request, "company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	get, baseErr := r.WarehouseMasterService.DropdownbyGroupId(warehouseDropDownGroupId, companyId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Warehouse Group ID not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
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
		"mtr_warehouse_master.warehouse_name":      queryValues.Get("warehouse_name"),
		"mtr_warehouse_master.warehouse_code":      queryValues.Get("warehouse_code"),
		"mtr_warehouse_group.warehouse_group_name": queryValues.Get("warehouse_group_name"),
		"mtr_warehouse_master.is_active":           queryValues.Get("is_active"),
		"mtr_warehouse_master.company_id":          queryValues.Get("company_id"),
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

	warehouseIdStr := chi.URLParam(request, "warehouse_id")
	warehouseId, err := strconv.Atoi(warehouseIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Warehouse ID", http.StatusBadRequest)
		return
	}

	getbyid, baseErr := r.WarehouseMasterService.GetById(warehouseId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Warehouse ID not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccess(writer, getbyid, "Get Data Successfully!", http.StatusOK)
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

	if code == "" {
		payloads.NewHandleError(writer, "Invalid Warehouse Code", http.StatusBadRequest)
		return
	}

	get, err := r.WarehouseMasterService.GetWarehouseMasterByCode(code)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(get), "Get Data Successfully!", http.StatusOK)

}

func (r *WarehouseMasterControllerImpl) GetWarehouseMasterByCodeCompany(writer http.ResponseWriter, request *http.Request) {
	warehouseCode := chi.URLParam(request, "warehouse_code")
	if warehouseCode == "" {
		payloads.NewHandleError(writer, "Invalid Warehouse Code", http.StatusBadRequest)
		return
	}

	companyId, errA := strconv.Atoi(chi.URLParam(request, "company_id"))
	if errA != nil || companyId == 0 {
		payloads.NewHandleError(writer, "Invalid Company Id", http.StatusBadRequest)
		return
	}

	get, err := r.WarehouseMasterService.GetWarehouseMasterByCodeCompany(warehouseCode, companyId)
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
	if warehouse_ids == "" {
		payloads.NewHandleError(writer, "Warehouse IDs are required", http.StatusBadRequest)
		return
	}

	sliceOfIdsStr := strings.Split(warehouse_ids, ",")
	var sliceOfIdsInt []int

	for _, idStr := range sliceOfIdsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			payloads.NewHandleError(writer, "Invalid warehouse ID: "+idStr, http.StatusBadRequest)
			return
		}
		sliceOfIdsInt = append(sliceOfIdsInt, id)
	}

	result, err := r.WarehouseMasterService.GetWarehouseWithMultiId(sliceOfIdsInt)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
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

	formRequest := masterwarehousepayloads.GetWarehouseMasterResponse{}
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	save, err := r.WarehouseMasterService.Save(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, save, "Create Data Successfully!", http.StatusCreated)
}

// @Summary Update Warehouse Master
// @Description Update Warehouse Master
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Param warehouse_id path int true "warehouse_id"
// @param reqBody body masterwarehousepayloads.GetWarehouseMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-master/{warehouse_id}/{company_id} [put]
func (r *WarehouseMasterControllerImpl) Update(writer http.ResponseWriter, request *http.Request) {

	warehouseId, err := strconv.Atoi(chi.URLParam(request, "warehouse_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Warehouse ID", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(chi.URLParam(request, "company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	formRequest := masterwarehousepayloads.UpdateWarehouseMasterRequest{}
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	update, baseErr := r.WarehouseMasterService.Update(warehouseId, companyId, formRequest)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Warehouse ID not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, update, "Update Data Successfully!", http.StatusOK)
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

	warehouseId, err := strconv.Atoi(chi.URLParam(request, "warehouse_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Warehouse ID", http.StatusBadRequest)
		return
	}

	changeStatus, baseErr := r.WarehouseMasterService.ChangeStatus(warehouseId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Warehouse ID not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, changeStatus, "Change Status Successfully!", http.StatusOK)

}

func (r *WarehouseMasterControllerImpl) GetAuthorizeUser(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filter := map[string]string{
		"mtr_warehouse_authorize.warehouse_authorize_id": queryValues.Get("warehouse_authorize_id"),
		"mtr_warehouse_authorize.employee_id":            queryValues.Get("employee_id"),
		"mtr_user_details.user_id":                       queryValues.Get("user_id"),
		"mtr_user_details.employee_name":                 queryValues.Get("employee_name"),
		"mtr_warehouse_authorize.company_id":             queryValues.Get("company_id"),
		"mtr_warehouse_authorize.warehouse_id":           queryValues.Get("warehouse_id"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(filter)

	result, err := r.WarehouseMasterService.GetAuthorizeUser(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

func (r *WarehouseMasterControllerImpl) PostAuthorizeUser(writer http.ResponseWriter, request *http.Request) {
	formRequest := masterwarehousepayloads.WarehouseAuthorize{}
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	save, err := r.WarehouseMasterService.PostAuthorizeUser(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, save, "data saved succesfully", http.StatusOK)
}

func (r *WarehouseMasterControllerImpl) DeleteMultiIdAuthorizeUser(writer http.ResponseWriter, request *http.Request) {
	warehouseAuthorizeId := chi.URLParam(request, "warehouse_authorize_id")
	delete, err := r.WarehouseMasterService.DeleteMultiIdAuthorizeUser(warehouseAuthorizeId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, delete, "data deleted succesfully", http.StatusOK)
}

func (r *WarehouseMasterControllerImpl) InTransitWarehouseCodeDropdown(writer http.ResponseWriter, request *http.Request) {

	companyID, err := strconv.Atoi(chi.URLParam(request, "company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	warehouseGroupId, err := strconv.Atoi(chi.URLParam(request, "warehouse_group_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Warehouse Group ID", http.StatusBadRequest)
		return
	}

	get, errResp := r.WarehouseMasterService.InTransitWarehouseCodeDropdown(companyID, warehouseGroupId)
	if errResp != nil {
		helper.ReturnError(writer, request, errResp)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}
