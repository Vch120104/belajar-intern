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
func (r *MarkupMasterController) GetMarkupMasterList() {

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

	result := r.markupMasterService.GetMarkupMasterList(filterCondition, pagination)




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
func (r *MarkupMasterController) GetMarkupMasterByCode() {

	markupMasterCode := c.Param("markup_master_code")
	result := r.markupMasterService.GetMarkupMasterByCode(markupMasterCode)

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
func (r *MarkupMasterController) SaveMarkupMaster() {

	var request masteritempayloads.MarkupMasterResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c.Error())
		return
	}


	create := r.markupMasterService.SaveMarkupMaster(request)

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
func (r *MarkupMasterController) ChangeStatusMarkupMaster() {

	markupMasterId, _ := strconv.Atoi(c.Param("markup_master_id"))


	response := r.markupMasterService.ChangeStatusMasterMarkupMaster(int(markupMasterId))

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}
