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

type BookingEstimationControllerImpl struct {
	bookingEstimationService transactionworkshopservice.BookingEstimationService
}

type BookingEstimationController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	New(writer http.ResponseWriter, request *http.Request)
	NewBooking(writer http.ResponseWriter, request *http.Request)
	NewAffiliated(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	CloseOrder(writer http.ResponseWriter, request *http.Request)
}

func NewBookingEstimationController(BookingEstimationService transactionworkshopservice.BookingEstimationService) BookingEstimationController {
	return &BookingEstimationControllerImpl{
		bookingEstimationService: BookingEstimationService,
	}
}

// GetAll gets all booking estimations
func (r *BookingEstimationControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	paginatedData, totalPages, totalRows, err := r.bookingEstimationService.GetAll(criteria, paginate)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *BookingEstimationControllerImpl) New(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

func (r *BookingEstimationControllerImpl) NewBooking(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

func (r *BookingEstimationControllerImpl) NewAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

// WithTrx handles the transaction for all booking estimations
func (r *BookingEstimationControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	// This function can be implemented to handle transaction-related logic if needed
	// For now, it's empty
}

// Save saves a new booking estimation
func (r *BookingEstimationControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	// Menginisialisasi koneksi database
	db := config.InitDB()

	// Mendekode payload request ke struct WorkOrderRequest
	var bookingEstimationRequest transactionworkshoppayloads.BookingEstimationRequest
	if err := json.NewDecoder(request.Body).Decode(&bookingEstimationRequest); err != nil {
		// Tangani kesalahan jika tidak dapat mendekode payload
		payloads.NewHandleError(writer, "Failed to decode request payload", http.StatusBadRequest)
		return
	}

	// Panggil fungsi Save dari layanan untuk menyimpan data booking estimation
	success, err := r.bookingEstimationService.Save(db, bookingEstimationRequest) // Memastikan untuk meneruskan db ke dalam metode Save
	if err != nil {
		// Tangani kesalahan dari layanan
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai hasil penyimpanan
	if success {
		payloads.NewHandleSuccess(writer, nil, "Work order saved successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to save booking estimation", http.StatusInternalServerError)
	}
}

func (r *BookingEstimationControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

func (r *BookingEstimationControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Cancel booking estimation
}

func (r *BookingEstimationControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close booking estimation
}
