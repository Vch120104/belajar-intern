package transactionworkshopcontroller

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
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
func (r *WorkOrderControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.brand_id": queryValues.Get("brand_id"),
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

func (r *WorkOrderControllerImpl) New(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
}

func (r *WorkOrderControllerImpl) NewBooking(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
}

func (r *WorkOrderControllerImpl) NewAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
}
func (r *WorkOrderControllerImpl) NewStatus(writer http.ResponseWriter, request *http.Request) {
	// Mengambil data entity WorkOrderMasterStatus dari body request
	var workOrderMasterStatus transactionworkshopentities.WorkOrderMasterStatus
	if err := json.NewDecoder(request.Body).Decode(&workOrderMasterStatus); err != nil {
		// Menangani kesalahan jika tidak dapat mendekode payload
		payloads.NewHandleError(writer, "Failed to decode request payload", http.StatusBadRequest)
		return
	}

	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Panggil fungsi NewStatus dari layanan untuk menyimpan data work order
	success, err := r.WorkOrderService.NewStatus(db, workOrderMasterStatus)
	if err != nil {
		// Menangani kesalahan dari layanan
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai hasil penyimpanan
	if success {
		payloads.NewHandleSuccess(writer, nil, "Work order status updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to update work order status", http.StatusInternalServerError)
	}
}

// WithTrx handles the transaction for all work orders
func (r *WorkOrderControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	// This function can be implemented to handle transaction-related logic if needed
	// For now, it's empty
}

// Save saves a new work order
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

func (r *WorkOrderControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Create new work order
}

func (r *WorkOrderControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Cancel work order
}

func (r *WorkOrderControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close work order
}
