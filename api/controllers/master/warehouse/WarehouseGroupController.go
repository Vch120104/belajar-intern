package masterwarehousecontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousegrouppayloads "after-sales/api/payloads/master/warehouse"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type WarehouseGroupControllerImpl struct {
	WarehouseGroupService masterwarehousegroupservice.WarehouseGroupService
}

type WarehouseGroupController interface {
	GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Save(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

func NewWarehouseGroupController(WarehouseGroupService masterwarehousegroupservice.WarehouseGroupService) WarehouseGroupController {
	return &WarehouseGroupControllerImpl{
		WarehouseGroupService: WarehouseGroupService,
	}
}

// @Summary Get All Warehouse Group
// @Description Get All Warehouse Group
// @Accept json
// @Produce json
// @Tags Master : Warehouse Group
// @Security BearerAuth
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param is_active query bool false "is_active"
// @Param warehouse_group_code query string false "Warehouse Group Code"
// @Param warehouse_group_name query string false "Warehouse Group Name"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-group [get]
func (r *WarehouseGroupControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()

	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	sortOf := queryValues.Get("sort_of")
	sortBy := queryValues.Get("sort_by")
	isActive := queryValues.Get("is_active")
	warehouseGroupCode := queryValues.Get("warehouse_group_code")
	warehouseGroupName := queryValues.Get("warehouse_group_name")

	get := r.WarehouseGroupService.GetAll(masterwarehousegrouppayloads.GetAllWarehouseGroupRequest{
		IsActive:           isActive,
		WarehouseGroupCode: warehouseGroupCode,
		WarehouseGroupName: warehouseGroupName,
	})

	result, totalPages, totalRows := utils.DataFramePaginate(get, page, limit, sortOf, sortBy)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", 200, limit, page, int64(totalRows), totalPages)
}

// @Summary Get Warehouse Group By Id
// @Description Get Warehouse Group By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Group
// @Security BearerAuth
// @Param warehouse_group_id path int true "warehouse_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-group/{warehouse_group_id} [get]
func (r *WarehouseGroupControllerImpl) GetById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	warehouseGroupId, _ := strconv.Atoi(params.ByName("warehouse_group_id"))

	get := r.WarehouseGroupService.GetById(warehouseGroupId)

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Save Warehouse Group
// @Description Save Warehouse Group
// @Accept json
// @Produce json
// @Tags Master : Warehouse Group
// @Security BearerAuth
// @param reqBody body masterwarehousegrouppayloads.GetWarehouseGroupResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-group [post]
func (r *WarehouseGroupControllerImpl) Save(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var message string
	var formRequest masterwarehousegrouppayloads.GetWarehouseGroupResponse
	helper.ReadFromRequestBody(request, &formRequest)

	save := r.WarehouseGroupService.Save(formRequest)

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
// @Security BearerAuth
// @Param warehouse_group_id path int true "warehouse_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-group/{warehouse_group_id} [patch]
func (r *WarehouseGroupControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	warehouseGroupId, _ := strconv.Atoi(params.ByName("warehouse_group_id"))

	change_status := r.WarehouseGroupService.ChangeStatus(warehouseGroupId)

	payloads.NewHandleSuccess(writer, change_status, "Updated successfully", http.StatusOK)

}
