package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	"fmt"
	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masteritemlevelservice "after-sales/api/services/master/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemLevelController struct {
	itemLevelService masteritemlevelservice.ItemLevelService
}

func OpenItemLevelRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	itemLevelService masteritemlevelservice.ItemLevelService,
) {
	handler := ItemLevelController{
		itemLevelService: itemLevelService,
	}

	// r.Use(middlewares.SetupAuthenticationMiddleware())
	r.POST("/save-item-level", middlewares.DBTransactionMiddleware(db), handler.Save)
	r.PUT("/update-item-level", middlewares.DBTransactionMiddleware(db), handler.Update)
	r.GET("/get-item-level-by-id", middlewares.DBTransactionMiddleware(db), handler.GetById)
	r.GET("/item-level", middlewares.DBTransactionMiddleware(db), handler.GetAll)
	r.PATCH("/change-item-level-status/:item_level_id", middlewares.DBTransactionMiddleware(db), handler.ChangeStatus)
}

// @Summary Save Item Level
// @Description Save Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @param reqBody body masteritemlevelpayloads.SaveItemLevelRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/save-item-level [post]
func (r *ItemLevelController) Save(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	requestBody := masteritemlevelpayloads.SaveItemLevelRequest{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	save, err := r.itemLevelService.WithTrx(trxHandle).Save(requestBody)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, save, "Success insert item level", http.StatusOK)

}

// @Summary Update Item Level
// @Description Update Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @param reqBody body masteritemlevelpayloads.SaveItemLevelRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/update-item-level [put]
func (r *ItemLevelController) Update(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	requestBody := masteritemlevelpayloads.SaveItemLevelRequest{}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if requestBody.ItemLevelId == 0 {
		exceptions.EntityException(c, "Item Level Id cannot be empty!")
		return
	}

	get, errGet := r.itemLevelService.WithTrx(trxHandle).GetById(requestBody.ItemLevelId)

	if errGet != nil {
		exceptions.AppException(c, errGet.Error())
		return
	}

	if get.ItemLevelId == 0 {
		exceptions.AppException(c, "Item Level Id does not exist")
		return
	}

	update, errUpdate := r.itemLevelService.WithTrx(trxHandle).Update(requestBody)

	if errUpdate != nil {
		exceptions.AppException(c, errUpdate.Error())
		return
	}

	payloads.HandleSuccess(c, update, "Success update item level", http.StatusOK)

}

// @Summary Get Item Level By Id
// @Description Get Item Level By Id
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @Param item_level_id path string true "item_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/get-item-level-by-id [get]
func (r *ItemLevelController) GetById(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	itemLevelId, _ := strconv.Atoi(c.Param("item_level_id"))

	get, err := r.itemLevelService.WithTrx(trxHandle).GetById(itemLevelId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if get.ItemLevelId == 0 {
		exceptions.NotFoundException(c, "Item Level Data Not Found!")
		return
	}

	payloads.HandleSuccess(c, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Get All Item Level
// @Description Get All Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Param item_level query string false "Item Level"
// @Param item_class_code query string false "Item Class Code"
// @Param item_level_parent query string false "Item Level Parent"
// @Param item_level_code query string false "Item Level Code"
// @Param item_level_name query string false "Item Level Name"
// @Param is_active query bool false "Is Active"
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/item-level [get]
func (r *ItemLevelController) GetAll(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")
	itemLevel := c.Query("item_level")
	itemClassCode := c.Query("item_class_code")
	itemLevelParent := c.Query("item_level_parent")
	itemLevelCode := c.Query("item_level_code")
	itemLevelName := c.Query("item_level_name")
	isActive := c.Query("is_active")

	get, err := r.itemLevelService.WithTrx(trxHandle).GetAll(masteritemlevelpayloads.GetAllItemLevelResponse{
		ItemLevel:       itemLevel,
		ItemClassCode:   itemClassCode,
		ItemLevelParent: itemLevelParent,
		ItemLevelCode:   itemLevelCode,
		ItemLevelName:   itemLevelName,
		IsActive:        isActive,
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

	fmt.Println(get.Rows)

	payloads.HandleSuccessPagination(c, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Change Item Level Status By Id
// @Description Change Item Level Status By Id
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @Param item_level_id path string true "item_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/change-item-level-status/{item_level_id} [patch]
func (r *ItemLevelController) ChangeStatus(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	itemLevelId, _ := strconv.Atoi(c.Param("item_level_id"))

	change_status, err := r.itemLevelService.WithTrx(trxHandle).ChangeStatus(itemLevelId)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, change_status, "Updated successfully", http.StatusOK)

}
