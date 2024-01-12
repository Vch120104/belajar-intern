package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	masteroperationservice "after-sales/api/services/master/operation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OperationKeyController struct {
	operationkeyservice masteroperationservice.OperationKeyService
}

func StartOperationKeyRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	operationkeyservice masteroperationservice.OperationKeyService,
) {
	operationKeyHandler := OperationKeyController{operationkeyservice: operationkeyservice}
	r.GET("/operation-key/:operation_key_id", middlewares.DBTransactionMiddleware(db), operationKeyHandler.GetOperationKeyByID)
	r.GET("/operation-key", middlewares.DBTransactionMiddleware(db), operationKeyHandler.GetAllOperationKeyList)
	r.GET("/operation-key-name", middlewares.DBTransactionMiddleware(db), operationKeyHandler.GetOperationKeyName)
	r.POST("/operation-key", middlewares.DBTransactionMiddleware(db), operationKeyHandler.SaveOperationKey)
	r.PATCH("/operation-key/:operation_key_id", middlewares.DBTransactionMiddleware(db), operationKeyHandler.ChangeStatusOperationKey)
}

// @Summary Get All Operation Key
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param operation_section_code query string false "operation_section_code"
// @Param operation_section_description query string false "operation_section_description"
// @Param operation_group_code query string false "operation_group_code"
// @Param operation_group_description query string false "operation_group_description"
// @Param operation_key_code query string false "operation_key_code"
// @Param operation_key_description query string false "operation_key_description"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-key [get]
func (r *OperationKeyController) GetAllOperationKeyList(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	queryParams := map[string]string{
		"mtr_operation_group.operation_group_code":            c.Query("operation_group_code"),
		"mtr_operation_group.operation_group_description":     c.Query("operation_group_description"),
		"mtr_operation_section.operation_section_code":        c.Query("operation_section_code"),
		"mtr_operation_section.operation_section_description": c.Query("operation_section_description"),
		"mtr_operation_key.is_active":                         c.Query("is_active"),
		"mtr_operation_key.operation_key_code":                c.Query("operation_key_code"),
		"mtr_operation_key.operation_key_description":         c.Query("operation_key_description"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.GetQueryInt(c, "limit"),
		Page:   utils.GetQueryInt(c, "page"),
		SortOf: c.Query("sort_of"),
		SortBy: c.Query("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.operationkeyservice.WithTrx(trxHandle).GetAllOperationKeyList(criteria, pagination)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccessPagination(c, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Operation Key By ID
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Param operation_key_id path int true "operation_key_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-key/{operation_key_id} [get]
func (r *OperationKeyController) GetOperationKeyByID(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationKeyId, _ := strconv.Atoi(c.Param("operation_key_id"))
	result, err := r.operationkeyservice.WithTrx(trxHandle).GetOperationKeyById(operationKeyId)
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Key Name
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @Param operation_group_id query int true "operation_group_id"
// @Param operation_section_id query int true "operation_section_id"
// @Param operation_key_code query string true "operation_key_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-key-name [get]
func (r *OperationKeyController) GetOperationKeyName(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationGroupId := utils.GetQueryInt(c, "operation_group_id")
	operationSectionId := utils.GetQueryInt(c, "operation_section_id")
	keyCode := c.Query("operation_key_code")

	result, err := r.operationkeyservice.WithTrx(trxHandle).GetOperationKeyName(masteroperationpayloads.OperationKeyRequest{
		OperationGroupId:   int32(operationGroupId),
		OperationSectionId: int32(operationSectionId),
		OperationKeyCode:   keyCode,
	})

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Key
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @param reqBody body masteroperationpayloads.OperationKeyResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-key [post]
func (r *OperationKeyController) SaveOperationKey(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteroperationpayloads.OperationKeyResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.OperationKeyId) != 0 {
		result, err := r.operationkeyservice.WithTrx(trxHandle).GetOperationKeyById(int(request.OperationKeyId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.OperationKeyId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.operationkeyservice.WithTrx(trxHandle).SaveOperationKey(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.OperationKeyId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Operation Key
// @Description REST API Operation Key
// @Accept json
// @Produce json
// @Tags Master : Operation Key
// @param operation_key_id path int true "operation_key_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-key/{operation_key_id} [patch]
func (r *OperationKeyController) ChangeStatusOperationKey(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationKeyId, err := strconv.Atoi(c.Param("operation_key_id"))

	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.operationkeyservice.WithTrx(trxHandle).GetOperationKeyById(int(operationKeyId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if result.OperationKeyId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.operationkeyservice.WithTrx(trxHandle).ChangeStatusOperationKey(int(operationKeyId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Change Status Successfully!", http.StatusOK)
}
