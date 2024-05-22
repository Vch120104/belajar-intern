package masterwarehousecontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"errors"

	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type WarehouseLocationControllerImpl struct {
	WarehouseLocationService masterwarehouseservice.WarehouseLocationService
}

type WarehouseLocationController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	ChangeStatus(writer http.ResponseWriter, request *http.Request)
}

func NewWarehouseLocationController(WarehouseLocationService masterwarehouseservice.WarehouseLocationService) WarehouseLocationController {
	return &WarehouseLocationControllerImpl{
		WarehouseLocationService: WarehouseLocationService,
	}
}

// @Summary Get All Warehouse Location
// @Description Get All Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param is_active query bool false "Is Active"
// @Param warehouse_location_code query string false "Warehouse Location Code"
// @Param warehouse_location_name query string false "Warehouse Location Name"
// @Param company_id query string false "Company Id"
// @Param warehouse_location_detail_name query string false "Warehouse Location Detail Name"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/ [get]
func (r *WarehouseLocationControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	sortOf := queryValues.Get("sort_of")
	sortBy := queryValues.Get("sort_by")
	warehouseLocationName := queryValues.Get("warehouse_location_code")
	companyId := queryValues.Get("company_id")
	warehouseLocationCode := queryValues.Get("warehouse_location_name")
	warehouseLocationDetailName := queryValues.Get("warehouse_location_detail_name")
	isActive := queryValues.Get("is_active")

	get, err := r.WarehouseLocationService.GetAll(masterwarehousepayloads.GetAllWarehouseLocationRequest{
		WarehouseLocationCode:       warehouseLocationName,
		WarehouseLocationName:       warehouseLocationCode,
		WarehouseLocationDetailName: warehouseLocationDetailName,
		CompanyId:                   companyId,
		IsActive:                    isActive,
	}, pagination.Pagination{
		Limit:  limit,
		SortOf: sortOf,
		SortBy: sortBy,
		Page:   page,
	})

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Get Warehouse Location By Id
// @Description Get Warehouse Location By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Param warehouse_location_id path int true "warehouse_location_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/{warehouse_location_id} [get]
func (r *WarehouseLocationControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {

	warehouseLocationId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_location_id"))

	get, err := r.WarehouseLocationService.GetById(warehouseLocationId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Save Warehouse Location
// @Description Save Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @param reqBody body masterwarehousepayloads.GetWarehouseLocationResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/ [post]
func (r *WarehouseLocationControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	var message string
	var formRequest masterwarehousepayloads.GetWarehouseLocationResponse
	helper.ReadFromRequestBody(request, &formRequest)

	save, err := r.WarehouseLocationService.Save(formRequest)

	if formRequest.WarehouseLocationId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, save, message, http.StatusOK)

}

// @Summary Change Warehouse Location Status By Id
// @Description Change Warehouse Location Status By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Param warehouse_location_id path int true "Warehouse Location Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/{warehouse_location_id} [patch]
func (r *WarehouseLocationControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request) {

	warehouseLocationId, _ := strconv.Atoi(chi.URLParam(request, "warehouse_location_id"))

	change_status, err := r.WarehouseLocationService.ChangeStatus(warehouseLocationId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, change_status, "Updated successfully", http.StatusOK)

}
