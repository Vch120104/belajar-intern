package masterwarehousecontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"strconv"
	"strings"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseMasterController struct {
	warehouseMasterService masterwarehouseservice.WarehouseMasterService
}

func OpenWarehouseMasterRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	warehouseMasterService masterwarehouseservice.WarehouseMasterService,
) {
	handler := WarehouseMasterController{
		warehouseMasterService: warehouseMasterService,
	}

	// r.Use(middlewares.SetupAuthenticationMiddleware())
	r.GET("/warehouse-master", middlewares.DBTransactionMiddleware(db), handler.GetAll)
	r.GET("/warehouse-master/:warehouse_id", middlewares.DBTransactionMiddleware(db), handler.GetById)
	r.GET("/warehouse-master-by-code/:warehouse_code", middlewares.DBTransactionMiddleware(db), handler.GetByCode)
	r.GET("/warehouse-master-multi-id/:warehouse_ids", middlewares.DBTransactionMiddleware(db), handler.GetWarehouseWithMultiId)
	r.GET("/warehouse-master-drop-down", middlewares.DBTransactionMiddleware(db), handler.GetAllIsActive)
	r.POST("/warehouse-master", middlewares.DBTransactionMiddleware(db), handler.Save)
	r.PATCH("/warehouse-master/:warehouse_id", middlewares.DBTransactionMiddleware(db), handler.ChangeStatus)
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
func (r *WarehouseMasterController) GetAll(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")
	warehouseName := c.Query("warehouse_name")
	warehouseCode := c.Query("warehouse_code")
	isActive := c.Query("is_active")

	get, err := r.warehouseMasterService.WithTrx(trxHandle).GetAll(masterwarehousepayloads.GetAllWarehouseMasterRequest{
		WarehouseName: warehouseName,
		WarehouseCode: warehouseCode,
		IsActive:      isActive,
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

// @Summary Get All Warehouse Master Is Active
// @Description Get All Warehouse Master Is Active
// @Accept json
// @Produce json
// @Tags Master : Warehouse Master
// @Security BearerAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-master-drop-down [get]
func (r *WarehouseMasterController) GetAllIsActive(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)

	get, err := r.warehouseMasterService.WithTrx(trxHandle).GetAllIsActive()

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, get, "Get Data Successfully!", http.StatusOK)
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
func (r *WarehouseMasterController) GetById(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	warehouseId, _ := strconv.Atoi(c.Param("warehouse_id"))

	get, err := r.warehouseMasterService.WithTrx(trxHandle).GetById(warehouseId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if get.WarehouseId == 0 {
		exceptions.NotFoundException(c, "Warehouse Master Data Not Found!")
		return
	}

	payloads.HandleSuccess(c, get, "Get Data Successfully!", http.StatusOK)
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
func (r *WarehouseMasterController) GetByCode(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	code := c.Param("warehouse_code")

	get, err := r.warehouseMasterService.WithTrx(trxHandle).GetWarehouseMasterByCode(code)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, utils.ModifyKeysInResponse(get), "Get Data Successfully!", http.StatusOK)

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
func (r *WarehouseMasterController) GetWarehouseWithMultiId(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	warehouse_ids := c.Param("warehouse_ids")

	sliceOfString := strings.Split(warehouse_ids, ",")

	result, err := r.warehouseMasterService.WithTrx(trxHandle).GetWarehouseWithMultiId(sliceOfString)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "success", 200)
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
func (r *WarehouseMasterController) Save(c *gin.Context) {
	var message string
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	requestBody := masterwarehousepayloads.GetWarehouseMasterResponse{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(requestBody.WarehouseId) != 0 {
		get, err := r.warehouseMasterService.WithTrx(trxHandle).GetById(requestBody.WarehouseId)

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if get.WarehouseId == 0 {
			exceptions.NotFoundException(c, "Warehouse Master Data Not Found!")
			return
		}
	}

	save, err := r.warehouseMasterService.WithTrx(trxHandle).Save(requestBody)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if requestBody.WarehouseId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, save, message, http.StatusOK)

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
func (r *WarehouseMasterController) ChangeStatus(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	warehouseId, _ := strconv.Atoi(c.Param("warehouse_id"))

	get, err := r.warehouseMasterService.WithTrx(trxHandle).GetById(warehouseId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if get.WarehouseId == 0 {
		exceptions.NotFoundException(c, "Warehouse Master Data Not Found!")
		return
	}

	change_status, err := r.warehouseMasterService.WithTrx(trxHandle).ChangeStatus(warehouseId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, change_status, "Updated successfully", http.StatusOK)

}
