package masterwarehousecontroller

import (
	"after-sales/api/payloads"

	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WarehouseLocationControllerImpl struct {
	WarehouseLocationService masterwarehouseservice.WarehouseLocationService
}

type WarehouseLocationController interface {
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
// @Security BearerAuth
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
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-location [get]
func (r *WarehouseLocationControllerImpl) GetAll(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")
	warehouseLocationName := c.Query("warehouse_location_code")
	companyId := c.Query("company_id")
	warehouseLocationCode := c.Query("warehouse_location_name")
	warehouseLocationDetailName := c.Query("warehouse_location_detail_name")
	isActive := c.Query("is_active")

	get := r.WarehouseLocationService.GetAll(masterwarehousepayloads.GetAllWarehouseLocationRequest{
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

	payloads.HandleSuccessPagination(c, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Get Warehouse Location By Id
// @Description Get Warehouse Location By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Security BearerAuth
// @Param warehouse_location_id path int true "warehouse_location_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-location/{warehouse_location_id} [get]
func (r *WarehouseLocationControllerImpl) GetById(c *gin.Context) {

	warehouseLocationId, _ := strconv.Atoi(c.Param("warehouse_location_id"))

	get := r.WarehouseLocationService.GetById(warehouseLocationId)

	payloads.HandleSuccess(c, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Save Warehouse Location
// @Description Save Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Security BearerAuth
// @param reqBody body masterwarehousepayloads.GetWarehouseLocationResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-location [post]
func (r *WarehouseLocationControllerImpl) Save(c *gin.Context) {

	var message string
	requestBody := masterwarehousepayloads.GetWarehouseLocationResponse{}

	save := r.WarehouseLocationService.Save(requestBody)

	if requestBody.WarehouseLocationId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, save, message, http.StatusOK)

}

// @Summary Change Warehouse Location Status By Id
// @Description Change Warehouse Location Status By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Security BearerAuth
// @Param warehouse_location_id path int true "Warehouse Location Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-location/{warehouse_location_id} [patch]
func (r *WarehouseLocationControllerImpl) ChangeStatus(c *gin.Context) {

	warehouseLocationId, _ := strconv.Atoi(c.Param("warehouse_location_id"))

	change_status := r.WarehouseLocationService.ChangeStatus(warehouseLocationId)

	payloads.HandleSuccess(c, change_status, "Updated successfully", http.StatusOK)

}
