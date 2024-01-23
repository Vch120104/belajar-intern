package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
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

func StartItemLevelRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	itemLevelService masteritemlevelservice.ItemLevelService,
) {
	handler := ItemLevelController{
		itemLevelService: itemLevelService,
	}
	r.GET("/item-level-by-id", middlewares.DBTransactionMiddleware(db), handler.GetById)
	r.GET("/item-level", middlewares.DBTransactionMiddleware(db), handler.GetAll)
	r.POST("/item-level", middlewares.DBTransactionMiddleware(db), handler.Save)
	r.PATCH("/item-level/:item_level_id", middlewares.DBTransactionMiddleware(db), handler.ChangeStatus)
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
// @Router /aftersales-service/api/aftersales/item-level [get]
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

	payloads.HandleSuccessPagination(c, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
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
// @Router /aftersales-service/api/aftersales/item-level-by-id [get]
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

// @Summary Save Item Level
// @Description Save Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security BearerAuth
// @param reqBody body masteritemlevelpayloads.SaveItemLevelRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-level [post]
func (r *ItemLevelController) Save(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteritemlevelpayloads.SaveItemLevelRequest
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.ItemLevelId) != 0 {
		result, err := r.itemLevelService.WithTrx(trxHandle).GetById(int(request.ItemLevelId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.ItemLevelId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.itemLevelService.WithTrx(trxHandle).Save(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.ItemLevelId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
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
// @Router /aftersales-service/api/aftersales/item-level/{item_level_id} [patch]
func (r *ItemLevelController) ChangeStatus(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	itemLevelId, err := strconv.Atoi(c.Param("item_level_id"))
	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.itemLevelService.WithTrx(trxHandle).GetById(int(itemLevelId))
	if err != nil || result.ItemLevelId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.itemLevelService.WithTrx(trxHandle).ChangeStatus(int(itemLevelId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}
