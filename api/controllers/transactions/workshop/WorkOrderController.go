package transactionworkshopcontroller

import (
	"after-sales/api/config"
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"encoding/json"
	"net/http"
)

type WorkOrderControllerImpl struct {
	WorkOrderService transactionworkshopservice.WorkOrderService
}

type WorkOrderController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	New(writer http.ResponseWriter, request *http.Request)
	NewBooking(writer http.ResponseWriter, request *http.Request)
	NewAffiliated(writer http.ResponseWriter, request *http.Request)
	NewStatus(writer http.ResponseWriter, request *http.Request)
	NewType(writer http.ResponseWriter, request *http.Request)
	VehicleLookup(writer http.ResponseWriter, request *http.Request)
	CampaignLookup(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	CloseOrder(writer http.ResponseWriter, request *http.Request)
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
// @Tags Transaction : Workshop Work Order
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order [get]
func (r *WorkOrderControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order.work_order_type_id":       queryValues.Get("work_order_type_id"),
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
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Tags Transaction : Workshop Work Order
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/normal [post]
func (r *WorkOrderControllerImpl) New(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi New dari layanan untuk mengisi data work order baru
	Create, err := r.WorkOrderService.New(db)
	if err != nil {
		// Menangani kesalahan dari layanan
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai dengan hasil pengambilan status
	payloads.NewHandleSuccess(writer, Create, "work order created", http.StatusCreated)

}

// NewBooking creates a new booking work order
// @Summary Create New Booking Work Order
// @Description Create a new booking work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/booking [post]
func (r *WorkOrderControllerImpl) NewBooking(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
}

// NewAffiliated creates a new affiliated work order
// @Summary Create New Affiliated Work Order
// @Description Create a new affiliated work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/affiliated [post]
func (r *WorkOrderControllerImpl) NewAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
}

// NewStatus gets the status of new work orders
// @Summary Get Work Order Statuses
// @Description Retrieve all work order statuses
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/dropdown-status [get]
func (r *WorkOrderControllerImpl) NewStatus(writer http.ResponseWriter, request *http.Request) {
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
	statuses, err := r.WorkOrderService.NewStatus(db)
	if err != nil {
		// Menangani kesalahan dari layanan
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai dengan hasil pengambilan status
	payloads.NewHandleSuccess(writer, statuses, "List of work order statuses", http.StatusOK)
}

// NewType gets the types of new work orders
// @Summary Get Work Order Types
// @Description Retrieve all work order types
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/dropdown-type [get]
func (r *WorkOrderControllerImpl) NewType(writer http.ResponseWriter, request *http.Request) {
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
	statuses, err := r.WorkOrderService.NewType(db)
	if err != nil {
		// Menangani kesalahan dari layanan
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai dengan hasil pengambilan status
	payloads.NewHandleSuccess(writer, statuses, "List of work order type", http.StatusOK)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
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
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
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
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetById handles the transaction for all work orders
// @Summary Get Work Order By ID
// @Description Retrieve work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/find/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	// This function can be implemented to handle transaction-related logic if needed
	// For now, it's empty
}

// Save saves a new work order
// @Summary Save New Work Order
// @Description Save a new work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order [put]
func (r *WorkOrderControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Mendekode payload request ke struct WorkOrderRequest
	var workOrderRequest transactionworkshoppayloads.WorkOrderRequest
	if err := json.NewDecoder(request.Body).Decode(&workOrderRequest); err != nil {
		// Tangani kesalahan jika tidak dapat mendekode payload
		payloads.NewHandleError(writer, "Failed to decode request payload", http.StatusBadRequest)
		return
	}

	// Panggil fungsi Save dari layanan untuk menyimpan data work order
	success, err := r.WorkOrderService.Save(db, workOrderRequest) // Memastikan untuk meneruskan db ke dalam metode Save
	if err != nil {
		// Tangani kesalahan dari layanan
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai hasil penyimpanan
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
// @Tags Transaction : Workshop Work Order
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/submit [post]
func (r *WorkOrderControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
}

// Void cancels a work order
// @Summary Cancel Work Order
// @Description Cancel an existing work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_id path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/{work_order_id} [delete]
func (r *WorkOrderControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Cancel work order
}

// CloseOrder closes a work order
// @Summary Close Work Order
// @Description Close an existing work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_id path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/work-order/{work_order_id}/close [put]
func (r *WorkOrderControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close work order
}
