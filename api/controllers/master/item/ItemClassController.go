package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemClassController struct {
	itemclassservice masteritemservice.ItemClassService
}

func StartItemClassRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	itemclassservice masteritemservice.ItemClassService,
) {
	itemclassGroupHandler := ItemClassController{itemclassservice: itemclassservice}
	r.GET("/item-class/pop-up", middlewares.DBTransactionMiddleware(db), itemclassGroupHandler.GetAllItemClassLookup)
	r.GET("/item-class/", middlewares.DBTransactionMiddleware(db), itemclassGroupHandler.GetAllItemClass)
	r.POST("/item-class/", middlewares.DBTransactionMiddleware(db), itemclassGroupHandler.SaveItemClass)
	r.PATCH("/item-class/:item-class_id", middlewares.DBTransactionMiddleware(db), itemclassGroupHandler.ChangeStatusItemClass)
}

// @Summary Get All Item Class Lookup
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param item_class_id query int false "item_class_id"
// @Param item_class_code query string false "item_class_code"
// @Param item_class_name query string false "item_class_name"
// @Param item_group_name query string false "item_group_name"
// @Param line_type_code query string false "line_type_code"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-class/pop-up [get]
func (r *ItemClassController) GetAllItemClassLookup(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)

	queryParams := map[string]string{
		"mtr_item_class.is_active":       c.Query("is_active"),
		"mtr_item_class.item_class_id":   c.Query("item_class_id"),
		"mtr_item_class.item_class_code": c.Query("item_class_code"),
		"mtr_item_class.item_class_name": c.Query("item_class_name"),
		"item_group_name":                c.Query("item_group_name"),
		"line_type_code":                 c.Query("line_type_code"),
	}

	limit := utils.GetQueryInt(c, "limit")
	page := utils.GetQueryInt(c, "page")
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.itemclassservice.WithTrx(trxHandle).GetAllItemClass(criteria)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	paginatedData, totalPages, totalRows := utils.DataFramePaginate(result, page, limit, utils.SnaketoPascalCase(sortOf), sortBy)

	payloads.HandleSuccessPagination(c, utils.ModifyKeysInResponse(paginatedData), "success", 200, limit, page, int64(totalRows), totalPages)
}

// @Summary Get All Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @Param is_active query string false "is_active" Enums(true, false)
// @Param item_class_id query int false "item_class_id"
// @Param item_class_code query string false "item_class_code"
// @Param item_class_name query string false "item_class_name"
// @Param item_group_name query string false "item_group_name"
// @Param line_type_code query string false "line_type_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-class/ [get]
func (r *ItemClassController) GetAllItemClass(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)

	queryParams := map[string]string{
		"mtr_item_class.is_active":       c.Query("is_active"),
		"mtr_item_class.item_class_id":   c.Query("item_class_id"),
		"mtr_item_class.item_class_code": c.Query("item_class_code"),
		"mtr_item_class.item_class_name": c.Query("item_class_name"),
		"item_group_name":                c.Query("item_group_name"),
		"line_type_code":                 c.Query("line_type_code"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.itemclassservice.WithTrx(trxHandle).GetAllItemClass(criteria)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, utils.ModifyKeysInResponse(result), "success", 200)
}

// @Summary Save Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @param reqBody body masteritempayloads.ItemClassResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-class [post]
func (r *ItemClassController) SaveItemClass(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteritempayloads.ItemClassResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.ItemClassId) != 0 {
		result, err := r.itemclassservice.WithTrx(trxHandle).GetItemClassById(int(request.ItemClassId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.ItemClassId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.itemclassservice.WithTrx(trxHandle).SaveItemClass(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.ItemClassId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @param item_class_id path int true "item_class_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-class/{item_class_id} [patch]
func (r *ItemClassController) ChangeStatusItemClass(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	itemclassGroupId, err := strconv.Atoi(c.Param("item_class_id"))
	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.itemclassservice.WithTrx(trxHandle).GetItemClassById(int(itemclassGroupId))
	if err != nil || result.ItemClassId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.itemclassservice.WithTrx(trxHandle).ChangeStatusItemClass(int(itemclassGroupId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}
