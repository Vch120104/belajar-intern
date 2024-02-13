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

	"github.com/julienschmidt/httprouter"
)

type WarehouseMasterControllerImpl struct {
	WarehouseMasterService masterwarehouseservice.WarehouseMasterService
}

type WarehouseMasterController interface {
	GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetWarehouseWithMultiId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Save(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

func NewWarehouseMasterController(WarehouseMasterService masterwarehouseservice.WarehouseMasterService) WarehouseMasterController {
	return &WarehouseMasterControllerImpl{
		WarehouseMasterService: WarehouseMasterService,
	}
}

// @Summary Get All Warehouse Master
// @Description Get All Warehouse Master
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Security BearerAuth
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param is_active query bool false "Is Active"
// @Param warehouse_name query string false "Warehouse Name"
// @Param warehouse_code query string false "Warehouse Code"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master [get]
func (r *WarehouseMasterControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	sortOf := queryValues.Get("sort_of")
	sortBy := queryValues.Get("sort_by")
	warehouseName := queryValues.Get("warehouse_name")
	warehouseCode := queryValues.Get("warehouse_code")
	isActive := queryValues.Get("is_active")

	get := r.WarehouseMasterService.GetAll(masterwarehousepayloads.GetAllWarehouseMasterRequest{
		WarehouseName: warehouseName,
		WarehouseCode: warehouseCode,
		IsActive:      isActive,
	}, pagination.Pagination{
		Limit:  limit,
		SortOf: sortOf,
		SortBy: sortBy,
		Page:   page,
	})

	payloads.NewHandleSuccessPagination(writer, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Get All Warehouse Master Is Active
// @Description Get All Warehouse Master Is Active
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Security BearerAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master-drop-down [get]
func (r *WarehouseMasterControllerImpl) GetAllIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	get := r.WarehouseMasterService.GetAllIsActive()

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Warehouse Master By Id
// @Description Get Warehouse Master By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Security BearerAuth
// @Param warehouse_id path int true "warehouse_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master/{warehouse_id} [get]
func (r *WarehouseMasterControllerImpl) GetById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	warehouseId, _ := strconv.Atoi(params.ByName("warehouse_id"))

	get := r.WarehouseMasterService.GetById(warehouseId)

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Warehouse Master By Code
// @Description Get Warehouse Master By Code
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Security BearerAuth
// @Param warehouse_code path string true "warehouse_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master-by-code/{warehouse_code} [get]
func (r *WarehouseMasterControllerImpl) GetByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	code := params.ByName("warehouse_code")

	get := r.WarehouseMasterService.GetWarehouseMasterByCode(code)

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(get), "Get Data Successfully!", http.StatusOK)

}

// @Summary Get Warehouse Master With MultiId
// @Description Get Warehouse Master
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Param warehouse_ids path string true "warehouse_ids"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master-multi-id/{warehouse_ids} [get]
func (r *WarehouseMasterControllerImpl) GetWarehouseWithMultiId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	warehouse_ids := params.ByName("warehouse_ids")

	sliceOfString := strings.Split(warehouse_ids, ",")

	result := r.WarehouseMasterService.GetWarehouseWithMultiId(sliceOfString)

	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Save Warehouse Master
// @Description Save Warehouse Master
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Security BearerAuth
// @param reqBody body masterwarehousepayloads.GetWarehouseMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master [post]
func (r *WarehouseMasterControllerImpl) Save(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var message string

	formRequest := masterwarehousepayloads.GetWarehouseMasterResponse{}
	helper.ReadFromRequestBody(request, &formRequest)

	save := r.WarehouseMasterService.Save(formRequest)

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
// @Security BearerAuth
// @Param warehouse_id path int true "warehouse_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master/{warehouse_id} [patch]
func (r *WarehouseMasterControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	warehouseId, _ := strconv.Atoi(params.ByName("warehouse_id"))

	change_status := r.WarehouseMasterService.ChangeStatus(warehouseId)

	payloads.NewHandleSuccess(writer, change_status, "Updated successfully", http.StatusOK)

}
