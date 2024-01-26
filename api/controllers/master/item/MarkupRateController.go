package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	masteritemservice "after-sales/api/services/master/item"

	// "after-sales/api/middlewares"

	// "strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MarkupRateController struct {
	markupRateService masteritemservice.MarkupRateService
}

func StartMarkupRateRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	markupRateService masteritemservice.MarkupRateService,
) {
	markupRateHandler := MarkupRateController{markupRateService: markupRateService}
	r.GET("/markup-rate", middlewares.DBTransactionMiddleware(db), markupRateHandler.GetAllMarkupRate)
	r.GET("/markup-rate/:markup_rate_id", middlewares.DBTransactionMiddleware(db), markupRateHandler.GetMarkupRateByID)
	r.POST("/markup-rate", middlewares.DBTransactionMiddleware(db), markupRateHandler.SaveMarkupRate)
	r.PATCH("/markup-rate/:markup_rate_id", middlewares.DBTransactionMiddleware(db), markupRateHandler.ChangeStatusMarkupRate)
}

// @Summary Get All Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param markup_master_code query string false "markup_master_code"
// @Param markup_master_description query string false "markup_master_description"
// @Param order_type_name query string false "order_type_name"
// @Param markup_rate query float64 false "markup_rate"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate [get]
func (r *MarkupRateController) GetAllMarkupRate(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)

	queryParams := map[string]string{
		"mtr_markup_master.markup_master_code":        c.Query("markup_master_code"),
		"mtr_markup_master.markup_master_description": c.Query("markup_master_description"),
		"order_type_name":                             c.Query("order_type_name"),
		"mtr_markup_rate.markup_rate":                 c.Query("markup_rate"),
		"mtr_markup_rate.is_active":                   c.Query("is_active"),
	}

	limit := utils.GetQueryInt(c, "limit")
	page := utils.GetQueryInt(c, "page")
	sortOf := c.Query("sort_of")
	sortBy := c.Query("sort_by")

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.markupRateService.WithTrx(trxHandle).GetAllMarkupRate(criteria)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	paginatedData, totalPages, totalRows := utils.DataFramePaginate(result, page, limit, utils.SnaketoPascalCase(sortOf), sortBy)

	payloads.HandleSuccessPagination(c, utils.ModifyKeysInResponse(paginatedData), "success", 200, limit, page, int64(totalRows), totalPages)
}

// @Summary Get Markup Rate By ID
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Param markup_rate_id path int true "markup_rate_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate/{markup_rate_id} [get]
func (r *MarkupRateController) GetMarkupRateByID(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	markupRateId, _ := strconv.Atoi(c.Param("markup_rate_id"))
	result, err := r.markupRateService.WithTrx(trxHandle).GetMarkupRateById(int(markupRateId))
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @param reqBody body masteritempayloads.MarkupRateRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate [post]
func (r *MarkupRateController) SaveMarkupRate(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteritempayloads.MarkupRateRequest
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.MarkupRateId) != 0 {
		result, err := r.markupRateService.WithTrx(trxHandle).GetMarkupRateById(int(request.MarkupRateId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.MarkupRateId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.markupRateService.WithTrx(trxHandle).SaveMarkupRate(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.MarkupRateId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @param markup_rate_id path int true "markup_rate_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate/{markup_rate_id} [patch]
func (r *MarkupRateController) ChangeStatusMarkupRate(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	markupMasterId, err := strconv.Atoi(c.Param("markup_rate_id"))
	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.markupRateService.WithTrx(trxHandle).GetMarkupRateById(int(markupMasterId))
	if err != nil || result.MarkupMasterId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.markupRateService.WithTrx(trxHandle).ChangeStatusMarkupRate(int(markupMasterId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}