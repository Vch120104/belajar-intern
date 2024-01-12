package masterwarehousecontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"

	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseLocationController struct {
	warehouseLocationService masterwarehouseservice.WarehouseLocationService
}

func OpenWarehouseLocationRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	WarehouseLocationService masterwarehouseservice.WarehouseLocationService,
) {
	handler := WarehouseLocationController{
		warehouseLocationService: WarehouseLocationService,
	}

	// r.Use(middlewares.SetupAuthenticationMiddleware())
	r.GET("/warehouse-location/:warehouse_location_id", middlewares.DBTransactionMiddleware(db), handler.GetById)
	r.GET("/warehouse-location", middlewares.DBTransactionMiddleware(db), handler.GetAll)
	r.POST("/warehouse-location", middlewares.DBTransactionMiddleware(db), handler.Save)
	r.PATCH("/warehouse-location/:warehouse_location_id", middlewares.DBTransactionMiddleware(db), handler.ChangeStatus)
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
func (r *WarehouseLocationController) GetAll(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")
	warehouseLocationName := c.Query("warehouse_location_code")
	companyId := c.Query("company_id")
	warehouseLocationCode := c.Query("warehouse_location_name")
	warehouseLocationDetailName := c.Query("warehouse_location_detail_name")
	isActive := c.Query("is_active")

	get, err := r.warehouseLocationService.WithTrx(trxHandle).GetAll(masterwarehousepayloads.GetAllWarehouseLocationRequest{
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
		exceptions.AppException(c, err.Error())
		return
	}

	if get.Rows == nil {
		exceptions.NotFoundException(c, "Nothing matching request")
		return
	}

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
func (r *WarehouseLocationController) GetById(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	warehouseLocationId, _ := strconv.Atoi(c.Param("warehouse_location_id"))

	get, err := r.warehouseLocationService.WithTrx(trxHandle).GetById(warehouseLocationId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if get.WarehouseLocationId == 0 {
		exceptions.NotFoundException(c, "Warehouse Location Data Not Found!")
		return
	}

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
func (r *WarehouseLocationController) Save(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var message string
	requestBody := masterwarehousepayloads.GetWarehouseLocationResponse{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(requestBody.WarehouseLocationId) != 0 {
		get, err := r.warehouseLocationService.WithTrx(trxHandle).GetById(requestBody.WarehouseLocationId)

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if get.WarehouseLocationId == 0 {
			exceptions.NotFoundException(c, "Warehouse Location Data Not Found!")
			return
		}
	}

	save, err := r.warehouseLocationService.WithTrx(trxHandle).Save(requestBody)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

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
func (r *WarehouseLocationController) ChangeStatus(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	warehouseLocationId, _ := strconv.Atoi(c.Param("warehouse_location_id"))

	get, err := r.warehouseLocationService.WithTrx(trxHandle).GetById(warehouseLocationId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if get.WarehouseLocationId == 0 {
		exceptions.NotFoundException(c, "Warehouse Location Data Not Found!")
		return
	}

	change_status, err := r.warehouseLocationService.WithTrx(trxHandle).ChangeStatus(warehouseLocationId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, change_status, "Updated successfully", http.StatusOK)

}
