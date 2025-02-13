package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkOrderBypassControllerImpl struct {
	WorkOrderBypassService transactionworkshopservice.WorkOrderBypassService
}

type WorkOrderBypassController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Bypass(writer http.ResponseWriter, request *http.Request)
}

func NewWorkOrderBypassController(WorkOrderBypassService transactionworkshopservice.WorkOrderBypassService) WorkOrderBypassController {
	return &WorkOrderBypassControllerImpl{
		WorkOrderBypassService: WorkOrderBypassService,
	}
}

// GetAll gets all work orders
// @Summary Get All Work Orders
// @Description Retrieve all work orders with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Bypass
// @Param work_order_system_number query string false "Work Order System Number"
// @Param work_order_document_number query string false "Work Order Document Number"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-bypass [get]
func (r *WorkOrderBypassControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.work_order_system_number":   queryValues.Get("work_order_system_number"),
		"trx_work_order.work_order_document_number": queryValues.Get("work_order_document_number"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.WorkOrderBypassService.GetAll(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)

}

// GetById gets a work order by id
// @Summary Get Work Order By Id
// @Description Retrieve a work order by id
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Bypass
// @Param work_order_system_number path int true "Work Order System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-bypass/{work_order_system_number} [get]
func (r *WorkOrderBypassControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	workOrderSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	workOrder, err := r.WorkOrderBypassService.GetById(workOrderSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(workOrder), "Get Data Successfully", http.StatusOK)

}

// Bypass bypasses a work order
// @Summary Bypass Work Order
// @Description Bypass a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Bypass
// @Param work_order_system_number path int true "Work Order System Number"
// @Param body body transactionworkshoppayloads.WorkOrderBypassRequestDetail true "Work Order Bypass Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-bypass/{work_order_system_number}/bypass [post]
func (r *WorkOrderBypassControllerImpl) Bypass(writer http.ResponseWriter, request *http.Request) {

	idstr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(idstr)
	if err != nil {

		payloads.NewHandleError(writer, "Invalid Work Order ID", http.StatusBadRequest)
		return
	}

	var detailRequest transactionworkshoppayloads.WorkOrderBypassRequestDetail
	helper.ReadFromRequestBody(request, &detailRequest)
	if validationErr := validation.ValidationForm(writer, request, &detailRequest); validationErr != nil {
		exceptions.NewConflictException(writer, request, validationErr)
		return
	}

	workOrder, baseErr := r.WorkOrderBypassService.Bypass(workOrderId, detailRequest)
	if baseErr != nil {

		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "id request not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, workOrder, "Bypass Work Order Successfully", http.StatusOK)

}
