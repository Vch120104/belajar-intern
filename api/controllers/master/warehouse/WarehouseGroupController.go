package masterwarehousecontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousegrouppayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WarehouseGroupControllerImpl struct {
	WarehouseGroupService masterwarehousegroupservice.WarehouseGroupService
}

type WarehouseGroupController interface {
	GetAllWarehouseGroup(writer http.ResponseWriter, request *http.Request)
	GetByIdWarehouseGroup(writer http.ResponseWriter, request *http.Request)
	SaveWarehouseGroup(writer http.ResponseWriter, request *http.Request)
	GetWarehouseGroupDropdownbyId(writer http.ResponseWriter, request *http.Request)
	GetWarehouseGroupDropDown(writer http.ResponseWriter, request *http.Request)
	ChangeStatusWarehouseGroup(writer http.ResponseWriter, request *http.Request)
	GetbyGroupCode(writer http.ResponseWriter, request *http.Request)
}

func NewWarehouseGroupController(WarehouseGroupService masterwarehousegroupservice.WarehouseGroupService) WarehouseGroupController {
	return &WarehouseGroupControllerImpl{
		WarehouseGroupService: WarehouseGroupService,
	}
}

// GetbyGroupCode implements WarehouseGroupController.
func (r *WarehouseGroupControllerImpl) GetbyGroupCode(writer http.ResponseWriter, request *http.Request) {
	groupCode := chi.URLParam(request, "warehouse_group_code")

	get, err := r.WarehouseGroupService.GetbyGroupCode(groupCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// GetWarehouseGroupDropdownbyId implements WarehouseGroupController.
func (r *WarehouseGroupControllerImpl) GetWarehouseGroupDropdownbyId(writer http.ResponseWriter, request *http.Request) {
	warehouseGroupId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_group_id"))

	get, err := r.WarehouseGroupService.GetWarehouseGroupDropdownbyId(warehouseGroupId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// GetWarehouseGroupDropDown implements WarehouseGroupController.
func (r *WarehouseGroupControllerImpl) GetWarehouseGroupDropDown(writer http.ResponseWriter, request *http.Request) {

	get, err := r.WarehouseGroupService.GetWarehouseGroupDropdown()
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Warehouse GroupF
// @Description Get All Warehouse Group
// @Accept json
// @Produce json
// @Tags Master : Warehouse Group
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param is_active query bool false "is_active"
// @Param warehouse_group_code query string false "Warehouse Group Code"
// @Param warehouse_group_name query string false "Warehouse Group Name"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-group/ [get]
func (r *WarehouseGroupControllerImpl) GetAllWarehouseGroup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	query_params := map[string]string{
		"is_active":            queryValues.Get("is_active"),
		"warehouse_group_code": queryValues.Get("warehouse_group_code"),
		"warehouse_group_name": queryValues.Get("warehouse_group_name"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(query_params)

	get, err := r.WarehouseGroupService.GetAllWarehouseGroup(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Get Warehouse Group By Id
// @Description Get Warehouse Group By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Group
// @Param warehouse_group_id path int true "warehouse_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-group/{warehouse_group_id} [get]
func (r *WarehouseGroupControllerImpl) GetByIdWarehouseGroup(writer http.ResponseWriter, request *http.Request) {

	warehouseGroupId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_group_id"))

	get, err := r.WarehouseGroupService.GetByIdWarehouseGroup(int(warehouseGroupId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Save Warehouse Group
// @Description Save Warehouse Group
// @Accept json
// @Produce json
// @Tags Master : Warehouse Group
// @param reqBody body masterwarehousegrouppayloads.GetWarehouseGroupResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-group/warehouse-group [post]
func (r *WarehouseGroupControllerImpl) SaveWarehouseGroup(writer http.ResponseWriter, request *http.Request) {

	var message string
	var formRequest masterwarehousegrouppayloads.GetWarehouseGroupResponse
	helper.ReadFromRequestBody(request, &formRequest)

	save, err := r.WarehouseGroupService.SaveWarehouseGroup(formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	if formRequest.WarehouseGroupId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, save, message, http.StatusOK)

}

// @Summary Change Warehouse Group Status By Id
// @Description Change Warehouse Group Status By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Group
// @Param warehouse_group_id path int true "warehouse_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-group/{warehouse_group_id} [patch]
func (r *WarehouseGroupControllerImpl) ChangeStatusWarehouseGroup(writer http.ResponseWriter, request *http.Request) {

	warehouseGroupId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_group_id"))

	change_status, err := r.WarehouseGroupService.ChangeStatusWarehouseGroup(int(warehouseGroupId))
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, change_status, "Updated successfully", http.StatusOK)

}
