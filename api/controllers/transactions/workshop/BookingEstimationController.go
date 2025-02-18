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
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	SaveBookEstimReq(writer http.ResponseWriter, request *http.Request)
	UpdateBookEstimReq(writer http.ResponseWriter, request *http.Request)
	GetByIdBookEstimReq(writer http.ResponseWriter, request *http.Request)
	GetAllBookEstimReq(writer http.ResponseWriter, request *http.Request)
	SaveBookEstimReminderServ(writer http.ResponseWriter, request *http.Request)
	SaveDetailBookEstim(writer http.ResponseWriter, request *http.Request)
	AddPackage(writer http.ResponseWriter, request *http.Request)
	AddContractService(writer http.ResponseWriter, request *http.Request)
	InputDiscount(writer http.ResponseWriter, request *http.Request)
	CopyFromHistory(writer http.ResponseWriter, request *http.Request)
	AddFieldAction(writer http.ResponseWriter, request *http.Request)
	GetByIdBookEstimDetail(writer http.ResponseWriter, request *http.Request)
	PostBookingEstimationCalculation(writer http.ResponseWriter, request *http.Request)
	SaveBookingEstimationFromPDI(writer http.ResponseWriter, request *http.Request)
	SaveBookingEstimationFromServiceRequest(writer http.ResponseWriter, request *http.Request)
	SaveBookingEstimationAllocation(Writer http.ResponseWriter, request *http.Request)
}

func NewBookingEstimationController(BookingEstimationService transactionworkshopservice.BookingEstimationService) BookingEstimationController {
	return &BookingEstimationControllerImpl{
		bookingEstimationService: BookingEstimationService,
	}
}

// GetAll gets all booking estimations
// @Summary Get All Booking Estimations
// @Description Retrieve all booking estimations with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param brand_id query string false "Brand ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation [get]
func (r *BookingEstimationControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"brand_id":                 queryValues.Get("brand_id"),
		"model_id":                 queryValues.Get("model_id"),
		"booking_system_number":    queryValues.Get("booking_system_number"),
		"estimation_system_number": queryValues.Get("estimation_system_number"),
		"vehicle_id":               queryValues.Get("vehicle_id"),
		"document_status_id":       queryValues.Get("document_status_id"),
		"contract_person_name":     queryValues.Get("contract_person_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.bookingEstimationService.GetAll(criteria, paginate)
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

// New creates a new booking estimation
// @Summary Create New Booking Estimation
// @Description Create a new booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/normal [post]
func (r *BookingEstimationControllerImpl) New(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

// NewBooking creates a new booking estimation for booking
// @Summary Create New Booking Estimation
// @Description Create a new booking estimation for booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/booking [post]
func (r *BookingEstimationControllerImpl) NewBooking(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

// NewAffiliated creates a new affiliated booking estimation
// @Summary Create New Affiliated Booking Estimation
// @Description Create a new affiliated booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/affiliated [post]
func (r *BookingEstimationControllerImpl) NewAffiliated(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

// GetById retrieves a booking estimation by ID
// @Summary Get Booking Estimation By ID
// @Description Retrieve a booking estimation by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param batch_system_number path int true "Booking Estimation ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/normal/{batch_system_number} [get]
func (r *BookingEstimationControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	bookestimid, _ := strconv.Atoi(chi.URLParam(request, "batch_system_number"))

	result, err := r.bookingEstimationService.GetById(bookestimid)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// Save saves a new booking estimation
// @Summary Save New Booking Estimation
// @Description Save a new booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param batch_system_number path int true "Booking Estimation ID"
// @Param reqBody body transactionworkshoppayloads.BookingEstimationRequest true "Booking Estimation Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/normal/{batch_system_number} [put]
func (r *BookingEstimationControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	// Get the Work Order ID from URL parameters and convert to int
	batchSystemNumberIdStr := chi.URLParam(request, "batch_system_number")
	batchSystemNumberId, err := strconv.Atoi(batchSystemNumberIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var BookingEstimationRequest transactionworkshoppayloads.BookingEstimationRequest
	helper.ReadFromRequestBody(request, &BookingEstimationRequest)
	if validationErr := validation.ValidationForm(writer, request, &BookingEstimationRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, baseErr := r.bookingEstimationService.Save(BookingEstimationRequest, batchSystemNumberId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Booking saved successfully", http.StatusOK)

}

// Submit submits a new booking estimation
// @Summary Submit Booking Estimation
// @Description Submit a new booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param batch_system_number path int true "Booking Estimation ID"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/submit/{batch_system_number} [post]
func (r *BookingEstimationControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

// Void cancels a booking estimation
// @Summary Cancel Booking Estimation
// @Description Cancel an existing booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param batch_system_number path int true "Booking Estimation ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/void/{batch_system_number} [delete]
func (r *BookingEstimationControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Cancel booking estimation
	batchSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "batch_system_number"))
	delete, err := r.bookingEstimationService.Void(batchSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, delete, "Get Data Successfully!", http.StatusOK)
}

// CloseOrder closes a booking estimation
// @Summary Close Booking Estimation
// @Description Close an existing booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param batch_system_number path int true "Booking Estimation ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/close/{batch_system_number} [patch]
func (r *BookingEstimationControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close booking estimation
}

func (r *BookingEstimationControllerImpl) SaveBookEstimReq(writer http.ResponseWriter, request *http.Request) {
	var formrequest transactionworkshoppayloads.BookEstimRemarkRequest
	helper.ReadFromRequestBody(request, &formrequest)
	if validationErr := validation.ValidationForm(writer, request, &formrequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	BookingEstimationId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))
	create, err := r.bookingEstimationService.SaveBookEstimReq(formrequest, BookingEstimationId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) UpdateBookEstimReq(writer http.ResponseWriter, request *http.Request) {
	var formrequest transactionworkshoppayloads.BookEstimRemarkRequest
	helper.ReadFromRequestBody(request, &formrequest)
	if validationErr := validation.ValidationForm(writer, request, &formrequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	BookingEstimationRequestId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))
	update, err := r.bookingEstimationService.UpdateBookEstimReq(formrequest, BookingEstimationRequestId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) GetByIdBookEstimReq(writer http.ResponseWriter, request *http.Request) {
	bookingestimationrequestid, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))
	get, err := r.bookingEstimationService.GetByIdBookEstimReq(bookingestimationrequestid)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) GetAllBookEstimReq(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"booking_estimation_request_id":   queryValues.Get("booking_estimation_request_id"),
		"booking_estimation_request_code": queryValues.Get("booking_estimation_request_code"),
		"booking_system_number":           queryValues.Get("booking_system_number"),
		"booking_service_request":         queryValues.Get("booking_service_request"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	bookestimid, err := strconv.Atoi(chi.URLParam(request, "booking_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid booking system number", http.StatusBadRequest)
		return
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, baseErr := r.bookingEstimationService.GetAllBookEstimReq(criteria, pagination, bookestimid)
	if baseErr != nil {
		exceptions.NewNotFoundException(writer, request, baseErr)
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

func (r *BookingEstimationControllerImpl) SaveBookEstimReminderServ(writer http.ResponseWriter, request *http.Request) {
	var formrequest transactionworkshoppayloads.ReminderServicePost
	helper.ReadFromRequestBody(request, &formrequest)
	if validationErr := validation.ValidationForm(writer, request, &formrequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	bookestimid, _ := strconv.Atoi(chi.URLParam(request, "booking_estimation_id"))
	create, err := r.bookingEstimationService.SaveBookEstimReminderServ(formrequest, bookestimid)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) SaveDetailBookEstim(writer http.ResponseWriter, request *http.Request) {
	var formrequest transactionworkshoppayloads.BookEstimDetailReq
	helper.ReadFromRequestBody(request, &formrequest)
	if validationErr := validation.ValidationForm(writer, request, &formrequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	id, _ := strconv.Atoi(chi.URLParam(request, "estimation_system_number"))

	create, err := r.bookingEstimationService.SaveDetailBookEstim(formrequest, id)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) AddPackage(writer http.ResponseWriter, request *http.Request) {
	bookingestiomationid, _ := strconv.Atoi(chi.URLParam(request, "booking_estimation_id"))
	packageid, _ := strconv.Atoi(chi.URLParam(request, "package_id"))
	create, err := r.bookingEstimationService.AddPackage(bookingestiomationid, packageid)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) AddContractService(writer http.ResponseWriter, request *http.Request) {
	bookingestiomationid, _ := strconv.Atoi(chi.URLParam(request, "booking_estimation_id"))
	contractserviceid, _ := strconv.Atoi(chi.URLParam(request, "contract_service_id"))
	create, err := r.bookingEstimationService.AddContractService(bookingestiomationid, contractserviceid)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) InputDiscount(writer http.ResponseWriter, request *http.Request) {
	var formrequest transactionworkshoppayloads.BookEstimationPayloadsDiscount
	bookingestiomationid, _ := strconv.Atoi(chi.URLParam(request, "booking_estimation_id"))
	helper.ReadFromRequestBody(request, &formrequest)
	if validationErr := validation.ValidationForm(writer, request, &formrequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	create, err := r.bookingEstimationService.InputDiscount(bookingestiomationid, formrequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) AddFieldAction(writer http.ResponseWriter, request *http.Request) {
	bookingestiomationid, _ := strconv.Atoi(chi.URLParam(request, "booking_estimation_id"))
	idfieldaction, _ := strconv.Atoi(chi.URLParam(request, "field_action_id"))
	create, err := r.bookingEstimationService.AddFieldAction(bookingestiomationid, idfieldaction)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) GetByIdBookEstimDetail(writer http.ResponseWriter, request *http.Request) {
	bookingestiomationid, _ := strconv.Atoi(chi.URLParam(request, "booking_estimation_id"))
	LineTypeId, _ := strconv.Atoi(chi.URLParam(request, "line_type_id"))
	get, err := r.bookingEstimationService.GetByIdBookEstimDetail(bookingestiomationid, LineTypeId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) PostBookingEstimationCalculation(writer http.ResponseWriter, request *http.Request) {
	bookingestiomationid, _ := strconv.Atoi(chi.URLParam(request, "booking_estimation_id"))
	create, err := r.bookingEstimationService.PostBookingEstimationCalculation(bookingestiomationid)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) SaveBookingEstimationFromPDI(writer http.ResponseWriter, request *http.Request) {
	var formrequest transactionworkshoppayloads.PdiServiceRequest
	helper.ReadFromRequestBody(request, &formrequest)
	if validationErr := validation.ValidationForm(writer, request, &formrequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	pdisystemnumber, _ := strconv.Atoi(chi.URLParam(request, "pdi_system_number"))

	save, err := r.bookingEstimationService.SaveBookingEstimationFromPDI(pdisystemnumber, formrequest)
	if err != nil || !save {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("data not found"),
		})
		return
	}
	payloads.NewHandleSuccess(writer, save, "Save Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) SaveBookingEstimationFromServiceRequest(writer http.ResponseWriter, request *http.Request) {
	var formrequest transactionworkshoppayloads.PdiServiceRequest
	helper.ReadFromRequestBody(request, &formrequest)
	if validationErr := validation.ValidationForm(writer, request, &formrequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	serviceRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "service_request_system_number"))
	save, err := r.bookingEstimationService.SaveBookingEstimationFromServiceRequest(serviceRequestSystemNumber, formrequest)
	if err != nil || !save {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errors.New("data not found"),
		})
		return
	}
	payloads.NewHandleSuccess(writer, save, "Save Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) CopyFromHistory(writer http.ResponseWriter, request *http.Request) {
	BatchSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "batch_system_number"))
	save, err := r.bookingEstimationService.CopyFromHistory(BatchSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, save, "Save Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) SaveBookingEstimationAllocation(Writer http.ResponseWriter, request *http.Request) {
	BatchSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "batch_system_number"))
	var allocationpayload transactionworkshoppayloads.BookEstimationAllocation
	helper.ReadFromRequestBody(request, &allocationpayload)
	if validationErr := validation.ValidationForm(Writer, request, &allocationpayload); validationErr != nil {
		exceptions.NewBadRequestException(Writer, request, validationErr)
		return
	}

	save, err := r.bookingEstimationService.SaveBookingEstimationAllocation(BatchSystemNumber, allocationpayload)
	if err != nil {
		exceptions.NewNotFoundException(Writer, request, err)
		return
	}
	payloads.NewHandleSuccess(Writer, save, "Save Data Successfully!", http.StatusOK)
}
