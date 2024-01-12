package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/utils"

	// "after-sales/api/middlewares"

	masteroperationservice "after-sales/api/services/master/operation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OperationEntriesController struct {
	operationentriesservice masteroperationservice.OperationEntriesService
}

func StartOperationEntriesRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	operationentriesservice masteroperationservice.OperationEntriesService,
) {
	operationEntriesHandler := OperationEntriesController{operationentriesservice: operationentriesservice}
	r.GET("/operation-entries/:operation_entries_id", middlewares.DBTransactionMiddleware(db), operationEntriesHandler.GetOperationEntriesByID)
	r.GET("/operation-entries-by-name", middlewares.DBTransactionMiddleware(db), operationEntriesHandler.GetOperationEntriesName)
	r.POST("/operation-entries", middlewares.DBTransactionMiddleware(db), operationEntriesHandler.SaveOperationEntries)
	r.PATCH("/operation-entries/:operation_entries_id", middlewares.DBTransactionMiddleware(db), operationEntriesHandler.ChangeStatusOperationEntries)
}

// @Summary Get Operation Entries By ID
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @Param operation_entries_id path int true "operation_entries_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-entries/{operation_entries_id} [get]
func (r *OperationEntriesController) GetOperationEntriesByID(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationId, _ := strconv.Atoi(c.Param("operation_entries_id"))
	result, err := r.operationentriesservice.WithTrx(trxHandle).GetOperationEntriesById(int32(operationId))
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Entries Name
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @Param operation_group_id query int true "operation_group_id"
// @Param operation_section_id query int true "operation_section_id"
// @Param operation_key_id query int true "operation_key_id"
// @Param operation_entries_code query string true "operation_entries_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-entries-by-name [get]
func (r *OperationEntriesController) GetOperationEntriesName(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationGroupId := utils.GetQueryInt(c, "operation_group_id")
	operationSectionId := utils.GetQueryInt(c, "operation_section_id")
	operationKeyId := utils.GetQueryInt(c, "operation_key_id")
	operationEntriesCode := c.Query("operation_entries_code")

	result, err := r.operationentriesservice.WithTrx(trxHandle).GetOperationEntriesName(masteroperationpayloads.OperationEntriesRequest{
		OperationGroupId:     operationGroupId,
		OperationSectionId:   operationSectionId,
		OperationKeyId:       operationKeyId,
		OperationEntriesCode: operationEntriesCode,
	})

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Entries
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @param reqBody body masteroperationpayloads.OperationEntriesResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-entries [post]
func (r *OperationEntriesController) SaveOperationEntries(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteroperationpayloads.OperationEntriesResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.OperationEntriesId) != 0 {
		result, err := r.operationentriesservice.WithTrx(trxHandle).GetOperationEntriesById(request.OperationEntriesId)

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.OperationEntriesId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.operationentriesservice.WithTrx(trxHandle).SaveOperationEntries(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.OperationEntriesId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Operation Entries
// @Description REST API Operation Entries
// @Accept json
// @Produce json
// @Tags Master : Operation Entries
// @param operation_entries_id path int true "operation_entries_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-entries/{operation_entries_id} [patch]
func (r *OperationEntriesController) ChangeStatusOperationEntries(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationEntriesId, err := strconv.Atoi(c.Param("operation_entries_id"))

	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.operationentriesservice.WithTrx(trxHandle).GetOperationEntriesById(int32(operationEntriesId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if result.OperationEntriesId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.operationentriesservice.WithTrx(trxHandle).ChangeStatusOperationEntries(operationEntriesId)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Change Status Successfully!", http.StatusOK)
}
