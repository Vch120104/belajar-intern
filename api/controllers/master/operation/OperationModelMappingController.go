package masteroperationcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	masteroperationpayloads "after-sales/api/payloads/master/operation"

	// "after-sales/api/middlewares"

	masteroperationservice "after-sales/api/services/master/operation"
)

type OperationModelMappingController struct {
	operationmodelmappingservice masteroperationservice.OperationModelMappingService
}

func StartOperationModelMappingRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	operationmodelmappingservice masteroperationservice.OperationModelMappingService,
) {
	operationModelMappingHandler := OperationModelMappingController{operationmodelmappingservice: operationmodelmappingservice}
	r.GET("/operation-model-mapping/", middlewares.DBTransactionMiddleware(db), operationModelMappingHandler.GetOperationModelMappingLookup)
	r.GET("/operation-model-mapping/:operation_model_mapping_id", middlewares.DBTransactionMiddleware(db), operationModelMappingHandler.GetOperationModelMappingById)
	r.GET("/operation-model-mapping-by-brand-model-operation-id/", middlewares.DBTransactionMiddleware(db), operationModelMappingHandler.GetOperationModelMappingByBrandModelOperationCode)
	r.POST("/operation-model-mapping/", middlewares.DBTransactionMiddleware(db), operationModelMappingHandler.SaveOperationModelMapping)
	r.PATCH("/operation-model-mapping/:operation_model_mapping_id", middlewares.DBTransactionMiddleware(db), operationModelMappingHandler.ChangeStatusOperationModelMapping)
}

// GetOperationModelMappingById(int) (masteroperationpayloads.OperationModelMappingResponse, error)
// GetOperationModelMappingLookup(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
// GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, error)
// SaveOperationModelMapping(masteroperationpayloads.OperationModelMappingResponse) (bool, error)
// ChangeStatusOperationModelMapping(int) (bool, error)

// @Summary Get All Operation Model Mapping Lookup
// @Description REST API Operation Model Mapping Lookup
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active"
// @Param operation_group_code query string false "operation_group_code"
// @Param operation_name query string false "operation_name"
// @Param operation_code query string false "operation_code"
// @Param brand_name query string false "brand_name"
// @Param model_code query string false "model_code"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-model-mapping [get]
func (r *OperationModelMappingController) GetOperationModelMappingLookup(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	queryParams := map[string]string{
		"mtr_operation_model_mapping.is_active":            c.Query("is_active"),
		"mtr_operation_model_mapping.operation_group_code": c.Query("operation_group_code"),
		"mtr_operation_code.operation_name":                c.Query("operation_name"),
		"mtr_operation_model_mapping.operation_code":       c.Query("operation_code"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingLookup(criteria)

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "success", 200)
}

// @Summary Get Operation Model Mapping By Id
// @Description REST API Operation Model Mapping
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param operation_model_mapping_id path string true "operation_model_mapping_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-model-mapping/{operation_model_mapping_id} [get]
func (r *OperationModelMappingController) GetOperationModelMappingById(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationModelMappingId, _ := strconv.Atoi(c.Param("operation_model_mapping_id"))
	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingById(operationModelMappingId)
	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}
	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Operation Model Mapping By Brand Model OperationId
// @Description REST API Operation Model Mapping
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @Param brand_id query string true "brand_id"
// @Param model_id query string true "model_id"
// @Param operation_id query string true "operation_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-model-mapping-by-brand-model-operation-id/ [get]
func (r *OperationModelMappingController) GetOperationModelMappingByBrandModelOperationCode(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	brandId, _ := strconv.Atoi(c.Query("brand_id"))
	modelId, _ := strconv.Atoi(c.Query("model_id"))
	operationId, _ := strconv.Atoi(c.Query("operation_id"))

	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingByBrandModelOperationCode(masteroperationpayloads.OperationModelModelBrandOperationCodeRequest{
		BrandId:     brandId,
		ModelId:     modelId,
		OperationId: operationId,
	})

	if err != nil {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Operation Model Mapping
// @Description REST API Operation Model Mapping
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @param reqBody body masteroperationpayloads.OperationModelMappingResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-model-mapping [post]
func (r *OperationModelMappingController) SaveOperationModelMapping(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	var request masteroperationpayloads.OperationModelMappingResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	if int(request.OperationModelMappingId) != 0 {
		result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingById(int(request.OperationModelMappingId))

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		if result.OperationModelMappingId == 0 {
			exceptions.NotFoundException(c, err.Error())
			return
		}
	}

	create, err := r.operationmodelmappingservice.WithTrx(trxHandle).SaveOperationModelMapping(request)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	if request.OperationModelMappingId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Operation Model Mapping
// @Description REST API Operation Model Mapping
// @Accept json
// @Produce json
// @Tags Master : Operation Model Mapping
// @param operation_model_mapping_id path string true "operation_model_mapping_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/operation-model-mapping/{operation_model_mapping_id} [patch]
func (r *OperationModelMappingController) ChangeStatusOperationModelMapping(c *gin.Context) {
	trxHandle := c.MustGet("db_trx").(*gorm.DB)
	operationModelMappingId, err := strconv.Atoi(c.Param("operation_model_mapping_id"))

	if err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}
	//id check
	result, err := r.operationmodelmappingservice.WithTrx(trxHandle).GetOperationModelMappingById(int(operationModelMappingId))
	if err != nil || result.OperationId == 0 {
		exceptions.NotFoundException(c, err.Error())
		return
	}

	response, err := r.operationmodelmappingservice.WithTrx(trxHandle).ChangeStatusOperationModelMapping(int(operationModelMappingId))
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, response, "Change Status Successfully!", http.StatusOK)
}
