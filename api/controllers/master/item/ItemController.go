package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
	"strconv"
	"strings"

	masteritemservice "after-sales/api/services/master/item"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemController struct {
	itemservice masteritemservice.ItemService
}

func StartItemRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	itemservice masteritemservice.ItemService,
) {
	itemHandler := ItemController{itemservice: itemservice}
	r.GET("/item/", middlewares.DBTransactionMiddleware(db), itemHandler.GetAllItem)
	r.GET("/item/pop-up", middlewares.DBTransactionMiddleware(db), itemHandler.GetAllItemLookup)
	r.GET("/item-multi-id/:item_ids", middlewares.DBTransactionMiddleware(db), itemHandler.GetItemWithMultiId)
	r.GET("/item/:item_code", middlewares.DBTransactionMiddleware(db), itemHandler.GetItemByCode)
	r.POST("/item", middlewares.DBTransactionMiddleware(db), itemHandler.SaveItem)
	r.PATCH("/item/:item_id", middlewares.DBTransactionMiddleware(db), itemHandler.ChangeStatusItem)
}

// @Summary Get All Item
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_code query string false "item_code"
// @Param item_name query string false "item_name"
// @Param item_type query string false "item_type"
// @Param is_active query string false "is_active"
// @Param item_class_code query string false "item_class_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item/ [get]
func (r *ItemController) GetAllItem(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)

	queryParams := map[string]string{
		"mtr_item.item_code":             c.Query("item_code"),
		"mtr_item.item_name":             c.Query("item_name"),
		"mtr_item.item_type":             c.Query("item_type"),
		"mtr_item_class.item_class_code": c.Query("item_class_code"),
		"mtr_item.is_active":             c.Query("is_active"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.itemservice.WithTrx(trxHandle).GetAllItem(criteria)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "success", 200)
}

// @Summary Get All Item Lookup
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_code query string false "item_code"
// @Param item_name query string false "item_name"
// @Param item_type query string false "item_type"
// @Param item_group_code query string false "item_group_code"
// @Param item_class_code query string false "item_class_code"
// @Param supplier_code query string false "supplier_code"
// @Param supplier_name query string false "supplier_name"
// @Param is_active query string false "is_active"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item/pop-up [get]
func (r *ItemController) GetAllItemLookup(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	queryParams := map[string]string{
		"item_code":       c.Query("item_code"),
		"item_name":       c.Query("item_name"),
		"item_type":       c.Query("item_type"),
		"item_group_code": c.Query("item_group_code"),
		"item_class_code": c.Query("item_class_code"),
		"supplier_code":   c.Query("supplier_code"),
		"supplier_name":   c.Query("supplier_name"),
		"is_active":       c.Query("is_active"),
		"sort_of":         c.Query("sort_of"),
		"sort_by":         c.Query("sort_by"),
		"limit":           c.Query("limit"),
		"page":            c.Query("page"),
	}

	result, err := r.itemservice.WithTrx(trxHandle).GetAllItemLookup(queryParams)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	// paginatedData, totalPages, totalRows := utils.DataFramePaginate(result, 0, 0, SnaketoPascalCase(sortOf), sortBy)

	// payloads.HandleSuccessPagination(c, modifyKeysInResponse(paginatedData), "Get Data Successfully!", 200, limit, page, int64(totalRows), totalPages)
	payloads.HandleSuccessPagination(c, utils.ModifyKeysInResponse(result), "Get Data Successfully!", 200, 0, 0, int64(0), 0)
}

// @Summary Get Item With MultiId
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_ids path string true "item_ids"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item-multi-id/{item_ids} [get]
func (r *ItemController) GetItemWithMultiId(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	item_ids := c.Param("item_ids")

	sliceOfString := strings.Split(item_ids, ",")

	result, err := r.itemservice.WithTrx(trxHandle).GetItemWithMultiId(sliceOfString)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "success", 200)
}

// @Summary Get Item By code
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_code path string true "item_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item/{item_code} [get]
func (r *ItemController) GetItemByCode(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	itemCode := c.Param("item_code")

	result, err := r.itemservice.WithTrx(trxHandle).GetItemCode(itemCode)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @param reqBody body masteritempayloads.ItemRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item [post]
func (r *ItemController) SaveItem(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteritempayloads.ItemResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.ItemId) != 0 {
		result, err := r.itemservice.WithTrx(trxHandle).GetItemById(int(request.ItemId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.ItemId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.itemservice.WithTrx(trxHandle).SaveItem(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.ItemId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Item
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @param item_id path int true "item_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/item/{item_id} [patch]
func (r *ItemController) ChangeStatusItem(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	ItemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.itemservice.WithTrx(trxHandle).GetItemById(int(ItemId))
	if err != nil || result.ItemId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.itemservice.WithTrx(trxHandle).ChangeStatusItem(int(ItemId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Change Status Successfully!", http.StatusOK)
}
