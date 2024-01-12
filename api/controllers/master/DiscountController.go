package mastercontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DiscountController struct {
	discountservice masterservice.DiscountService
}

func StartDiscountRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	discountservice masterservice.DiscountService,
) {
	discountHandler := DiscountController{discountservice: discountservice}
	r.GET("/discount", middlewares.DBTransactionMiddleware(db), discountHandler.GetAllDiscount)
	r.GET("/discount-drop-down/", middlewares.DBTransactionMiddleware(db), discountHandler.GetAllDiscountIsActive)
	r.GET("/discount-by-code/:discount_code", middlewares.DBTransactionMiddleware(db), discountHandler.GetDiscountByCode)
	r.POST("/discount", middlewares.DBTransactionMiddleware(db), discountHandler.SaveDiscount)
	r.PATCH("/discount/:discount_code_id", middlewares.DBTransactionMiddleware(db), discountHandler.ChangeStatusDiscount)
}

// @Summary Get All Discount
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param discount_code_value query string false "discount_code_value"
// @Param discount_code_description query string false "discount_code_description"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount [get]
func (r *DiscountController) GetAllDiscount(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	queryParams := map[string]string{
		"is_active":                 c.Query("is_active"),
		"discount_code_value":       c.Query("discount_code_value"),
		"discount_code_description": c.Query("discount_code_description"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.GetQueryInt(c, "limit"),
		Page:   utils.GetQueryInt(c, "page"),
		SortOf: c.Query("sort_of"),
		SortBy: c.Query("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.discountservice.WithTrx(trxHandle).GetAllDiscount(filterCondition, pagination)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if result.Rows == nil {
		exceptions.NotFoundException(c, "Nothing matching request")
		return
	}

	payloads.HandleSuccessPagination(c, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Discount drop down
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-drop-down/ [get]
func (r *DiscountController) GetAllDiscountIsActive(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	result, err := r.discountservice.WithTrx(trxHandle).GetAllDiscountIsActive()
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Discount By Code
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @Param discount_code path string true "discount_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount-by-code/{discount_code} [get]
func (r *DiscountController) GetDiscountByCode(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationGroupCode := c.Param("discount_code")
	result, err := r.discountservice.WithTrx(trxHandle).GetDiscountByCode(operationGroupCode)
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Discount
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @param reqBody body masterpayloads.DiscountResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount [post]
func (r *DiscountController) SaveDiscount(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masterpayloads.DiscountResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.DiscountCodeId) != 0 {
		result, err := r.discountservice.WithTrx(trxHandle).GetDiscountById(int(request.DiscountCodeId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.DiscountCodeId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.discountservice.WithTrx(trxHandle).SaveDiscount(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.DiscountCodeId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Discount
// @Description REST API Discount
// @Accept json
// @Produce json
// @Tags Master : Discount
// @param discount_code_id path int true "discount_code_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/discount/{discount_code_id} [patch]
func (r *DiscountController) ChangeStatusDiscount(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	uomId, err := strconv.Atoi(c.Param("discount_code_id"))
	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.discountservice.WithTrx(trxHandle).GetDiscountById(int(uomId))
	if err != nil || result.DiscountCodeId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.discountservice.WithTrx(trxHandle).ChangeStatusDiscount(int(uomId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}
