package masterwarehousecontroller

import (
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masterwarehousegrouppayloads "after-sales/api/payloads/master/warehouse"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WarehouseGroupControllerImpl struct {
	WarehouseGroupService masterwarehousegroupservice.WarehouseGroupService
}

type WarehouseGroupController interface {
}

func NewWarehouseGroupController(WarehouseGroupService masterwarehousegroupservice.WarehouseGroupService) WarehouseGroupController {
	return &WarehouseGroupControllerImpl{
		WarehouseGroupService: WarehouseGroupService,
	}
}

// @Summary Get All Warehouse Groupfil
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
func (r *WarehouseGroupControllerImpl) GetAll(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")
	warehouseGroupCode := c.Query("warehouse_group_code")
	warehouseGroupName := c.Query("warehouse_group_name")

	get := r.WarehouseGroupService.GetAll(masterwarehousegrouppayloads.GetAllWarehouseGroupRequest{
		WarehouseGroupCode: warehouseGroupCode,
		WarehouseGroupName: warehouseGroupName,
	})

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
func (r *WarehouseGroupControllerImpl) GetById(c *gin.Context) {

	warehouseGroupId, _ := strconv.Atoi(c.Param("warehouse_group_id"))

	get := r.WarehouseGroupService.GetById(warehouseGroupId)

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
func (r *WarehouseGroupControllerImpl) Save(c *gin.Context) {

	var message string
	requestBody := masterwarehousegrouppayloads.GetWarehouseGroupResponse{}

	save := r.WarehouseGroupService.Save(requestBody)

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
func (r *WarehouseGroupControllerImpl) ChangeStatus(c *gin.Context) {

	warehouseGroupId, _ := strconv.Atoi(c.Param("warehouse_group_id"))

	change_status := r.WarehouseGroupService.ChangeStatus(warehouseGroupId)

	payloads.HandleSuccess(c, change_status, "Updated successfully", http.StatusOK)

}
