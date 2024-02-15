package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/middlewares"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UnitOfMeasurementController struct {
	unitofmeasurementservice masteritemservice.UnitOfMeasurementService
}

func StartUnitOfMeasurementRoutes(
	db *gorm.DB,
	r *gin.RouterGroup,
	unitofmeasurementservice masteritemservice.UnitOfMeasurementService,
) {
	unitOfMeasurementHandler := UnitOfMeasurementController{unitofmeasurementservice: unitofmeasurementservice}
	r.GET("/unit-of-measurement", middlewares.DBTransactionMiddleware(db), unitOfMeasurementHandler.GetAllUnitOfMeasurement)
	r.GET("/unit-of-measurement-drop-down", middlewares.DBTransactionMiddleware(db), unitOfMeasurementHandler.GetAllUnitOfMeasurementIsActive)
	r.GET("/unit-of-measurement-by-code/:uom_code", middlewares.DBTransactionMiddleware(db), unitOfMeasurementHandler.GetUnitOfMeasurementByCode)
	r.POST("/unit-of-measurement", middlewares.DBTransactionMiddleware(db), unitOfMeasurementHandler.SaveUnitOfMeasurement)
	r.PATCH("/unit-of-measurement/:uom_id", middlewares.DBTransactionMiddleware(db), unitOfMeasurementHandler.ChangeStatusUnitOfMeasurement)
}

// @Summary Get All Unit Of Measurement
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param uom_code query string false "uom_code"
// @Param uom_description query string false "uom_description"
// @Param uom_type_desc query string false "uom_type_desc"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement [get]
func (r *UnitOfMeasurementController) GetAllUnitOfMeasurement() {

	queryParams := map[string]string{
		"mtr_uom.is_active":          c.Query("is_active"),
		"mtr_uom.uom_code":           c.Query("uom_code"),
		"mtr_uom.uom_description":    c.Query("uom_description"),
		"mtr_uom_type.uom_type_desc": c.Query("uom_type_desc"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.GetQueryInt(c, "limit"),
		Page:   utils.GetQueryInt(c, "page"),
		SortOf: c.Query("sort_of"),
		SortBy: c.Query("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.unitofmeasurementservice.GetAllUnitOfMeasurement(filterCondition, pagination)

	if result.Rows == nil {
		exceptions.NotFoundException(c, "Nothing matching request")
		return
	}

	payloads.HandleSuccessPagination(c, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get All Unit Of Measurement drop down
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement-drop-down [get]
func (r *UnitOfMeasurementController) GetAllUnitOfMeasurementIsActive() {

	result := r.unitofmeasurementservice.GetAllUnitOfMeasurementIsActive()

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Unit Of Measurement By Code
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @Param uom_code path string true "uom_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement-by-code/{uom_code} [get]
func (r *UnitOfMeasurementController) GetUnitOfMeasurementByCode() {

	operationGroupCode := c.Param("uom_code")
	result := r.unitofmeasurementservice.GetUnitOfMeasurementByCode(operationGroupCode)

	payloads.HandleSuccess(c, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Unit Of Measurement
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @param reqBody body masteritempayloads.UomResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement [post]
func (r *UnitOfMeasurementController) SaveUnitOfMeasurement() {

	var request masteritempayloads.UomResponse
	var message = ""

	if err := c.ShouldBindJSON(&request); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	create := r.unitofmeasurementservice.SaveUnitOfMeasurement(request)

	if request.UomId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.HandleSuccess(c, create, message, http.StatusOK)
}

// @Summary Change Status Unit Of Measurement
// @Description REST API Unit Of Measurement
// @Accept json
// @Produce json
// @Tags Master : Unit Of Measurement
// @param uom_id path int true "uom_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/unit-of-measurement/{uom_id} [patch]
func (r *UnitOfMeasurementController) ChangeStatusUnitOfMeasurement() {

	uomId, _ := strconv.Atoi(c.Param("uom_id"))

	response := r.unitofmeasurementservice.ChangeStatusUnitOfMeasurement(int(uomId))

	payloads.HandleSuccess(c, response, "Update Data Successfully!", http.StatusOK)
}
