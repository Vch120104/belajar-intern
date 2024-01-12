package masterwarehousecontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousegrouppayloads "after-sales/api/payloads/master/warehouse"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseGroupController struct {
	warehouseGroupService masterwarehousegroupservice.WarehouseGroupService
}

func OpenWarehouseGroupRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	warehouseGroupService masterwarehousegroupservice.WarehouseGroupService,
) {
	handler := WarehouseGroupController{
		warehouseGroupService: warehouseGroupService,
	}

	// r.Use(middlewares.SetupAuthenticationMiddleware())
	r.GET("/warehouse-group/:warehouse_group_id", middlewares.DBTransactionMiddleware(db), handler.GetById)
	r.GET("/warehouse-group", middlewares.DBTransactionMiddleware(db), handler.GetAll)
	r.POST("/warehouse-group", middlewares.DBTransactionMiddleware(db), handler.Save)
	r.PATCH("/warehouse-group/:warehouse_group_id", middlewares.DBTransactionMiddleware(db), handler.ChangeStatus)
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
// @Param is_active query bool false "Is Active"
// @Param warehouse_group_code query string false "Warehouse Group Code"
// @Param warehouse_group_name query string false "Warehouse Group Name"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/warehouse-group [get]
func (r *WarehouseGroupController) GetAll(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")
	warehouseGroupCode := c.Query("warehouse_group_code")
	warehouseGroupName := c.Query("warehouse_group_name")

	get, err := r.warehouseGroupService.WithTrx(trxHandle).GetAll(masterwarehousegrouppayloads.GetAllWarehouseGroupRequest{
		WarehouseGroupCode: warehouseGroupCode,
		WarehouseGroupName: warehouseGroupName,
	})

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if len(get) == 0 {
		exceptions.NotFoundException(c, "Nothing matching request")
		return
	}

	result, totalPages, totalRows := utils.DataFramePaginate(get, page, limit, sortOf, sortBy)

	payloads.HandleSuccessPagination(c, utils.ModifyKeysInResponse(result), "Get Data Successfully!", 200, limit, page, int64(totalRows), totalPages)
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
func (r *WarehouseGroupController) GetById(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	warehouseGroupId, _ := strconv.Atoi(c.Param("warehouse_group_id"))

	get, err := r.warehouseGroupService.WithTrx(trxHandle).GetById(warehouseGroupId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if get.WarehouseGroupId == 0 {
		exceptions.NotFoundException(c, "Warehouse Group Data Not Found!")
		return
	}

	payloads.HandleSuccess(c, get, "Get Data Successfully!", http.StatusOK)

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
func (r *WarehouseGroupController) Save(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var message string
	requestBody := masterwarehousegrouppayloads.GetWarehouseGroupResponse{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(requestBody.WarehouseGroupId) != 0 {
		result, err := r.warehouseGroupService.WithTrx(trxHandle).GetById(int(requestBody.WarehouseGroupId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.WarehouseGroupId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	save, err := r.warehouseGroupService.WithTrx(trxHandle).Save(requestBody)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if requestBody.WarehouseGroupId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, save, message, http.StatusOK)

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
func (r *WarehouseGroupController) ChangeStatus(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	warehouseGroupId, _ := strconv.Atoi(c.Param("warehouse_group_id"))

	get, err := r.warehouseGroupService.WithTrx(trxHandle).GetById(warehouseGroupId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if get.WarehouseGroupId == 0 {
		exceptions.NotFoundException(c, "Warehouse Group Data Not Found!")
		return
	}

	change_status, err := r.warehouseGroupService.WithTrx(trxHandle).ChangeStatus(warehouseGroupId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, change_status, "Updated successfully", http.StatusOK)

}
