package transactionworkshopcontroller

import (
	"after-sales/api/config"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"encoding/json"
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
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// New creates a new booking estimation
// @Summary Create New Booking Estimation
// @Description Create a new booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation [post]
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
// @Param work_order_system_number path string true "Booking Estimation ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/find/{work_order_system_number} [get]
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
// @Param reqBody body transactionworkshoppayloads.BookingEstimationRequest true "Booking Estimation Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation [put]
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
	if _, err := r.bookingEstimationService.Save(db, bookingEstimationRequest); err != nil {
		// Tangani kesalahan dari layanan
		exceptions.NewAppException(writer, request, err)
		return
	}

	// Kirim respons ke klien sesuai hasil penyimpanan
	payloads.NewHandleSuccess(writer, nil, "Work order saved successfully", http.StatusOK)
}

// Submit submits a new booking estimation
// @Summary Submit Booking Estimation
// @Description Submit a new booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/submit [post]
func (r *BookingEstimationControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Create new booking estimation
}

// Void cancels a booking estimation
// @Summary Cancel Booking Estimation
// @Description Cancel an existing booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param booking_estimation_id path string true "Booking Estimation ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/{booking_estimation_id} [delete]
func (r *BookingEstimationControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Cancel booking estimation
}

// CloseOrder closes a booking estimation
// @Summary Close Booking Estimation
// @Description Close an existing booking estimation
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Booking Estimation
// @Param booking_estimation_id path string true "Booking Estimation ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/booking-estimation/{booking_estimation_id}/close [put]
func (r *BookingEstimationControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close booking estimation
}


func (r *BookingEstimationControllerImpl) SaveBookEstimReq(writer http.ResponseWriter, request *http.Request){
	var formrequest transactionworkshoppayloads.BookEstimRemarkRequest
	helper.ReadFromRequestBody(request,&formrequest)
	BookingEstimationId,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	create,err := r.bookingEstimationService.SaveBookEstimReq(formrequest,BookingEstimationId)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) UpdateBookEstimReq(writer http.ResponseWriter, request *http.Request){
	var formrequest transactionworkshoppayloads.BookEstimRemarkRequest
	helper.ReadFromRequestBody(request,&formrequest)
	BookingEstimationId,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	update,err := r.bookingEstimationService.UpdateBookEstimReq(formrequest,BookingEstimationId)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl)GetByIdBookEstimReq(writer http.ResponseWriter, request *http.Request){
	bookingestimationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	get,err:= r.bookingEstimationService.GetByIdBookEstimReq(bookingestimationid)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) GetAllBookEstimReq(writer http.ResponseWriter,request *http.Request){
	queryValues := request.URL.Query()
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	bookestimid,_ := strconv.Atoi(chi.URLParam(request,"book_estim_id"))
	get,err := r.bookingEstimationService.GetAllBookEstimReq(&pagination,bookestimid)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

func ( r *BookingEstimationControllerImpl) SaveBookEstimReminderServ(writer http.ResponseWriter, request *http.Request){
	var formrequest transactionworkshoppayloads.ReminderServicePost
	helper.ReadFromRequestBody(request,&formrequest)
	bookestimid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	create,err:= r.bookingEstimationService.SaveBookEstimReminderServ(formrequest,bookestimid)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl)SaveDetailBookEstim(writer http.ResponseWriter, request *http.Request){
	var formrequest transactionworkshoppayloads.BookEstimDetailReq
	helper.ReadFromRequestBody(request,&formrequest)

	create,err := r.bookingEstimationService.SaveDetailBookEstim(formrequest)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl)AddPackage(writer http.ResponseWriter, request *http.Request){
	bookingestiomationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	packageid,_ := strconv.Atoi(chi.URLParam(request,"package_id"))
	create,err := r.bookingEstimationService.AddPackage(bookingestiomationid,packageid)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) AddContractService(writer http.ResponseWriter, request *http.Request){
	bookingestiomationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	contractserviceid,_ := strconv.Atoi(chi.URLParam(request,"contract_service_id"))
	create,err := r.bookingEstimationService.AddPackage(bookingestiomationid,contractserviceid)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) InputDiscount(writer http.ResponseWriter, request *http.Request){
	var formrequest transactionworkshoppayloads.BookEstimationPayloadsDiscount
	bookingestiomationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	helper.ReadFromRequestBody(request,&formrequest)
	create,err := r.bookingEstimationService.InputDiscount(bookingestiomationid,formrequest)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) AddFieldAction(writer http.ResponseWriter, request *http.Request){
	bookingestiomationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	idfieldaction,_ := strconv.Atoi(chi.URLParam(request,"field_action_id"))
	create,err:= r.bookingEstimationService.AddFieldAction(bookingestiomationid,idfieldaction)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) GetByIdBookEstimDetail(writer http.ResponseWriter, request *http.Request){
	bookingestiomationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	LineTypeId,_:= strconv.Atoi(chi.URLParam(request,"line_type_id"))
	get,err := r.bookingEstimationService.GetByIdBookEstimDetail(bookingestiomationid,LineTypeId)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) PostBookingEstimationCalculation(writer http.ResponseWriter, request *http.Request){
	bookingestiomationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	create,err := r.bookingEstimationService.PostBookingEstimationCalculation(bookingestiomationid)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, create, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) PutBookingEstimationCalculation(writer http.ResponseWriter, request *http.Request){
	var formrequest transactionworkshoppayloads.BookingEstimationCalculationPayloads
	bookingestiomationid,_ := strconv.Atoi(chi.URLParam(request,"booking_estimation_id"))
	LineTypeId,_:=strconv.Atoi(chi.URLParam(request,"line_type_id"))
	helper.ReadFromRequestBody(request,&formrequest)
	update,err := r.bookingEstimationService.PutBookingEstimationCalculation(bookingestiomationid,LineTypeId,formrequest)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Get Data Successfully!", http.StatusOK)
}

func (r *BookingEstimationControllerImpl) SaveBookingEstimationFromPDI(writer http.ResponseWriter, request *http.Request){
	pdisystemnumber,_ := strconv.Atoi(chi.URLParam(request,"pdi_system_number"))
	save,err:= r.bookingEstimationService.SaveBookingEstimationFromPDI(pdisystemnumber)
	if err != nil{
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, save, "Get Data Successfully!", http.StatusOK)
}