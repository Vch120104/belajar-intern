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

type DiscountPercentController struct {
	discountpercentservice masteritemservice.DiscountPercentService
}

func StartDiscountPercentRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	discountpercentservice masteritemservice.DiscountPercentService,
) {
	discountpercentGroupHandler := DiscountPercentController{discountpercentservice: discountpercentservice}
	r.GET("/discount-percent", middlewares.DBTransactionMiddleware(db), discountpercentGroupHandler.GetAllDiscountPercent)
	r.GET("/discount-percent/:discount_percent_id", middlewares.DBTransactionMiddleware(db), discountpercentGroupHandler.GetDiscountPercentByID)
	r.POST("/discount-percent", middlewares.DBTransactionMiddleware(db), discountpercentGroupHandler.SaveDiscountPercent)
	r.PATCH("/discount-percent/:discount_percent_id", middlewares.DBTransactionMiddleware(db), discountpercentGroupHandler.ChangeStatusDiscountPercent)
}

// @Summary Get All Discount Percent
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param discount_code_value query string false "discount_code_value"
// @Param discount_code_description query string false "discount_code_description"
// @Param order_type_name query string false "order_type_name"
// @Param discount query float64 false "discount"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent [get]
func (r *DiscountPercentController) GetAllDiscountPercent(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)

	queryParams := map[string]string{
		"mtr_discount.discount_code_value":       c.Query("discount_code_value"),
		"mtr_discount.discount_code_description": c.Query("discount_code_description"),
		"order_type_name":                        c.Query("order_type_name"),
		"mtr_discount_percent.discount":          c.Query("discount"),
		"mtr_discount_percent.is_active":         c.Query("is_active"),
	}

	limit := utils.GetQueryInt(c, "limit")
	page := utils.GetQueryInt(c, "page")
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.discountpercentservice.WithTrx(trxHandle).GetAllDiscountPercent(criteria)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	paginatedData, totalPages, totalRows := utils.DataFramePaginate(result, page, limit, utils.SnaketoPascalCase(sortOf), sortBy)

	payloads.HandleSuccessPagination(c, utils.ModifyKeysInResponse(paginatedData), "success", 200, limit, page, int64(totalRows), totalPages)
}

// @Summary Get Discount Percent By ID
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @Param discount_percent_id path int true "discount_percent_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent/{discount_percent_id} [get]
func (r *DiscountPercentController) GetDiscountPercentByID(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	discountPercentId, _ := strconv.Atoi(c.Param("discount_percent_id"))
	result, err := r.discountpercentservice.WithTrx(trxHandle).GetDiscountPercentById(int(discountPercentId))
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Discount Percent
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @param reqBody body masteritempayloads.DiscountPercentResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent [post]
func (r *DiscountPercentController) SaveDiscountPercent(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteritempayloads.DiscountPercentResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.DiscountPercentId) != 0 {
		result, err := r.discountpercentservice.WithTrx(trxHandle).GetDiscountPercentById(int(request.DiscountPercentId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.DiscountPercentId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.discountpercentservice.WithTrx(trxHandle).SaveDiscountPercent(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.DiscountPercentId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Discount Percent
// @Description REST API Discount Percent
// @Accept json
// @Produce json
// @Tags Master : Discount Percent
// @param discount_percent_id path int true "discount_percent_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-percent/{discount_percent_id} [patch]
func (r *DiscountPercentController) ChangeStatusDiscountPercent(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	discountPercentId, err := strconv.Atoi(c.Param("discount_percent_id"))
	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.discountpercentservice.WithTrx(trxHandle).GetDiscountPercentById(int(discountPercentId))
	if err != nil || result.DiscountPercentId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.discountpercentservice.WithTrx(trxHandle).ChangeStatusDiscountPercent(int(discountPercentId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}