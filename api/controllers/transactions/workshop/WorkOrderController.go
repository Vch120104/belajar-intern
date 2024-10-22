package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	utils "after-sales/api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type WorkOrderControllerImpl struct {
	WorkOrderService transactionworkshopservice.WorkOrderService
}

type WorkOrderController interface {
	GetAllRequest(writer http.ResponseWriter, request *http.Request)
	GetRequestById(writer http.ResponseWriter, request *http.Request)
	UpdateRequest(writer http.ResponseWriter, request *http.Request)
	AddRequest(writer http.ResponseWriter, request *http.Request)
	AddRequestMultiId(writer http.ResponseWriter, request *http.Request)
	DeleteRequest(writer http.ResponseWriter, request *http.Request)
	DeleteRequestMultiId(writer http.ResponseWriter, request *http.Request)

	GetAllVehicleService(writer http.ResponseWriter, request *http.Request)
	GetVehicleServiceById(writer http.ResponseWriter, request *http.Request)
	UpdateVehicleService(writer http.ResponseWriter, request *http.Request)
	AddVehicleService(writer http.ResponseWriter, request *http.Request)
	DeleteVehicleService(writer http.ResponseWriter, request *http.Request)
	DeleteVehicleServiceMultiId(writer http.ResponseWriter, request *http.Request)

	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	New(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	CloseOrder(writer http.ResponseWriter, request *http.Request)

	NewStatus(writer http.ResponseWriter, request *http.Request)
	AddStatus(writer http.ResponseWriter, request *http.Request)
	UpdateStatus(writer http.ResponseWriter, request *http.Request)
	DeleteStatus(writer http.ResponseWriter, request *http.Request)

	NewType(writer http.ResponseWriter, request *http.Request)
	AddType(writer http.ResponseWriter, request *http.Request)
	UpdateType(writer http.ResponseWriter, request *http.Request)
	DeleteType(writer http.ResponseWriter, request *http.Request)

	NewLineType(writer http.ResponseWriter, request *http.Request)
	AddLineType(writer http.ResponseWriter, request *http.Request)
	UpdateLineType(writer http.ResponseWriter, request *http.Request)
	DeleteLineType(writer http.ResponseWriter, request *http.Request)

	NewBill(writer http.ResponseWriter, request *http.Request)
	AddBill(writer http.ResponseWriter, request *http.Request)
	UpdateBill(writer http.ResponseWriter, request *http.Request)
	DeleteBill(writer http.ResponseWriter, request *http.Request)

	NewTrxType(writer http.ResponseWriter, request *http.Request)
	AddTrxType(writer http.ResponseWriter, request *http.Request)
	UpdateTrxType(writer http.ResponseWriter, request *http.Request)
	DeleteTrxType(writer http.ResponseWriter, request *http.Request)

	NewTrxTypeSo(writer http.ResponseWriter, request *http.Request)
	AddTrxTypeSo(writer http.ResponseWriter, request *http.Request)
	UpdateTrxTypeSo(writer http.ResponseWriter, request *http.Request)
	DeleteTrxTypeSo(writer http.ResponseWriter, request *http.Request)

	NewJobType(writer http.ResponseWriter, request *http.Request)
	AddJobType(writer http.ResponseWriter, request *http.Request)
	UpdateJobType(writer http.ResponseWriter, request *http.Request)
	DeleteJobType(writer http.ResponseWriter, request *http.Request)

	NewDropPoint(writer http.ResponseWriter, request *http.Request)
	NewVehicleBrand(writer http.ResponseWriter, request *http.Request)
	NewVehicleModel(writer http.ResponseWriter, request *http.Request)
	GenerateDocumentNumber(writer http.ResponseWriter, request *http.Request)

	GetAllDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	GetDetailByIdWorkOrder(writer http.ResponseWriter, request *http.Request)
	AddDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	UpdateDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	DeleteDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	DeleteDetailWorkOrderMultiId(writer http.ResponseWriter, request *http.Request)

	GetAllBooking(writer http.ResponseWriter, request *http.Request)
	GetBookingById(writer http.ResponseWriter, request *http.Request)
	NewBooking(writer http.ResponseWriter, request *http.Request)
	SaveBooking(writer http.ResponseWriter, request *http.Request)

	GetAllAffiliated(writer http.ResponseWriter, request *http.Request)
	GetAffiliatedById(writer http.ResponseWriter, request *http.Request)
	NewAffiliated(writer http.ResponseWriter, request *http.Request)
	SaveAffiliated(writer http.ResponseWriter, request *http.Request)

	DeleteCampaign(writer http.ResponseWriter, request *http.Request)
	AddContractService(writer http.ResponseWriter, request *http.Request)
	AddGeneralRepairPackage(writer http.ResponseWriter, request *http.Request)
	AddFieldAction(writer http.ResponseWriter, request *http.Request)
	ChangeBillTo(writer http.ResponseWriter, request *http.Request)
	ChangePhoneNo(writer http.ResponseWriter, request *http.Request)
	ConfirmPrice(writer http.ResponseWriter, request *http.Request)
	GetServiceRequestByWO(writer http.ResponseWriter, request *http.Request)
	GetClaimByWO(writer http.ResponseWriter, request *http.Request)
	GetClaimItemByWO(writer http.ResponseWriter, request *http.Request)
	GetWOByBillCode(writer http.ResponseWriter, request *http.Request)
	GetDetailWOByClaimBillCode(writer http.ResponseWriter, request *http.Request)
	GetDetailWOByBillCode(writer http.ResponseWriter, request *http.Request)
	GetDetailWOByATPMBillCode(writer http.ResponseWriter, request *http.Request)
	GetSupplyByWO(writer http.ResponseWriter, request *http.Request)
}

func NewWorkOrderController(WorkOrderService transactionworkshopservice.WorkOrderService) WorkOrderController {
	return &WorkOrderControllerImpl{
		WorkOrderService: WorkOrderService,
	}
}

// GetAll gets all work orders
// @Summary Get All Work Orders
// @Description Retrieve all work orders with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number query string false "Work Order System Number"
// @Param work_order_type_id query string false "Work Order Type ID"
// @Param brand_id query string false "Brand ID"
// @Param model_id query string false "Model ID"
// @Param vehicle_id query string false "Vehicle ID"
// @Param work_order_date query string false "Work Order Date"
// @Param work_order_close_date query string false "Work Order Close Date"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order [get]
func (r *WorkOrderControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.work_order_document_number":        queryValues.Get("work_order_document_number"),
		"trx_work_order.work_order_system_number":          queryValues.Get("work_order_system_number"),
		"trx_work_order.work_order_date_from":              queryValues.Get("work_order_date_from"),
		"trx_work_order.work_order_date_to":                queryValues.Get("work_order_date_to"),
		"trx_work_order.work_order_type_id":                queryValues.Get("work_order_type_id"),
		"trx_work_order.work_order_type_description":       queryValues.Get("work_order_type_description"),
		"trx_work_order.brand_id":                          queryValues.Get("brand_id"),
		"trx_work_order.brand_name":                        queryValues.Get("brand_name"),
		"trx_work_order.model_id":                          queryValues.Get("model_id"),
		"trx_work_order.model_name":                        queryValues.Get("model_name"),
		"trx_work_order.vehicle_id":                        queryValues.Get("vehicle_id"),
		"trx_work_order.vehicle_chassis_number":            queryValues.Get("vehicle_chassis_number"),
		"trx_work_order.vehicle_tnkb":                      queryValues.Get("vehicle_tnkb"),
		"trx_work_order.work_order_status_id":              queryValues.Get("work_order_status_id"),
		"trx_work_order.work_order_status_name":            queryValues.Get("work_order_status_name"),
		"trx_work_order.work_order_repeated_system_number": queryValues.Get("work_order_repeated_system_number"),
		"trx_work_order.variant_id":                        queryValues.Get("variant_id"),
		"trx_work_order.foreman_id":                        queryValues.Get("foreman_id"),
		"trx_work_order.service_advisor_id":                queryValues.Get("service_advisor_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAll(criteria, paginate)
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

// New creates a new work order
// @Summary Create New Work Order
// @Description Create a new work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param reqBody body transactionworkshoppayloads.WorkOrderNormalRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal [post]
func (r *WorkOrderControllerImpl) New(writer http.ResponseWriter, request *http.Request) {

	var workOrderRequest transactionworkshoppayloads.WorkOrderNormalRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, err := r.WorkOrderService.New(workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success.WorkOrderSystemNumber > 0 {
		payloads.NewHandleSuccess(writer, success, "Work order created successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

}

// NewStatus gets the status of new work orders
// @Summary Get Work Order Statuses
// @Description Retrieve all work order statuses
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-status [get]
func (r *WorkOrderControllerImpl) NewStatus(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	var filters []utils.FilterCondition

	for key, values := range queryParams {
		for _, value := range values {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	statuses, err := r.WorkOrderService.NewStatus(filters)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order statuses", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddStatus adds a new status to a work order
// @Summary Add Work Order Status
// @Description Add a new status to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderStatusRequest true "Work Order Status Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-status [post]
func (r *WorkOrderControllerImpl) AddStatus(writer http.ResponseWriter, request *http.Request) {
	// Add status to work order
	var statusRequest transactionworkshoppayloads.WorkOrderStatusRequest
	helper.ReadFromRequestBody(request, &statusRequest)

	success, err := r.WorkOrderService.AddStatus(statusRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Status added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// UpdateStatus updates a status of a work order
// @Summary Update Work Order Status
// @Description Update a status of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_status_id path string true "Work Order Status ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderStatusRequest true "Work Order Status Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-status/{work_order_status_id} [put]
func (r *WorkOrderControllerImpl) UpdateStatus(writer http.ResponseWriter, request *http.Request) {
	// Update status of a work order
	statusID, _ := strconv.Atoi(chi.URLParam(request, "work_order_status_id"))

	var statusRequest transactionworkshoppayloads.WorkOrderStatusRequest
	helper.ReadFromRequestBody(request, &statusRequest)

	update, err := r.WorkOrderService.UpdateStatus(int(statusID), statusRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if update {
		payloads.NewHandleSuccess(writer, update, "Status updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// DeleteStatus deletes a status from a work order
// @Summary Delete Work Order Status
// @Description Delete a status from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_status_id path string true "Work Order Status ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-status/{work_order_status_id} [delete]
func (r *WorkOrderControllerImpl) DeleteStatus(writer http.ResponseWriter, request *http.Request) {
	// Delete status from work order
	statusID, _ := strconv.Atoi(chi.URLParam(request, "work_order_status_id"))

	delete, err := r.WorkOrderService.DeleteStatus(int(statusID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, delete, "Status deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// NewLineType gets the LineType of new work orders
// @Summary Get Work Order LineType
// @Description Retrieve all work order LineType
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-line-type [get]
func (r *WorkOrderControllerImpl) NewLineType(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	var filters []utils.FilterCondition

	for key, values := range queryParams {
		for _, value := range values {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	statuses, err := r.WorkOrderService.NewLineType(filters)
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}
	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order line type", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddLineType adds a new LineType to a work order
// @Summary Add Work Order LineType
// @Description Add a new LineType to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.Linetype true "Work Order Bill Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-line-type [post]
func (r *WorkOrderControllerImpl) AddLineType(writer http.ResponseWriter, request *http.Request) {
	// Add line type to work order
	var linetypeRequest transactionworkshoppayloads.Linetype
	helper.ReadFromRequestBody(request, &linetypeRequest)

	success, err := r.WorkOrderService.AddLineType(linetypeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Line Type added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// UpdateLineType updates a LineType of a work order
// @Summary Update Work Order LineType
// @Description Update a LineType of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param line_type_id path string true "Work Order LineType ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderLinetypeRequest true "Work Order Bill Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-line-type/{line_type_id} [put]
func (r *WorkOrderControllerImpl) UpdateLineType(writer http.ResponseWriter, request *http.Request) {
	// Update lineTypeID of a work order
	lineTypeID, _ := strconv.Atoi(chi.URLParam(request, "line_type_id"))

	var linetypeRequest transactionworkshoppayloads.Linetype
	helper.ReadFromRequestBody(request, &linetypeRequest)

	update, err := r.WorkOrderService.UpdateLineType(int(lineTypeID), linetypeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	if update {
		payloads.NewHandleSuccess(writer, update, "Line Type updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// DeleteLineType deletes a LineType from a work order
// @Summary Delete Work Order LineType
// @Description Delete a LineType from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param line_type_id path string true "Work Order LineType ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-line-type/{line_type_id} [delete]
func (r *WorkOrderControllerImpl) DeleteLineType(writer http.ResponseWriter, request *http.Request) {
	// Delete DeleteLineType from work order
	lineTypeID, _ := strconv.Atoi(chi.URLParam(request, "line_type_id"))

	delete, err := r.WorkOrderService.DeleteLineType(int(lineTypeID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, delete, "LineType deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// NewBill gets the bill of new work orders
// @Summary Get Work Order Bill
// @Description Retrieve all work order bill
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-bill [get]
func (r *WorkOrderControllerImpl) NewBill(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	var filters []utils.FilterCondition

	for key, values := range queryParams {
		for _, value := range values {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	statuses, err := r.WorkOrderService.NewBill(filters)
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order bill able to", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddBill adds a new bill to a work order
// @Summary Add Work Order Bill
// @Description Add a new bill to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderBillableRequest true "Work Order Bill Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-bill [post]
func (r *WorkOrderControllerImpl) AddBill(writer http.ResponseWriter, request *http.Request) {
	// Add bill to work order
	var billRequest transactionworkshoppayloads.WorkOrderBillableRequest
	helper.ReadFromRequestBody(request, &billRequest)

	success, err := r.WorkOrderService.AddBill(billRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Bill added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// UpdateBill updates a bill of a work order
// @Summary Update Work Order Bill
// @Description Update a bill of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_bill_id path string true "Work Order Bill ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBillableRequest true "Work Order Bill Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-bill/{work_order_bill_id} [put]
func (r *WorkOrderControllerImpl) UpdateBill(writer http.ResponseWriter, request *http.Request) {
	// Update bill of a work order
	billID, _ := strconv.Atoi(chi.URLParam(request, "work_order_bill_id"))

	var billRequest transactionworkshoppayloads.WorkOrderBillableRequest
	helper.ReadFromRequestBody(request, &billRequest)

	update, err := r.WorkOrderService.UpdateBill(int(billID), billRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	if update {
		payloads.NewHandleSuccess(writer, update, "Bill updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// DeleteBill deletes a bill from a work order
// @Summary Delete Work Order Bill
// @Description Delete a bill from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_bill_id path string true "Work Order Bill ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-bill/{work_order_bill_id} [delete]
func (r *WorkOrderControllerImpl) DeleteBill(writer http.ResponseWriter, request *http.Request) {
	// Delete bill from work order
	billID, _ := strconv.Atoi(chi.URLParam(request, "work_order_bill_id"))

	delete, err := r.WorkOrderService.DeleteBill(int(billID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, delete, "Bill deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// NewType gets the types of new work orders
// @Summary Get Work Order Types
// @Description Retrieve all work order types
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-type [get]
func (r *WorkOrderControllerImpl) NewType(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	var filters []utils.FilterCondition

	for key, values := range queryParams {
		for _, value := range values {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	statuses, err := r.WorkOrderService.NewType(filters)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order type", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddType adds a new type to a work order
// @Summary Add Work Order Type
// @Description Add a new type to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderTypeRequest true "Work Order Type Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-type [post]
func (r *WorkOrderControllerImpl) AddType(writer http.ResponseWriter, request *http.Request) {
	// Add type to work order
	var typeRequest transactionworkshoppayloads.WorkOrderTypeRequest
	helper.ReadFromRequestBody(request, &typeRequest)

	success, err := r.WorkOrderService.AddType(typeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Type added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// UpdateType updates a type of a work order
// @Summary Update Work Order Type
// @Description Update a type of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_type_id path string true "Work Order Type ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderTypeRequest true "Work Order Type Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-type/{work_order_type_id} [put]
func (r *WorkOrderControllerImpl) UpdateType(writer http.ResponseWriter, request *http.Request) {
	// Update type of a work order
	typeID, _ := strconv.Atoi(chi.URLParam(request, "work_order_type_id"))

	var typeRequest transactionworkshoppayloads.WorkOrderTypeRequest
	helper.ReadFromRequestBody(request, &typeRequest)

	update, err := r.WorkOrderService.UpdateType(int(typeID), typeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if update {
		payloads.NewHandleSuccess(writer, update, "Type updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// DeleteType deletes a type from a work order
// @Summary Delete Work Order Type
// @Description Delete a type from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_type_id path string true "Work Order Type ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-type/{work_order_type_id} [delete]
func (r *WorkOrderControllerImpl) DeleteType(writer http.ResponseWriter, request *http.Request) {
	// Delete type from work order
	typeID, _ := strconv.Atoi(chi.URLParam(request, "work_order_type_id"))

	delete, err := r.WorkOrderService.DeleteType(int(typeID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, delete, "Type deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// NewDropPoint gets the drop points of new work orders
// @Summary Get Work Order Drop Points
// @Description Retrieve all work order drop points
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-drop-point [get]
func (r *WorkOrderControllerImpl) NewDropPoint(writer http.ResponseWriter, request *http.Request) {

	statuses, err := r.WorkOrderService.NewDropPoint()
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order drop point", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// NewVehicleBrand gets the vehicle brands of new work orders
// @Summary Get Work Order Vehicle Brands
// @Description Retrieve all work order vehicle brands
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-vehicle-brand [get]
func (r *WorkOrderControllerImpl) NewVehicleBrand(writer http.ResponseWriter, request *http.Request) {

	statuses, err := r.WorkOrderService.NewVehicleBrand()
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order vehicle brand", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// NewVehicleModel gets the vehicle models of new work orders
// @Summary Get Work Order Vehicle Models Based Brand ID
// @Description Retrieve all work order vehicle models based brand ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param brand_id query string true "Brand ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-vehicle-model/{brand_id} [get]
func (r *WorkOrderControllerImpl) NewVehicleModel(writer http.ResponseWriter, request *http.Request) {
	brandIdStr := chi.URLParam(request, "brand_id")
	brandId, err := strconv.Atoi(brandIdStr)
	if err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{Message: "Invalid brand ID"})
		return
	}

	create, baseErr := r.WorkOrderService.NewVehicleModel(brandId)

	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	if len(create) > 0 {
		payloads.NewHandleSuccess(writer, create, "List of work order vehicle model", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetAllService gets all services of a work order
// @Summary Get All Services of Work Order
// @Description Retrieve all services of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/requestservice [get]
func (r *WorkOrderControllerImpl) GetAllRequest(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	excludeParams := map[string]bool{
		"page":    true,
		"limit":   true,
		"sort_of": true,
		"sort_by": true,
	}

	filterConditions := make([]utils.FilterCondition, 0)
	for key, values := range queryValues {
		if len(values) > 0 && !excludeParams[key] {
			filterConditions = append(filterConditions, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: values[0],
			})
		}
	}

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAllRequest(filterConditions, paginate)
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

// GetServiceById gets a service of a work order by ID
// @Summary Get Service of Work Order By ID
// @Description Retrieve a service of a work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_id path string true "Work Order Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservice/{work_order_service_id} [get]
func (r *WorkOrderControllerImpl) GetRequestById(writer http.ResponseWriter, request *http.Request) {
	// Get service of a work order by ID
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	serviceID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_id"))

	service, err := r.WorkOrderService.GetRequestById(int(workorderID), int(serviceID))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, service, "Get Data Successfully", http.StatusOK)
}

// UpdateRequest updates a request of a work order
// @Summary Update Request of Work Order
// @Description Update a request of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_id path string true "Work Order Service ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservice/{work_order_service_id} [put]
func (r *WorkOrderControllerImpl) UpdateRequest(writer http.ResponseWriter, request *http.Request) {
	// Update request of a work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	requestID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_id"))

	var groupRequest transactionworkshoppayloads.WorkOrderServiceRequest
	helper.ReadFromRequestBody(request, &groupRequest)

	update, err := r.WorkOrderService.UpdateRequest(int(workorderID), int(requestID), groupRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Request updated successfully", http.StatusOK)

}

// AddRequest adds a new request to a work order
// @Summary Add Request to Work Order
// @Description Add a new request to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservice [post]
func (r *WorkOrderControllerImpl) AddRequest(writer http.ResponseWriter, request *http.Request) {
	// Add request to work order\
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var groupRequest transactionworkshoppayloads.WorkOrderServiceRequest
	helper.ReadFromRequestBody(request, &groupRequest)

	success, err := r.WorkOrderService.AddRequest(int(workorderID), groupRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success.WorkOrderServiceId > 0 {
		payloads.NewHandleSuccess(writer, success, "Request added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddRequest Multi adds multiple request to a work order
// @Summary Add Multiple Request to Work Order
// @Description Add multiple request to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservicemulti [post]
func (r *WorkOrderControllerImpl) AddRequestMultiId(writer http.ResponseWriter, request *http.Request) {
	// Add request to work order
	workorderID, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	var groupRequests []transactionworkshoppayloads.WorkOrderServiceRequest
	helper.ReadFromRequestBody(request, &groupRequests)

	entities, baseErr := r.WorkOrderService.AddRequestMultiId(workorderID, groupRequests) // Call the modified service method
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, entities, "Requests added successfully", http.StatusCreated)
}

// DeleteRequest deletes a request from a work order
// @Summary Delete Request from Work Order
// @Description Delete a request from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_id path string true "Work Order Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservice/{work_order_service_id} [delete]
func (r *WorkOrderControllerImpl) DeleteRequest(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	requestID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_id"))

	delete, err := r.WorkOrderService.DeleteRequest(int(workorderID), int(requestID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, delete, "Request deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

}

// DeleteRequestMultiId deletes multiple request from a work order
// @Summary Delete Multiple Request from Work Order
// @Description Delete multiple request from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @param multi_id query string true "Multiple Request ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservice/{multi_id} [delete]
func (r *WorkOrderControllerImpl) DeleteRequestMultiId(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderstrID := chi.URLParam(request, "work_order_system_number")
	workorderID, err := strconv.Atoi(workorderstrID)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var intIds []int
	for _, element := range elements {
		num, err := strconv.Atoi(strings.TrimSpace(element))
		if err != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		intIds = append(intIds, num)
	}

	success, baseErr := r.WorkOrderService.DeleteRequestMultiId(workorderID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "request detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Service Detail deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete service detail", http.StatusInternalServerError)
	}

}

// GetAllVehicleService gets all vehicle services of a work order
// @Summary Get All Vehicle Services of Work Order
// @Description Retrieve all vehicle services of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/vehicleservice [get]
func (r *WorkOrderControllerImpl) GetAllVehicleService(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	excludeParams := map[string]bool{
		"page":    true,
		"limit":   true,
		"sort_of": true,
		"sort_by": true,
	}

	filterConditions := make([]utils.FilterCondition, 0)
	for key, values := range queryValues {
		if len(values) > 0 && !excludeParams[key] {
			filterConditions = append(filterConditions, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: values[0],
			})
		}
	}

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAllVehicleService(filterConditions, paginate)
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

// GetVehicleServiceById gets a vehicle service of a work order by ID
// @Summary Get Vehicle Service of Work Order By ID
// @Description Retrieve a vehicle service of a work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_vehicle_id path string true "Work Order Vehicle Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice/{work_order_service_vehicle_id} [get]
func (r *WorkOrderControllerImpl) GetVehicleServiceById(writer http.ResponseWriter, request *http.Request) {
	// Get vehicle service of a work order by ID
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	vehicleServiceID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_vehicle_id"))

	service, err := r.WorkOrderService.GetVehicleServiceById(int(workorderID), int(vehicleServiceID))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, service, "Get Data Successfully", http.StatusOK)
}

// UpdateVehicleService updates a vehicle service of a work order
// @Summary Update Vehicle Service of Work Order
// @Description Update a vehicle service of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_vehicle_service_id path string true "Work Order Vehicle Service ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceVehicleRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice/{work_order_vehicle_service_id} [put]
func (r *WorkOrderControllerImpl) UpdateVehicleService(writer http.ResponseWriter, request *http.Request) {
	// Update vehicle service of a work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	vehicleServiceID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_vehicle_id"))

	var vehicleRequest transactionworkshoppayloads.WorkOrderServiceVehicleRequest
	helper.ReadFromRequestBody(request, &vehicleRequest)

	update, err := r.WorkOrderService.UpdateVehicleService(int(workorderID), int(vehicleServiceID), vehicleRequest)

	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Vehicle service updated successfully", http.StatusOK)
}

// AddVehicleService adds a new vehicle service to a work order
// @Summary Add Vehicle Service to Work Order
// @Description Add a new vehicle service to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceVehicleRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice [post]
func (r *WorkOrderControllerImpl) AddVehicleService(writer http.ResponseWriter, request *http.Request) {
	// Add vehicle service to work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var vehicleRequest transactionworkshoppayloads.WorkOrderServiceVehicleRequest
	helper.ReadFromRequestBody(request, &vehicleRequest)

	success, err := r.WorkOrderService.AddVehicleService(int(workorderID), vehicleRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success.WorkOrderServiceVehicleId > 0 {
		payloads.NewHandleSuccess(writer, success, "Vehicle service added successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to add vehicle service", http.StatusInternalServerError)
	}
}

// DeleteVehicleService deletes a vehicle service from a work order
// @Summary Delete Vehicle Service from Work Order
// @Description Delete a vehicle service from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_vehicle_service_id path string true "Work Order Vehicle Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice/{work_order_vehicle_service_id} [delete]
func (r *WorkOrderControllerImpl) DeleteVehicleService(writer http.ResponseWriter, request *http.Request) {
	// Delete vehicle service from work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	vehicleServiceID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_vehicle_id"))

	delete, err := r.WorkOrderService.DeleteVehicleService(int(workorderID), int(vehicleServiceID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, delete, "Vehicle service deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

}

// DeleteVehicleServiceMultiId deletes multiple a vehicle service from a work order
// @Summary Delete multiple vehicle service
// @Description  Delete multiple vehicle service
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Detail
// @Param work_order_system_number path string true "Service Detail System ID"
// @Param multi_id path string true "Service Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice/{multi_id} [delete]
func (r *WorkOrderControllerImpl) DeleteVehicleServiceMultiId(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderstrID := chi.URLParam(request, "work_order_system_number")
	workorderID, err := strconv.Atoi(workorderstrID)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var intIds []int
	for _, element := range elements {
		num, err := strconv.Atoi(strings.TrimSpace(element))
		if err != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		intIds = append(intIds, num)
	}

	success, baseErr := r.WorkOrderService.DeleteVehicleServiceMultiId(workorderID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "request detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Vehicle service deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete Vehicle detail", http.StatusInternalServerError)
	}

}

// GetById handles the transaction for all work orders
// @Summary Get Work Order By ID
// @Description Retrieve work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	workOrder, baseErr := r.WorkOrderService.GetById(workOrderId, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, workOrder, "Get Data Successfully", http.StatusOK)
}

// Save saves a new work order
// @Summary Save Work Order
// @Description Save a new work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderNormalSaveRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number} [put]
func (r *WorkOrderControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	// Get the Work Order ID from URL parameters and convert to int
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.WorkOrderNormalSaveRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, baseErr := r.WorkOrderService.Save(workOrderRequest, workOrderId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Work order saved successfully", http.StatusOK)

}

// Submit submits a new work order
// @Summary Submit Work Order
// @Description Submit a new work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path int true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/submit/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Retrieve work order ID from URL parameters
	workOrderId := chi.URLParam(request, "work_order_system_number")
	workOrderIdInt, err := strconv.Atoi(workOrderId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, newDocumentNumber, baseErr := r.WorkOrderService.Submit(workOrderIdInt)
	if baseErr != nil {
		if baseErr.Message == "Document number has already been generated" {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		} else if baseErr.Message == "No work order data found" {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		responseData := transactionworkshoppayloads.SubmitWorkOrderResponse{
			DocumentNumber:        newDocumentNumber,
			WorkOrderSystemNumber: workOrderIdInt,
		}
		payloads.NewHandleSuccess(writer, responseData, "Work order submitted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to submit work order", http.StatusInternalServerError)
	}
}

// Void delete or cancel a work order
// @Summary Void Work Order
// @Description Delete or cancel a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path int true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/void/{work_order_system_number} [delete]
func (r *WorkOrderControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Void work order
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.WorkOrderService.Void(workOrderId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Work order voided successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to void work order", http.StatusInternalServerError)
	}
}

// CloseOrder closes a work order
// @Summary Close Work Order
// @Description Close an existing work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path int true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/close/{work_order_system_number} [patch]
func (r *WorkOrderControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close work order
	workOrderId := chi.URLParam(request, "work_order_system_number")
	workOrderIdInt, err := strconv.Atoi(workOrderId)

	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.WorkOrderService.CloseOrder(workOrderIdInt)
	if baseErr != nil {
		switch baseErr.Message {
		case "Work order not found":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		case "Work order cannot be closed because status is draft":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "There is still DP payment that has not been settled":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Detail Work Order without Invoice No must be deleted":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Warranty Item (PTP) must be supplied":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Warranty Item (PTM)/Operation must be Invoiced":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Service Mileage must be larger than Last Mileage.":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		default:
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Work order closed successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to close work order", http.StatusInternalServerError)
	}

}

// GetWorkOrderDetail gets the detail of a work order
// @Summary Get Work Order Detail
// @Description Retrieve the detail of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal Detail
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/detail [get]
func (r *WorkOrderControllerImpl) GetAllDetailWorkOrder(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order_detail.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order_detail.work_order_detail_id":     queryValues.Get("work_order_detail_id"),
		"trx_work_order_detail.transaction_type_id":      queryValues.Get("transaction_type_id"),
		"trx_work_order_detail.supply_quantity":          queryValues.Get("supply_quantity"),
		"trx_work_order_detail.price_list_id":            queryValues.Get("price_list_id"),
		"trx_work_order_detail.line_type_id":             queryValues.Get("line_type_id"),
		"trx_work_order_detail.job_type_id":              queryValues.Get("job_type_id"),
		"trx_work_order_detail.frt_quantity":             queryValues.Get("frt_quantity"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAllDetailWorkOrder(criteria, paginate)
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

// GetDetailWorkOrderById gets the detail of a work order by ID
// @Summary Get Work Order Detail By ID
// @Description Retrieve the detail of a work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_detail_id path string true "Work Order Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail/{work_order_detail_id} [get]
func (r *WorkOrderControllerImpl) GetDetailByIdWorkOrder(writer http.ResponseWriter, request *http.Request) {
	// Get the detail of a work order by ID
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}
	detailId, err := strconv.Atoi(chi.URLParam(request, "work_order_detail_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order detail ID", http.StatusBadRequest)
		return
	}

	detail, baseErr := r.WorkOrderService.GetDetailByIdWorkOrder(workOrderId, detailId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, detail, "Get Data Successfully", http.StatusOK)

}

// UpdateDetailWorkOrder updates the detail of a work order
// @Summary Update Work Order Detail
// @Description Update the detail of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_detail_id path string true "Work Order Detail ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderDetailRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail/{work_order_detail_id} [put]
func (r *WorkOrderControllerImpl) UpdateDetailWorkOrder(writer http.ResponseWriter, request *http.Request) {
	// Update the detail of a work order
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	detailId, _ := strconv.Atoi(chi.URLParam(request, "work_order_detail_id"))

	var detailRequest transactionworkshoppayloads.WorkOrderDetailRequest
	helper.ReadFromRequestBody(request, &detailRequest)

	update, err := r.WorkOrderService.UpdateDetailWorkOrder(int(workOrderId), int(detailId), detailRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if update.WorkOrderSystemNumber > 0 {
		payloads.NewHandleSuccess(writer, update, "Detail updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddDetailWorkOrder adds a new detail to a work order
// @Summary Add Work Order Detail
// @Description Add a new detail to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderDetailRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail [post]
func (r *WorkOrderControllerImpl) AddDetailWorkOrder(writer http.ResponseWriter, request *http.Request) {
	workOrderStrId := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var detailRequest transactionworkshoppayloads.WorkOrderDetailRequest
	helper.ReadFromRequestBody(request, &detailRequest)

	success, serviceErr := r.WorkOrderService.AddDetailWorkOrder(workOrderId, detailRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Detail added successfully", http.StatusCreated)

}

// DeleteDetailWorkOrder deletes a detail from a work order
// @Summary Delete Work Order Detail
// @Description Delete a detail from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_detail_id path string true "Work Order Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail/{work_order_detail_id} [delete]
func (r *WorkOrderControllerImpl) DeleteDetailWorkOrder(writer http.ResponseWriter, request *http.Request) {
	// Delete a detail from a work order
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	detailId, _ := strconv.Atoi(chi.URLParam(request, "work_order_detail_id"))

	delete, err := r.WorkOrderService.DeleteDetailWorkOrder(int(workOrderId), int(detailId))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, delete, "Detail deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

}

// DeleteDetailWorkOrder deletes a multiple detail from a work order
// @Summary Delete multiple Work Order Detail
// @Description  Delete multiple Work Order Detail
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal Detail
// @Param work_order_system_number path string true "Work Order System Number ID"
// @Param multi_id path string true "Work Order Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail/{multi_id} [delete]
func (r *WorkOrderControllerImpl) DeleteDetailWorkOrderMultiId(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderstrID := chi.URLParam(request, "work_order_system_number")
	workorderID, err := strconv.Atoi(workorderstrID)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var intIds []int
	for _, element := range elements {
		num, err := strconv.Atoi(strings.TrimSpace(element))
		if err != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		intIds = append(intIds, num)
	}

	success, baseErr := r.WorkOrderService.DeleteDetailWorkOrderMultiId(workorderID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "request detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Vehicle service deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete Vehicle detail", http.StatusInternalServerError)
	}

}

// GetAllWorkOrderBooking gets all work order bookings
// @Summary Get All Work Order Booking
// @Description Retrieve all work order bookings
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Booking
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking [get]
func (r *WorkOrderControllerImpl) GetAllBooking(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.work_order_system_number":      chi.URLParam(request, "work_order_system_number"),
		"trx_work_order.booking_system_number":         chi.URLParam(request, "booking_system_number"),
		"trx_work_order.estimation_system_number":      chi.URLParam(request, "estimation_system_number"),
		"trx_work_order.service_request_system_number": chi.URLParam(request, "service_request_system_number"),
		"trx_work_order.brand_id":                      chi.URLParam(request, "brand_id"),
		"trx_work_order.model_id":                      chi.URLParam(request, "model_id"),
		"trx_work_order.vehicle_id":                    chi.URLParam(request, "vehicle_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAllBooking(criteria, paginate)
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

// GetWorkOrderBookingById gets a work
// @Summary Get Work Order Booking By ID
// @Description Retrieve a work order booking by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Booking
// @Param work_order_system_number path string true "Work Order ID"
// @Param booking_system_number path string true "Work Order Booking ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/{work_order_system_number}/{booking_system_number} [get]
func (r *WorkOrderControllerImpl) GetBookingById(writer http.ResponseWriter, request *http.Request) {
	// Get a work order booking by ID
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	workOrderBookId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	workOrder, baseErr := r.WorkOrderService.GetBookingById(workOrderId, workOrderBookId, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}
	payloads.NewHandleSuccess(writer, workOrder, "Get Data Successfully", http.StatusOK)
}

// UpdateWorkOrderBooking updates a work order booking
// @Summary Update Work Order Booking
// @Description Update a work order booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Booking
// @Param work_order_system_number path string true "Work Order ID"
// @Param booking_system_number path string true "Work Order Booking ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBookingRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/{work_order_system_number}/{booking_system_number} [put]
func (r *WorkOrderControllerImpl) SaveBooking(writer http.ResponseWriter, request *http.Request) {
	// Update a work order booking
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	workOrderBookId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	var workOrderRequest transactionworkshoppayloads.WorkOrderBookingRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, err := r.WorkOrderService.SaveBooking(workOrderId, workOrderBookId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if result {
		payloads.NewHandleSuccess(writer, result, "Work order saved successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to save work order", http.StatusInternalServerError)
	}

}

// AddWorkOrderBooking adds a new work order booking
// @Summary Add Work Order Booking
// @Description Add a new work order booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Booking
// @Param reqBody body transactionworkshoppayloads.WorkOrderBookingRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking [post]
func (r *WorkOrderControllerImpl) NewBooking(writer http.ResponseWriter, request *http.Request) {

	var workOrderRequest transactionworkshoppayloads.WorkOrderBookingRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, err := r.WorkOrderService.NewBooking(workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order saved successfully", http.StatusCreated)

}

// GetAllAffiliated gets all affiliated work orders
// @Summary Get All Affiliated Work Orders
// @Description Retrieve all affiliated work orders
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Affiliated
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated [get]
func (r *WorkOrderControllerImpl) GetAllAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Get all affiliated work orders
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAllAffiliated(criteria, paginate)
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

// GetAffiliatedById gets an affiliated work order by ID
// @Summary Get Affiliated Work Order By ID
// @Description Retrieve an affiliated work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Affiliated
// @Param work_order_system_number path string true "Work Order ID"
// @Param service_request_system_number path string true "Affiliated Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated/{work_order_system_number}/{service_request_system_number} [get]
func (r *WorkOrderControllerImpl) GetAffiliatedById(writer http.ResponseWriter, request *http.Request) {
	// Get affiliated work order by ID
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	affiliatedWorkOrderId, _ := strconv.Atoi(chi.URLParam(request, "service_request_system_number"))

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	workOrder, baseErr := r.WorkOrderService.GetAffiliatedById(workOrderId, affiliatedWorkOrderId, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, workOrder, "Get Data Successfully", http.StatusOK)
}

// NewAffiliated creates a new affiliated work order
// @Summary Create New Affiliated Work Order
// @Description Create a new affiliated work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Affiliated
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated [post]
func (r *WorkOrderControllerImpl) NewAffiliated(writer http.ResponseWriter, request *http.Request) {
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var workOrderRequest transactionworkshoppayloads.WorkOrderAffiliatedRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, err := r.WorkOrderService.NewAffiliated(workOrderId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order added successfully", http.StatusCreated)
}

// UpdateAffiliated updates an affiliated work order
// @Summary Update Affiliated Work Order
// @Description Update an affiliated work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Affiliated
// @Param work_order_system_number path string true "Work Order ID"
// @Param affiliated_work_order_system_number path string true "Affiliated Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderAffiliatedRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated/{work_order_system_number}/{affiliated_work_order_system_number} [put]
func (r *WorkOrderControllerImpl) SaveAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Update an affiliated work order
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	affiliatedWorkOrderId, _ := strconv.Atoi(chi.URLParam(request, "affiliated_work_order_system_number"))

	var workOrderRequest transactionworkshoppayloads.WorkOrderAffiliatedRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, err := r.WorkOrderService.SaveAffiliated(workOrderId, affiliatedWorkOrderId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order updated successfully", http.StatusOK)
}

// GenerateWorkOrderDocumentNumber generates a new work order document number
// @Summary Generate Work Order Document Number
// @Description Generate a new work order document number
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path string true "Work Order ID"
// Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/document-number/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) GenerateDocumentNumber(writer http.ResponseWriter, request *http.Request) {
	// Generate a new work order document number
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	result, err := r.WorkOrderService.GenerateDocumentNumber(workOrderId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Document number generated successfully", http.StatusOK)

}

// ChangeBillTo changes the bill to of a work order
// @Summary Change Bill To
// @Description Change the bill to of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.ChangeBillToRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/change-bill-to/{work_order_system_number} [put]
func (r *WorkOrderControllerImpl) ChangeBillTo(writer http.ResponseWriter, request *http.Request) {
	// Change the bill to of a work order
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.ChangeBillToRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, baseErr := r.WorkOrderService.ChangeBillTo(workOrderId, workOrderRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Bill to changed successfully", http.StatusOK)
}

// ChangePhoneNo changes the phone number of a work order
// @Summary Change Phone Number
// @Description Change the phone number of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.ChangePhoneNoRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/change-phone-no/{work_order_system_number} [put]
func (r *WorkOrderControllerImpl) ChangePhoneNo(writer http.ResponseWriter, request *http.Request) {
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.ChangePhoneNoRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	updatedPayload, baseErr := r.WorkOrderService.ChangePhoneNo(workOrderId, workOrderRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, updatedPayload, "Phone number changed successfully", http.StatusOK)
}

// ConfirmPriceList confirms the price list of a work order
// @Summary Confirm Price List
// @Description Confirm the price list of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.ConfirmPriceListRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/confirm-price/{work_order_system_number}/{multi_id} [put]
func (r *WorkOrderControllerImpl) ConfirmPrice(writer http.ResponseWriter, request *http.Request) {
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var intIds []int
	for _, element := range elements {
		num, err := strconv.Atoi(strings.TrimSpace(element))
		if err != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		intIds = append(intIds, num)
	}

	var confirmPriceRequest transactionworkshoppayloads.WorkOrderConfirmPriceRequest
	helper.ReadFromRequestBody(request, &confirmPriceRequest)

	result, baseErr := r.WorkOrderService.ConfirmPrice(workOrderId, intIds, confirmPriceRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Confirm Price list confirmed successfully", http.StatusOK)
}

// NewTrxType gets the Trx Type of new work orders
// @Summary Get Work Order Trx Type
// @Description Retrieve all work order Trx Type
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type [get]
func (r *WorkOrderControllerImpl) NewTrxType(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	var filters []utils.FilterCondition

	for key, values := range queryParams {
		for _, value := range values {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	statuses, err := r.WorkOrderService.NewTrxType(filters)
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order Trx Type", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddTrxType adds a new  Trx Type to a work order
// @Summary Add Work Order  Trx Type
// @Description Add a new  Trx Type to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderTransactionType true "Work Order Transaction Type Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type [post]
func (r *WorkOrderControllerImpl) AddTrxType(writer http.ResponseWriter, request *http.Request) {
	var trxTypeRequest transactionworkshoppayloads.WorkOrderTransactionType
	helper.ReadFromRequestBody(request, &trxTypeRequest)

	success, err := r.WorkOrderService.AddTrxType(trxTypeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Trx Type added successfully", http.StatusCreated)
}

// UpdateTrxType updates a Trx Type of a work order
// @Summary Update Work Order Trx Type
// @Description Update a Trx Type of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_transaction_type_id path string true "Work Order Bill ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBillableRequest true "Work Order Bill Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type/{work_order_transaction_type_id} [put]
func (r *WorkOrderControllerImpl) UpdateTrxType(writer http.ResponseWriter, request *http.Request) {
	// Update a Trx Type of a work order
	trxTypeId, err := strconv.Atoi(chi.URLParam(request, "work_order_transaction_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order Trx Type ID", http.StatusBadRequest)
		return
	}

	var trxTypeRequest transactionworkshoppayloads.WorkOrderTransactionType
	helper.ReadFromRequestBody(request, &trxTypeRequest)

	success, serviceErr := r.WorkOrderService.UpdateTrxType(trxTypeId, trxTypeRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Trx Type updated successfully", http.StatusOK)
}

// DeleteTrxType deletes a Trx Type from a work order
// @Summary Delete Work Order Trx Type
// @Description Delete a Trx Type from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_transaction_type_id path string true "Work Order Trx Type ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type/{work_order_transaction_type_id} [delete]
func (r *WorkOrderControllerImpl) DeleteTrxType(writer http.ResponseWriter, request *http.Request) {

	trxTypeId, err := strconv.Atoi(chi.URLParam(request, "work_order_transaction_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order Trx Type ID", http.StatusBadRequest)
		return
	}

	success, serviceErr := r.WorkOrderService.DeleteTrxType(trxTypeId)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Trx Type deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// NewTrxTypeSo gets the Trx Type of new work orders
// @Summary Get Work Order Trx Type
// @Description Retrieve all work order Trx Type
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type-so [get]
func (r *WorkOrderControllerImpl) NewTrxTypeSo(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	var filters []utils.FilterCondition

	for key, values := range queryParams {
		for _, value := range values {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	statuses, err := r.WorkOrderService.NewTrxTypeSo(filters)
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order Trx Type So", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddTrxTypeSo adds a new  Trx Type to a work order
// @Summary Add Work Order  Trx Type
// @Description Add a new  Trx Type to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderTransactionType true "Work Order Transaction Type Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type-so [post]
func (r *WorkOrderControllerImpl) AddTrxTypeSo(writer http.ResponseWriter, request *http.Request) {
	var trxTypeRequest transactionworkshoppayloads.WorkOrderTransactionType
	helper.ReadFromRequestBody(request, &trxTypeRequest)

	success, err := r.WorkOrderService.AddTrxTypeSo(trxTypeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Trx Type added successfully", http.StatusCreated)
}

// UpdateTrxTypeSo updates a Trx Type of a work order
// @Summary Update Work Order Trx Type
// @Description Update a Trx Type of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_transaction_type_id path string true "Work Order Bill ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBillableRequest true "Work Order Bill Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type-so/{work_order_transaction_type_id} [put]
func (r *WorkOrderControllerImpl) UpdateTrxTypeSo(writer http.ResponseWriter, request *http.Request) {
	// Update a Trx Type of a work order
	trxTypeId, err := strconv.Atoi(chi.URLParam(request, "work_order_transaction_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order Trx Type ID", http.StatusBadRequest)
		return
	}

	var trxTypeRequest transactionworkshoppayloads.WorkOrderTransactionType
	helper.ReadFromRequestBody(request, &trxTypeRequest)

	success, serviceErr := r.WorkOrderService.UpdateTrxTypeSo(trxTypeId, trxTypeRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Trx Type updated successfully", http.StatusOK)
}

// DeleteTrxTypeSo deletes a Trx Type from a work order
// @Summary Delete Work Order Trx Type
// @Description Delete a Trx Type from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_transaction_type_id path string true "Work Order Trx Type ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type/{work_order_transaction_type_id} [delete]
func (r *WorkOrderControllerImpl) DeleteTrxTypeSo(writer http.ResponseWriter, request *http.Request) {

	trxTypeId, err := strconv.Atoi(chi.URLParam(request, "work_order_transaction_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order Trx Type ID", http.StatusBadRequest)
		return
	}

	success, serviceErr := r.WorkOrderService.DeleteTrxTypeSo(trxTypeId)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Trx Type deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// DeleteCampaign deletes a campaign from a work order
// @Summary Delete Work Order Campaign
// @Description Delete a campaign from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order Campaign ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/delete-campaign/{work_order_system_number} [delete]
func (r *WorkOrderControllerImpl) DeleteCampaign(writer http.ResponseWriter, request *http.Request) {

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, serviceErr := r.WorkOrderService.DeleteCampaign(workOrderId)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Campaign deleted successfully", http.StatusOK)

}

// AddContractService adds a new contract service to a work order
// @Summary Add Work Order Contract Service
// @Description Add a new contract service to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderContractServiceRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/contract-service/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) AddContractService(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.WorkOrderContractServiceRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, serviceErr := r.WorkOrderService.AddContractService(workOrderId, workOrderRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Contract service added successfully", http.StatusCreated)
}

// AddGeneralRepairPackage adds a new general repair package to a work order
// @Summary Add Work Order General Repair Package
// @Description Add a new general repair package to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/general-repair-package/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) AddGeneralRepairPackage(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, serviceErr := r.WorkOrderService.AddGeneralRepairPackage(workOrderId, workOrderRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "General repair package added successfully", http.StatusCreated)
}

// AddFieldAction adds a new field action to a work order
// @Summary Add Work Order Field Action
// @Description Add a new field action to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderFieldActionRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/field-action/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) AddFieldAction(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.WorkOrderFieldActionRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, serviceErr := r.WorkOrderService.AddFieldAction(workOrderId, workOrderRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Field action added successfully", http.StatusCreated)
}

// NewJobType gets the Job Type of new work orders
// @Summary Get Work Order Job Type
// @Description Retrieve all work order Job Type
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-job-type [get]
func (r *WorkOrderControllerImpl) NewJobType(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()
	var filters []utils.FilterCondition

	for key, values := range queryParams {
		for _, value := range values {
			filters = append(filters, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: value,
			})
		}
	}

	statuses, err := r.WorkOrderService.NewJobType(filters)
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order Job Type", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddJobType adds a new  Job Type to a work order
// @Summary Add Work Order  Job Type
// @Description Add a new  Job Type to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderTransactionType true "Work Order Transaction Type Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-job-type [post]
func (r *WorkOrderControllerImpl) AddJobType(writer http.ResponseWriter, request *http.Request) {
	var jobTypeRequest transactionworkshoppayloads.WorkOrderJobType
	helper.ReadFromRequestBody(request, &jobTypeRequest)

	success, err := r.WorkOrderService.AddJobType(jobTypeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Job Type added successfully", http.StatusCreated)
}

// UpdateJobType updates a Job Type of a work order
// @Summary Update Work Order Job Type
// @Description Update a Job Type of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_transaction_type_id path string true "Work Order Bill ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBillableRequest true "Work Order Bill Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-transaction-type/{work_order_transaction_type_id} [put]
func (r *WorkOrderControllerImpl) UpdateJobType(writer http.ResponseWriter, request *http.Request) {
	// Update a Trx Type of a work order
	jobTypeId, err := strconv.Atoi(chi.URLParam(request, "job_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order Job Type ID", http.StatusBadRequest)
		return
	}

	var jobTypeRequest transactionworkshoppayloads.WorkOrderJobType
	helper.ReadFromRequestBody(request, &jobTypeRequest)

	success, serviceErr := r.WorkOrderService.UpdateJobType(jobTypeId, jobTypeRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Trx Type updated successfully", http.StatusOK)
}

// DeleteJobType deletes a Job Type from a work order
// @Summary Delete Work Order Job Type
// @Description Delete a Job Type from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_transaction_type_id path string true "Work Order Job Type ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-job-type/{work_order_transaction_type_id} [delete]
func (r *WorkOrderControllerImpl) DeleteJobType(writer http.ResponseWriter, request *http.Request) {

	jobTypeId, err := strconv.Atoi(chi.URLParam(request, "job_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order Job Type ID", http.StatusBadRequest)
		return
	}

	success, serviceErr := r.WorkOrderService.DeleteJobType(jobTypeId)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Trx Type deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetServiceRequestByWO gets all service request by work order
// @Summary Get Service Request By Work Order
// @Description Retrieve all service request by work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/service-request/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetServiceRequestByWO(writer http.ResponseWriter, request *http.Request) {

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.work_order_system_number": queryValues.Get("work_order_system_number"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, baseErr := r.WorkOrderService.GetServiceRequestByWO(workOrderId, criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetClaimByWO gets all claim by work order
// @Summary Get Claim By Work Order
// @Description Retrieve all claim by work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/claim-service/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetClaimByWO(writer http.ResponseWriter, request *http.Request) {

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order_detail.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order_detail.transaction_type_id":      queryValues.Get("transaction_type_id"),
		"trx_work_order_detail.atpm_claim_number":        queryValues.Get("atpm_claim_number"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, baseErr := r.WorkOrderService.GetClaimByWO(workOrderId, criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetClaimItemByWO gets all claim item by work order
// @Summary Get Claim Item By Work Order
// @Description Retrieve all claim item by work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/claim-item-service/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetClaimItemByWO(writer http.ResponseWriter, request *http.Request) {

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"trx_work_order_detail.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order_detail.transaction_type_id":      queryValues.Get("transaction_type_id"),
		"trx_work_order_detail.atpm_claim_number":        queryValues.Get("atpm_claim_number"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, baseErr := r.WorkOrderService.GetClaimItemByWO(workOrderId, criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetWOByBillCode gets all work order by bill code
// @Summary Get Work Order By Bill Code
// @Description Retrieve all work order by bill code
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path int true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/transactiontype-service/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetWOByBillCode(writer http.ResponseWriter, request *http.Request) {

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"trx_work_order_detail.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order_detail.transaction_type_id":      queryValues.Get("transaction_type_id"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, baseErr := r.WorkOrderService.GetWOByBillCode(workOrderId, criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetDetailWOByClaimBillCode gets all claim detail by work order
// @Summary Get Claim Detail By Work Order
// @Description Retrieve all claim detail by work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path int true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/claim-detail-service/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetDetailWOByClaimBillCode(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	transactionTypeId, err := strconv.Atoi(chi.URLParam(request, "transaction_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid transaction type ID", http.StatusBadRequest)
		return
	}

	atpmClaimNumber := chi.URLParam(request, "atpm_claim_number")
	if atpmClaimNumber == "" {
		payloads.NewHandleError(writer, "Invalid ATPM claim number", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Call the service method
	paginatedData, baseErr := r.WorkOrderService.GetDetailWOByClaimBillCode(workOrderId, transactionTypeId, atpmClaimNumber, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, paginatedData, "Get Data Successfully", http.StatusOK)
}

// GetDetailWOByBillCode gets all claim bill by work order
// @Summary Get Claim Bill By Work Order
// @Description Retrieve all claim bill by work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path int true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Param transaction_type_id query string false "Transaction Type ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/claim-bill-service/{work_order_system_number}/{transaction_type_id} [get]
func (r *WorkOrderControllerImpl) GetDetailWOByBillCode(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	transactionTypeId, err := strconv.Atoi(chi.URLParam(request, "transaction_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid transaction type ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Call the service method
	paginatedData, baseErr := r.WorkOrderService.GetDetailWOByBillCode(workOrderId, transactionTypeId, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, paginatedData, "Get Data Successfully", http.StatusOK)
}

// GetDetailWOByATPMBillCode gets all claim bill by work order
// @Summary Get Claim Bill By Work Order
// @Description Retrieve all claim bill by work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path int true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Param transaction_type_id query string false "Transaction Type ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/atpm-claim-bill-service/{work_order_system_number}/{transaction_type_id} [get]
func (r *WorkOrderControllerImpl) GetDetailWOByATPMBillCode(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	transactionTypeId, err := strconv.Atoi(chi.URLParam(request, "transaction_type_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid transaction type ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Call the service method
	paginatedData, baseErr := r.WorkOrderService.GetDetailWOByATPMBillCode(workOrderId, transactionTypeId, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, paginatedData, "Get Data Successfully", http.StatusOK)
}

// GetSupplyByWO gets all supply by work order
// @Summary Get Supply By Work Order
// @Description Retrieve all supply by work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path int true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/supply-service/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetSupplyByWO(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"trx_work_order_detail.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order_detail.transaction_type_id":      queryValues.Get("transaction_type_id"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	// Call the service method
	paginatedData, totalPages, totalRows, baseErr := r.WorkOrderService.GetSupplyByWO(workOrderId, criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}
