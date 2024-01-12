package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"

	// "after-sales/api/middlewares"

	masteroperationservice "after-sales/api/services/master/operation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OperationCodeController struct {
	operationcodeservice masteroperationservice.OperationCodeService
}

func StartOperationCodeRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	operationcodeservice masteroperationservice.OperationCodeService,
) {
	operationCodeHandler := OperationCodeController{operationcodeservice: operationcodeservice}
	r.GET("/operation-code/:operation_id", middlewares.DBTransactionMiddleware(db), operationCodeHandler.GetOperationCodeByID)
}

// // @Summary Get All Operation Code
// // @Description REST API Operation Code
// // @Accept json
// // @Produce json
// // @Tags Operation Code
// // @Param page query string true "page"
// // @Param limit query string true "limit"
// // @Param operation_group_code query string false "operation_group_code"
// // @Param operation_group_description query string false "operation_group_description"
// // @Param is_active query string false "is_active" Enums(true, false)
// // @Param sort_by query string false "sort_by"
// // @Param sort_of query string false "sort_of"
// // @Success 200 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.Error
// // @Router /aftersales-service/api/aftersales/operation-code [get]
// func (r *OperationGroupController) GetAllOperationCode(c *gin.Context) {
// 	trxHandle := c.MustGet("db_trx").(*gorm.DB)
// 	queryParams := map[string]string{
// 		"operation_group_code":        c.Query("operation_group_code"),
// 		"operation_group_description": c.Query("operation_group_description"),
// 		"is_active":                   c.Query("is_active"),
// 	}

// 	pagination := pagination.Pagination{
// 		Limit:  utils.GetQueryInt(c, "limit"),
// 		Page:   utils.GetQueryInt(c, "page"),
// 		SortOf: c.Query("sort_of"),
// 		SortBy: c.Query("sort_by"),
// 	}

// 	filterCondition := utils.BuildFilterCondition(queryParams)

// 	result, err := r.operationgroupservice.WithTrx(trxHandle).GetAllOperationGroup(filterCondition, pagination)

// 	if err != nil {
// 		exceptions.AppException(c, err.Error())
// 		return
// 	}

// 	if result.Rows == nil {
// 		exceptions.NotFoundException(c, "Nothing matching request")
// 		return
// 	}

// 	payloads.HandleSuccessPagination(c, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
// }

// @Summary Get Operation Code By ID
// @Description REST API Operation Code
// @Accept json
// @Produce json
// @Tags Master : Operation Code
// @Param operation_id path int true "operation_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-code/{operation_id} [get]
func (r *OperationCodeController) GetOperationCodeByID(c *gin.Context) {
	operationId, _ := strconv.Atoi(c.Param("operation_id"))
	result, err := r.operationcodeservice.GetOperationCodeById(int32(operationId))
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}
