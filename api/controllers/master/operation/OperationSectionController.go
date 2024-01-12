package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OperationSectionController struct {
	operationsectionservice masteroperationservice.OperationSectionService
}

func StartOperationSectionRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	operationsectionservice masteroperationservice.OperationSectionService,
) {
	operationSectionHandler := OperationSectionController{operationsectionservice: operationsectionservice}
	r.GET("/operation-section", middlewares.DBTransactionMiddleware(db), operationSectionHandler.GetAllOperationSectionList)
	r.GET("/operation-section/:operation_section_id", middlewares.DBTransactionMiddleware(db), operationSectionHandler.GetOperationSectionByID)
	r.GET("/operation-section-name", middlewares.DBTransactionMiddleware(db), operationSectionHandler.GetOperationSectionName)
	r.GET("/operation-section-code-by-group-id", middlewares.DBTransactionMiddleware(db), operationSectionHandler.GetSectionCodeByGroupId)
	r.PUT("/operation-section", middlewares.DBTransactionMiddleware(db), operationSectionHandler.SaveOperationSection)
	r.PATCH("/operation-section/:operation_section_id", middlewares.DBTransactionMiddleware(db), operationSectionHandler.ChangeStatusOperationSection)
}

// @Summary Get All Operation Section
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_section_code query string false "operation_section_code"
// @Param operation_section_description query string false "operation_section_description"
// @Param operation_group_code query string false "operation_group_code"
// @Param operation_group_description query string false "operation_group_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-section [get]
func (r *OperationSectionController) GetAllOperationSectionList(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	queryParams := map[string]string{
		"mtr_operation_group.operation_group_code":            c.Query("operation_group_code"),
		"mtr_operation_group.operation_group_description":     c.Query("operation_group_description"),
		"mtr_operation_section.peration_section_code":         c.Query("operation_section_code"),
		"mtr_operation_section.operation_section_description": c.Query("operation_section_description"),
		"mtr_operation_section.is_active":                     c.Query("is_active"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.GetQueryInt(c, "limit"),
		Page:   utils.GetQueryInt(c, "page"),
		SortOf: c.Query("sort_of"),
		SortBy: c.Query("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.operationsectionservice.WithTrx(trxHandle).GetAllOperationSectionList(criteria, pagination)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	if result.Rows == nil {
		exceptions.NotFoundException(c, "Nothing matching request")
		return
	}

	payloads.HandleSuccessPagination(c, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Operation Section By ID
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param operation_section_id path int true "operation_section_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-section/{operation_section_id} [get]
func (r *OperationSectionController) GetOperationSectionByID(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationSectionId, _ := strconv.Atoi(c.Param("operation_section_id"))

	result, err := r.operationsectionservice.WithTrx(trxHandle).GetOperationSectionById(int(operationSectionId))

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Section Code By Group Id
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param operation_group_id query int true "operation_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-section-code-by-group-id [get]
func (r *OperationSectionController) GetSectionCodeByGroupId(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	groupId := c.Query("operation_group_id")

	result, err := r.operationsectionservice.WithTrx(trxHandle).GetSectionCodeByGroupId(groupId)
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Section Name
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @Param operation_group_id query int true "operation_group_id"
// @Param operation_section_code query string true "operation_section_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-section-name [get]
func (r *OperationSectionController) GetOperationSectionName(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationGroupId := utils.GetQueryInt(c, "operation_group_id")
	section_code := c.Query("operation_section_code")

	result, err := r.operationsectionservice.WithTrx(trxHandle).GetOperationSectionName(operationGroupId, section_code)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Section
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @param reqBody body masteroperationpayloads.OperationSectionRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-section [put]
func (r *OperationSectionController) SaveOperationSection(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteroperationpayloads.OperationSectionRequest
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	create, err := r.operationsectionservice.WithTrx(trxHandle).SaveOperationSection(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.OperationSectionId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Operation Section
// @Description REST API Operation Section
// @Accept json
// @Produce json
// @Tags Master : Operation Section
// @param operation_section_id path int true "operation_section_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-section/{operation_section_id} [patch]
func (r *OperationSectionController) ChangeStatusOperationSection(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationSectionId, err := strconv.Atoi(c.Param("operation_section_id"))

	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.operationsectionservice.WithTrx(trxHandle).GetOperationSectionById(int(operationSectionId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if result.OperationSectionId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.operationsectionservice.WithTrx(trxHandle).ChangeStatusOperationSection(int(operationSectionId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Change Status Successfully!", http.StatusOK)
}
