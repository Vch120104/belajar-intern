package transactionworkshopcontroller

import (
	"after-sales/api/config"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	utils "after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkOrderControllerImpl struct {
	WorkOrderService transactionworkshopservice.WorkOrderService
}

type WorkOrderController interface {
	VehicleLookup(writer http.ResponseWriter, request *http.Request)
	CampaignLookup(writer http.ResponseWriter, request *http.Request)

	GetAllRequest(writer http.ResponseWriter, request *http.Request)
	GetRequestById(writer http.ResponseWriter, request *http.Request)
	UpdateRequest(writer http.ResponseWriter, request *http.Request)
	AddRequest(writer http.ResponseWriter, request *http.Request)
	DeleteRequest(writer http.ResponseWriter, request *http.Request)

	GetAllVehicleService(writer http.ResponseWriter, request *http.Request)
	GetVehicleServiceById(writer http.ResponseWriter, request *http.Request)
	UpdateVehicleService(writer http.ResponseWriter, request *http.Request)
	AddVehicleService(writer http.ResponseWriter, request *http.Request)
	DeleteVehicleService(writer http.ResponseWriter, request *http.Request)

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

	NewBill(writer http.ResponseWriter, request *http.Request)
	AddBill(writer http.ResponseWriter, request *http.Request)
	UpdateBill(writer http.ResponseWriter, request *http.Request)
	DeleteBill(writer http.ResponseWriter, request *http.Request)

	NewDropPoint(writer http.ResponseWriter, request *http.Request)
	// AddDropPoint(writer http.ResponseWriter, request *http.Request)
	// UpdateDropPoint(writer http.ResponseWriter, request *http.Request)
	// DeleteDropPoint(writer http.ResponseWriter, request *http.Request)

	NewVehicleBrand(writer http.ResponseWriter, request *http.Request)
	NewVehicleModel(writer http.ResponseWriter, request *http.Request)
	GenerateDocumentNumber(writer http.ResponseWriter, request *http.Request)

	GetAllDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	GetDetailByIdWorkOrder(writer http.ResponseWriter, request *http.Request)
	AddDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	UpdateDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	DeleteDetailWorkOrder(writer http.ResponseWriter, request *http.Request)

	GetAllBooking(writer http.ResponseWriter, request *http.Request)
	GetBookingById(writer http.ResponseWriter, request *http.Request)
	NewBooking(writer http.ResponseWriter, request *http.Request)
	SaveBooking(writer http.ResponseWriter, request *http.Request)
	VoidBooking(writer http.ResponseWriter, request *http.Request)
	CloseBooking(writer http.ResponseWriter, request *http.Request)

	GetAllAffiliated(writer http.ResponseWriter, request *http.Request)
	GetAffiliatedById(writer http.ResponseWriter, request *http.Request)
	NewAffiliated(writer http.ResponseWriter, request *http.Request)
	SaveAffiliated(writer http.ResponseWriter, request *http.Request)
	VoidAffiliated(writer http.ResponseWriter, request *http.Request)
	CloseAffiliated(writer http.ResponseWriter, request *http.Request)
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
		"trx_work_order.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order.work_order_type_id":       queryValues.Get("work_order_type_id"),
		"work_order_status_id":                    queryValues.Get("work_order_status_id"),
		"trx_work_order.brand_id":                 queryValues.Get("brand_id"),
		"trx_work_order.model_id":                 queryValues.Get("model_id"),
		"trx_work_order.vehicle_id":               queryValues.Get("vehicle_id"),
		"trx_work_order.work_order_date":          queryValues.Get("work_order_date"),
		"trx_work_order.work_order_close_date":    queryValues.Get("work_order_close_date"),
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

	// Menginisialisasi koneksi database
	db := config.InitDB()

	var workOrderRequest transactionworkshoppayloads.WorkOrderNormalRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, err := r.WorkOrderService.New(db, workOrderRequest)
	if err != nil {

		exceptions.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai hasil penyimpanan
	if success {
		payloads.NewHandleSuccess(writer, nil, "Work order saved successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Failed to save work order", http.StatusInternalServerError)
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

	db := config.InitDB()

	statuses, err := r.WorkOrderService.NewStatus(db, filters)
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

	db := config.InitDB()
	success, err := r.WorkOrderService.AddStatus(db, statusRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, nil, "Status added successfully", http.StatusCreated)
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

	db := config.InitDB()
	update, err := r.WorkOrderService.UpdateStatus(db, int(statusID), statusRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if update {
		payloads.NewHandleSuccess(writer, nil, "Status updated successfully", http.StatusOK)
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

	db := config.InitDB()
	delete, err := r.WorkOrderService.DeleteStatus(db, int(statusID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, nil, "Status deleted successfully", http.StatusOK)
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
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
	statuses, err := r.WorkOrderService.NewBill(db)
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

	db := config.InitDB()
	success, err := r.WorkOrderService.AddBill(db, billRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, nil, "Bill added successfully", http.StatusCreated)
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

	db := config.InitDB()
	update, err := r.WorkOrderService.UpdateBill(db, int(billID), billRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	if update {
		payloads.NewHandleSuccess(writer, nil, "Bill updated successfully", http.StatusOK)
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

	db := config.InitDB()
	delete, err := r.WorkOrderService.DeleteBill(db, int(billID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, nil, "Bill deleted successfully", http.StatusOK)
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

	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
	statuses, err := r.WorkOrderService.NewType(db, filters)
	if err != nil {
		// Menangani kesalahan dari layanan
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

	db := config.InitDB()
	success, err := r.WorkOrderService.AddType(db, typeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, nil, "Type added successfully", http.StatusCreated)
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

	db := config.InitDB()
	update, err := r.WorkOrderService.UpdateType(db, int(typeID), typeRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if update {
		payloads.NewHandleSuccess(writer, nil, "Type updated successfully", http.StatusOK)
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

	db := config.InitDB()
	delete, err := r.WorkOrderService.DeleteType(db, int(typeID))
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if delete {
		payloads.NewHandleSuccess(writer, nil, "Type deleted successfully", http.StatusOK)
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
	db := config.InitDB()

	statuses, err := r.WorkOrderService.NewDropPoint(db)
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

// // AddDropPoint adds a new drop point to a work order
// // @Summary Add Work Order Drop Point
// // @Description Add a new drop point to a work order
// // @Accept json
// // @Produce json
// // @Tags Transaction : Workshop Work Order
// // @Param reqBody body transactionworkshoppayloads.WorkOrderDropPointRequest true "Work Order Drop Point Data"
// // @Success 201 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// // @Router /v1/work-order/dropdown-drop-point [post]
// func (r *WorkOrderControllerImpl) AddDropPoint(writer http.ResponseWriter, request *http.Request) {
// 	// Add drop point to work order
// 	var dropPointRequest transactionworkshoppayloads.WorkOrderDropPointRequest
// 	helper.ReadFromRequestBody(request, &dropPointRequest)

// 	db := config.InitDB()
// 	err := r.WorkOrderService.AddDropPoint(db, dropPointRequest)
// 	if err != nil {
// 		exceptions.NewAppException(writer, request, err)
// 		return
// 	}

// 	payloads.NewHandleSuccess(writer, nil, "Drop point added successfully", http.StatusCreated)
// }

// // UpdateDropPoint updates a drop point of a work order
// // @Summary Update Work Order Drop Point
// // @Description Update a drop point of a work order
// // @Accept json
// // @Produce json
// // @Tags Transaction : Workshop Work Order
// // @Param work_order_drop_point_id path string true "Work Order Drop Point ID"
// // @Param reqBody body transactionworkshoppayloads.WorkOrderDropPointRequest true "Work Order Drop Point Data"
// // @Success 200 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// // @Router /v1/work-order/dropdown-drop-point/{work_order_drop_point_id} [put]
// func (r *WorkOrderControllerImpl) UpdateDropPoint(writer http.ResponseWriter, request *http.Request) {
// 	// Update drop point of a work order
// 	dropPointID, _ := strconv.Atoi(chi.URLParam(request, "work_order_drop_point_id"))

// 	var dropPointRequest transactionworkshoppayloads.WorkOrderDropPointRequest
// 	helper.ReadFromRequestBody(request, &dropPointRequest)

// 	db := config.InitDB()
// 	err := r.WorkOrderService.UpdateDropPoint(db, int(dropPointID), dropPointRequest)
// 	if err != nil {
// 		exceptions.NewAppException(writer, request, err)
// 		return
// 	}

// 	payloads.NewHandleSuccess(writer, nil, "Drop point updated successfully", http.StatusOK)
// }

// // DeleteDropPoint deletes a drop point from a work order
// // @Summary Delete Work Order Drop Point
// // @Description Delete a drop point from a work order
// // @Accept json
// // @Produce json
// // @Tags Transaction : Workshop Work Order
// // @Param work_order_drop_point_id path string true "Work Order Drop Point ID"
// // @Success 204 {object} payloads.Response
// // @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// // @Router /v1/work-order/dropdown-drop-point/{work_order_drop_point_id} [delete]
// func (r *WorkOrderControllerImpl) DeleteDropPoint(writer http.ResponseWriter, request *http.Request) {
// 	// Delete drop point from work order
// 	dropPointID, _ := strconv.Atoi(chi.URLParam(request, "work_order_drop_point_id"))

// 	db := config.InitDB()
// 	err := r.WorkOrderService.DeleteDropPoint(db, int(dropPointID))
// 	if err != nil {
// 		exceptions.NewAppException(writer, request, err)
// 		return
// 	}

// 	payloads.NewHandleSuccess(writer, nil, "Drop point deleted successfully", http.StatusNoContent)
// }

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
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
	statuses, err := r.WorkOrderService.NewVehicleBrand(db)
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

	db := config.InitDB()
	create, baseErr := r.WorkOrderService.NewVehicleModel(db, brandId)

	// Periksa apakah ada error yang dikembalikan
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

// VehicleLookup looks up vehicles
// @Summary Vehicle Lookup
// @Description Look up vehicles with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param vehicle_id query string false "Vehicle ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/lookup-vehicle [get]
func (r *WorkOrderControllerImpl) VehicleLookup(writer http.ResponseWriter, request *http.Request) {
	// Menginisialisasi koneksi database
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.vehicle_id": queryValues.Get("vehicle_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.VehicleLookup(criteria, paginate)
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

// CampaignLookup looks up campaign
// @Summary Campaign Lookup
// @Description Look up campaign with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param campaign_id query string false "Campaign ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/lookup-campaign [get]
func (r *WorkOrderControllerImpl) CampaignLookup(writer http.ResponseWriter, request *http.Request) {
	// Menginisialisasi koneksi database
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.campaign_id": queryValues.Get("campaign_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.CampaignLookup(criteria, paginate)
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
	// Get all services of a work order
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"trx_work_order_service.work_order_system_number": chi.URLParam(request, "work_order_system_number"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAllRequest(criteria, paginate)
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

	db := config.InitDB()
	err := r.WorkOrderService.UpdateRequest(db, int(workorderID), int(requestID), groupRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Request updated successfully", http.StatusOK)

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

	if success {
		payloads.NewHandleSuccess(writer, nil, "Request added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
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
		payloads.NewHandleSuccess(writer, nil, "Request deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
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
	// Get all vehicle services of a work order
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"trx_work_order_vehicle_service.work_order_system_number": chi.URLParam(request, "work_order_system_number"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderService.GetAllVehicleService(criteria, paginate)
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

	db := config.InitDB()
	err := r.WorkOrderService.UpdateVehicleService(db, int(workorderID), int(vehicleServiceID), vehicleRequest)

	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Vehicle service updated successfully", http.StatusOK)
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

	if success {
		payloads.NewHandleSuccess(writer, nil, "Vehicle service added successfully", http.StatusOK)
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
// @Router /v1/work-order/{work_order_system_number}/vehicleservice/{work_order_vehicle_service_id} [delete]
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
		payloads.NewHandleSuccess(writer, nil, "Vehicle service deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
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
	// Get work order by ID
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	workOrder, baseErr := r.WorkOrderService.GetById(workOrderId)
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

	// Read the request body and convert to WorkOrderRequest struct
	var workOrderRequest transactionworkshoppayloads.WorkOrderNormalSaveRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	// Initialize the database connection
	db := config.InitDB()

	// Save the work order
	success, baseErr := r.WorkOrderService.Save(db, workOrderRequest, workOrderId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	// Send response to client based on the save result
	if success {
		payloads.NewHandleSuccess(writer, nil, "Work order saved successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to save work order", http.StatusInternalServerError)
	}

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
// @Router /v1/work-order/normal/{work_order_system_number}/submit [post]
func (r *WorkOrderControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Retrieve work order ID from URL parameters
	workOrderId := chi.URLParam(request, "work_order_system_number")
	workOrderIdInt, err := strconv.Atoi(workOrderId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	// Initialize database connection
	db := config.InitDB()

	// Submit work order
	success, newDocumentNumber, baseErr := r.WorkOrderService.Submit(db, workOrderIdInt)
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

	// Handle success and failure responses
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
	workOrderId := chi.URLParam(request, "work_order_system_number")
	workOrderIdInt, err := strconv.Atoi(workOrderId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	db := config.InitDB()
	success, baseErr := r.WorkOrderService.Void(db, workOrderIdInt)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, nil, "Work order voided successfully", http.StatusOK)
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
// @Router /v1/work-order/normal/{work_order_system_number}/close [patch]
func (r *WorkOrderControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close work order
	workOrderId := chi.URLParam(request, "work_order_system_number")
	workOrderIdInt, err := strconv.Atoi(workOrderId)

	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	db := config.InitDB()
	success, baseErr := r.WorkOrderService.CloseOrder(db, workOrderIdInt)

	if baseErr != nil {
		if baseErr.Message == "Work order cannot be closed because status is draft" {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, nil, "Work order closed successfully", http.StatusOK)
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
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/detail [get]
func (r *WorkOrderControllerImpl) GetAllDetailWorkOrder(writer http.ResponseWriter, request *http.Request) {
	// Get the detail of a work order
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"trx_work_order_detail.work_order_system_number": chi.URLParam(request, "work_order_system_number"),
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
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	detailId, _ := strconv.Atoi(chi.URLParam(request, "work_order_detail_id"))

	detail, err := r.WorkOrderService.GetDetailByIdWorkOrder(int(workOrderId), int(detailId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
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

	db := config.InitDB()
	update, err := r.WorkOrderService.UpdateDetailWorkOrder(db, int(workOrderId), int(detailId), detailRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if update {
		payloads.NewHandleSuccess(writer, nil, "Detail updated successfully", http.StatusOK)
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
	// Add a new detail to a work order
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var detailRequest transactionworkshoppayloads.WorkOrderDetailRequest
	helper.ReadFromRequestBody(request, &detailRequest)

	success, err := r.WorkOrderService.AddDetailWorkOrder(int(workOrderId), detailRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, nil, "Detail added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

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
		payloads.NewHandleSuccess(writer, nil, "Detail deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
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
	// Get all work order bookings
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"trx_work_order.work_order_system_number": chi.URLParam(request, "work_order_system_number"),
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
// @Param booking_system_number path string true "Work Order Booking ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/{work_order_system_number}/{booking_system_number} [get]
func (r *WorkOrderControllerImpl) GetBookingById(writer http.ResponseWriter, request *http.Request) {
	// Get a work order booking by ID
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	workOrderBookId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	workOrder, err := r.WorkOrderService.GetBookingById(workOrderId, workOrderBookId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
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
// @Router /v1/work-order/booking/{booking_system_number}/{booking_system_number} [put]
func (r *WorkOrderControllerImpl) SaveBooking(writer http.ResponseWriter, request *http.Request) {
	// Update a work order booking
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	workOrderBookId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	var workOrderRequest transactionworkshoppayloads.WorkOrderBookingRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	db := config.InitDB()
	result, err := r.WorkOrderService.SaveBooking(db, workOrderId, workOrderBookId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order updated successfully", http.StatusOK)
}

// AddWorkOrderBooking adds a new work order booking
// @Summary Add Work Order Booking
// @Description Add a new work order booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Booking
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBookingRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) NewBooking(writer http.ResponseWriter, request *http.Request) {
	// Add a new work order booking
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var workOrderRequest transactionworkshoppayloads.WorkOrderBookingRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	db := config.InitDB()
	result, err := r.WorkOrderService.NewBooking(db, workOrderId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order added successfully", http.StatusCreated)
}

// Void WorkOrderBooking deletes a work order booking
// @Summary Void Work Order Booking
// @Description Void a work order booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Booking
// @Param work_order_system_number path string true "Work Order ID"
// @Param booking_system_number path string true "Work Order Booking ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/void/{work_order_system_number}/{booking_system_number} [delete]
func (r *WorkOrderControllerImpl) VoidBooking(writer http.ResponseWriter, request *http.Request) {
	// Void a work order booking
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	bookingId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	db := config.InitDB()
	result, err := r.WorkOrderService.VoidBooking(db, workOrderId, bookingId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order voided successfully", http.StatusOK)
}

// CloseWorkOrderBooking closes a work order booking
// @Summary Close Work Order Booking
// @Description Close a work order booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Booking
// @Param work_order_system_number path string true "Work Order ID"
// @Param booking_system_number path string true "Work Order Booking ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBookingRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/{work_order_system_number}/close/{booking_system_number} [patch]
func (r *WorkOrderControllerImpl) CloseBooking(writer http.ResponseWriter, request *http.Request) {
	// Close a work order booking
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	bookingId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	db := config.InitDB()
	close, err := r.WorkOrderService.CloseBooking(db, workOrderId, bookingId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, close, "Work order closed successfully", http.StatusOK)
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
// @Param affiliated_work_order_system_number path string true "Affiliated Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated/{work_order_system_number}/{affiliated_work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetAffiliatedById(writer http.ResponseWriter, request *http.Request) {
	// Get affiliated work order by ID
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	affiliatedWorkOrderId, _ := strconv.Atoi(chi.URLParam(request, "affiliated_work_order_system_number"))

	workOrder, err := r.WorkOrderService.GetAffiliatedById(workOrderId, affiliatedWorkOrderId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
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

	db := config.InitDB()
	result, err := r.WorkOrderService.NewAffiliated(db, workOrderId, workOrderRequest)
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

	db := config.InitDB()
	result, err := r.WorkOrderService.SaveAffiliated(db, workOrderId, affiliatedWorkOrderId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order updated successfully", http.StatusOK)
}

// Void Affiliated Work Order deletes an affiliated work order
// @Summary Void Affiliated Work Order
// @Description Void an affiliated work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Affiliated
// @Param work_order_system_number path string true "Work Order ID"
// @Param affiliated_work_order_system_number path string true "Affiliated Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderAffiliatedRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated/void/{work_order_system_number}/{affiliated_work_order_system_number} [delete]
func (r *WorkOrderControllerImpl) VoidAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Void an affiliated work order
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	affiliatedWorkOrderId, _ := strconv.Atoi(chi.URLParam(request, "affiliated_work_order_system_number"))

	db := config.InitDB()
	result, err := r.WorkOrderService.VoidAffiliated(db, workOrderId, affiliatedWorkOrderId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order voided successfully", http.StatusOK)
}

// CloseAffiliated closes an affiliated work order
// @Summary Close Affiliated Work Order
// @Description Close an affiliated work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order Affiliated
// @Param work_order_system_number path string true "Work Order ID"
// @Param affiliated_work_order_system_number path string true "Affiliated Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderAffiliatedRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated/{work_order_system_number}/close/{affiliated_work_order_system_number} [patch]
func (r *WorkOrderControllerImpl) CloseAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Close an affiliated work order
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	affiliatedWorkOrderId, _ := strconv.Atoi(chi.URLParam(request, "affiliated_work_order_system_number"))

	db := config.InitDB()
	result, err := r.WorkOrderService.CloseAffiliated(db, workOrderId, affiliatedWorkOrderId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order closed successfully", http.StatusOK)
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

	db := config.InitDB()

	result, err := r.WorkOrderService.GenerateDocumentNumber(db, workOrderId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Document number generated successfully", http.StatusOK)

}
