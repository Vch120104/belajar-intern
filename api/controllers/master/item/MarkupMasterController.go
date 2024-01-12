package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	"net/http"
	"strconv"

	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"

	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	// "strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MarkupMasterController struct {
	markupMasterService masteritemservice.MarkupMasterService
}

func StartMarkupMasterRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	markupMasterService masteritemservice.MarkupMasterService,
) {
	markupMasterHandler := MarkupMasterController{markupMasterService: markupMasterService}
	r.GET("/markup-master", middlewares.DBTransactionMiddleware(db), markupMasterHandler.GetMarkupMasterList)
	r.GET("/markup-master-by-code/:markup_master_code", middlewares.DBTransactionMiddleware(db), markupMasterHandler.GetMarkupMasterByCode)
	r.POST("/markup-master", middlewares.DBTransactionMiddleware(db), markupMasterHandler.SaveMarkupMaster)
	r.PATCH("/markup-master/:markup_master_id", middlewares.DBTransactionMiddleware(db), markupMasterHandler.ChangeStatusMarkupMaster)
}

// @Summary Get All Markup Master
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param markup_master_code query string false "markup_master_code"
// @Param markup_master_description query string false "markup_master_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master [get]
func (r *MarkupMasterController) GetMarkupMasterList(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	queryParams := map[string]string{
		"markup_master_code":        c.Query("markup_master_code"),
		"markup_master_description": c.Query("markup_master_description"),
		"is_active":                 c.Query("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.GetQueryInt(c, "limit"),
		Page:   utils.GetQueryInt(c, "page"),
		SortOf: c.Query("sort_of"),
		SortBy: c.Query("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.markupMasterService.WithTrx(trxHandle).GetMarkupMasterList(filterCondition, pagination)

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

// @Summary Get Markup Master Description by code
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @Param markup_master_code path string true "markup_master_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master-by-code/{markup_master_code} [get]
func (r *MarkupMasterController) GetMarkupMasterByCode(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	markupMasterCode := c.Param("markup_master_code")
	result, err := r.markupMasterService.WithTrx(trxHandle).GetMarkupMasterByCode(markupMasterCode)
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Markup Master
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @param reqBody body masteritempayloads.MarkupMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master [post]
func (r *MarkupMasterController) SaveMarkupMaster(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteritempayloads.MarkupMasterResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.MarkupMasterId) != 0 {
		result, err := r.markupMasterService.WithTrx(trxHandle).GetMarkupMasterById(int(request.MarkupMasterId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.MarkupMasterId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.markupMasterService.WithTrx(trxHandle).SaveMarkupMaster(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.MarkupMasterId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Markup Master
// @Description REST API Markup Master
// @Accept json
// @Produce json
// @Tags Master : Markup Master
// @param markup_master_id path int true "markup_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-master/{markup_master_id} [patch]
func (r *MarkupMasterController) ChangeStatusMarkupMaster(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	markupMasterId, err := strconv.Atoi(c.Param("markup_master_id"))
	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.markupMasterService.WithTrx(trxHandle).GetMarkupMasterById(int(markupMasterId))
	if err != nil || result.MarkupMasterId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.markupMasterService.WithTrx(trxHandle).ChangeStatusMasterMarkupMaster(int(markupMasterId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}
