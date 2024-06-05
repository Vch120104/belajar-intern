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
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkOrderControllerImpl struct {
	WorkOrderService transactionworkshopservice.WorkOrderService
}

type WorkOrderController interface {
	New(writer http.ResponseWriter, request *http.Request)
	NewBooking(writer http.ResponseWriter, request *http.Request)
	NewAffiliated(writer http.ResponseWriter, request *http.Request)

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
	Save(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	CloseOrder(writer http.ResponseWriter, request *http.Request)

	NewStatus(writer http.ResponseWriter, request *http.Request)
	NewType(writer http.ResponseWriter, request *http.Request)
	NewBill(writer http.ResponseWriter, request *http.Request)
	NewDropPoint(writer http.ResponseWriter, request *http.Request)
	NewVehicleBrand(writer http.ResponseWriter, request *http.Request)
	NewVehicleModel(writer http.ResponseWriter, request *http.Request)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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
// @Tags Transaction : Workshop Work Order
// @Param reqBody body transactionworkshoppayloads.WorkOrderRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal [post]
func (r *WorkOrderControllerImpl) New(writer http.ResponseWriter, request *http.Request) {

	// Menginisialisasi koneksi database
	db := config.InitDB()

	var workOrderRequest transactionworkshoppayloads.WorkOrderRequest
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

// NewBooking creates a new booking work order
// @Summary Create New Booking Work Order
// @Description Create a new booking work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/dropdown-status [get]
func (r *WorkOrderControllerImpl) NewStatus(writer http.ResponseWriter, request *http.Request) {
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
	statuses, err := r.WorkOrderService.NewStatus(db)
	if err != nil {
		// Menangani kesalahan dari layanan
		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of work order statuses", http.StatusOK)
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
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
	statuses, err := r.WorkOrderService.NewType(db)
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
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi GetAll dari layanan untuk mendapatkan semua status work order
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
// @Tags Transaction : Workshop Work Order
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
// @Tags Transaction : Workshop Work Order
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
// @Tags Transaction : Workshop Work Order
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
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/{work_order_system_number}/requestservice [post]
func (r *WorkOrderControllerImpl) AddRequest(writer http.ResponseWriter, request *http.Request) {
	// Add request to work order\
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var groupRequest transactionworkshoppayloads.WorkOrderServiceRequest
	helper.ReadFromRequestBody(request, &groupRequest)

	if err := r.WorkOrderService.AddRequest(int(workorderID), groupRequest); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Request added successfully", http.StatusOK)
}

// DeleteRequest deletes a request from a work order
// @Summary Delete Request from Work Order
// @Description Delete a request from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_id path string true "Work Order Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/{work_order_system_number}/requestservice/{work_order_service_id} [delete]
func (r *WorkOrderControllerImpl) DeleteRequest(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	requestID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_id"))

	if err := r.WorkOrderService.DeleteRequest(int(workorderID), int(requestID)); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Request deleted successfully", http.StatusOK)
}

// GetAllVehicleService gets all vehicle services of a work order
// @Summary Get All Vehicle Services of Work Order
// @Description Retrieve all vehicle services of a work order
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
// @Tags Transaction : Workshop Work Order
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
// @Tags Transaction : Workshop Work Order
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
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceVehicleRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/{work_order_system_number}/vehicleservice [post]
func (r *WorkOrderControllerImpl) AddVehicleService(writer http.ResponseWriter, request *http.Request) {
	// Add vehicle service to work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var vehicleRequest transactionworkshoppayloads.WorkOrderServiceVehicleRequest
	helper.ReadFromRequestBody(request, &vehicleRequest)

	if err := r.WorkOrderService.AddVehicleService(int(workorderID), vehicleRequest); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Vehicle service added successfully", http.StatusOK)
}

// DeleteVehicleService deletes a vehicle service from a work order
// @Summary Delete Vehicle Service from Work Order
// @Description Delete a vehicle service from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_vehicle_service_id path string true "Work Order Vehicle Service ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/{work_order_system_number}/vehicleservice/{work_order_vehicle_service_id} [delete]
func (r *WorkOrderControllerImpl) DeleteVehicleService(writer http.ResponseWriter, request *http.Request) {
	// Delete vehicle service from work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	vehicleServiceID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_vehicle_id"))

	if err := r.WorkOrderService.DeleteVehicleService(int(workorderID), int(vehicleServiceID)); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Vehicle service deleted successfully", http.StatusOK)
}

// GetById handles the transaction for all work orders
// @Summary Get Work Order By ID
// @Description Retrieve work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Param work_order_system_number path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number} [get]
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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
		exceptions.NewAppException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/{work_order_id}/close [put]
func (r *WorkOrderControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close work order
}
