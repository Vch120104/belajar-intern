package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
)

type QualityControlControllerImpl struct {
	QualityControlService transactionworkshopservice.QualityControlService
}

type QualityControlController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Qcpass(writer http.ResponseWriter, request *http.Request)
	Reorder(writer http.ResponseWriter, request *http.Request)
}

func NewQualityControlController(QualityControlService transactionworkshopservice.QualityControlService) QualityControlController {
	return &QualityControlControllerImpl{
		QualityControlService: QualityControlService,
	}
}

// GetAll gets all quality controls
// @Summary Get All Quality Controls
// @Description Retrieve all quality controls with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Quality Control
// @Param quality_control_system_number query string false "Quality Control System Number"
// @Param quality_control_document_number query string false "Quality Control Document Number"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/quality-control [get]
func (r *QualityControlControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_customer.customer_name":         queryValues.Get("customer_name"),
		"mtr_unit_model.model_code":          queryValues.Get("model_code"),
		"mtr_unit_variant.varian_code":       queryValues.Get("varian_code"),
		"mtr_vehicle.vehicle_chassis_number": queryValues.Get("vehicle_chassis_number"),
		"mtr_vehicle_registration_certificate.vehicle_registration_certificate_tnkb": queryValues.Get("vehicle_registration_certificate_tnkb"),
		"trx_work_order.work_order_date":                                             queryValues.Get("work_order_date"),
		"trx_work_order.work_order_system_number":                                    queryValues.Get("work_order_system_number"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.QualityControlService.GetAll(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetById gets quality control by id
// @Summary Get Quality Control By Id
// @Description Retrieve quality control by id
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Quality Control
// @Param id path string true "Quality Control Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/quality-control/{work_order_system_number} [get]
func (r *QualityControlControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	idstr := chi.URLParam(request, "work_order_system_number")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid request ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	// Convert map to []utils.FilterCondition
	var filterConditions []utils.FilterCondition
	for field, value := range map[string]string{
		"trx_work_order.work_order_system_number": queryValues.Get("work_order_system_number"),
	} {
		if value != "" {
			filterConditions = append(filterConditions, utils.FilterCondition{
				ColumnField: field,
				ColumnValue: value,
			})
		}
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	data, baseErr := r.QualityControlService.GetById(id, filterConditions, paginate)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Get Data Successfully", http.StatusOK)
}

// Bypass bypasses quality control
// @Summary Bypass Quality Control
// @Description Bypass quality control
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Quality Control
// @Param id path string true "Quality Control Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/quality-control/{work_order_system_number}/{work_order_detail_id}/qcpass [put]
func (r *QualityControlControllerImpl) Qcpass(writer http.ResponseWriter, request *http.Request) {

	idstr := chi.URLParam(request, "work_order_system_number")
	id, err := strconv.Atoi(idstr)
	if err != nil {

		payloads.NewHandleError(writer, "Invalid Work Order ID", http.StatusBadRequest)
		return
	}

	iddetstr := chi.URLParam(request, "work_order_detail_id")
	iddet, err := strconv.Atoi(iddetstr)
	if err != nil {

		payloads.NewHandleError(writer, "Invalid Work Order Detail ID", http.StatusBadRequest)
		return
	}

	bypass, baseErr := r.QualityControlService.Qcpass(id, iddet)
	if baseErr != nil {

		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "id request not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, bypass, "Quality Control Qcpassed Successfully", http.StatusOK)
}

// Reorder reorders quality control
// @Summary Reorder Quality Control
// @Description Reorder quality control
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Quality Control
// @Param id path string true "Quality Control Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/quality-control/{work_order_system_number}/{work_order_detail_id}/reorder [put]
func (r *QualityControlControllerImpl) Reorder(writer http.ResponseWriter, request *http.Request) {
	idstr := chi.URLParam(request, "work_order_system_number")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Work Order ID", http.StatusBadRequest)
		return
	}

	iddetstr := chi.URLParam(request, "work_order_detail_id")
	iddet, err := strconv.Atoi(iddetstr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Work Order Detail ID", http.StatusBadRequest)
		return
	}

	var ReorderRequest transactionworkshoppayloads.QualityControlReorder
	helper.ReadFromRequestBody(request, &ReorderRequest)

	reorder, baseErr := r.QualityControlService.Reorder(id, iddet, ReorderRequest)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "id request not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, reorder, "Quality Control Reordered Successfully", http.StatusOK)
}
