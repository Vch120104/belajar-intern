package transactionbodyshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionbodyshopservice "after-sales/api/services/transaction/bodyshop"
	"after-sales/api/utils"
	"fmt"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServiceBodyshopControllerImp struct {
	ServiceBodyshopService transactionbodyshopservice.ServiceBodyshopService
}

type ServiceBodyshopController interface {
	GetAllByTechnicianWOBodyshop(writer http.ResponseWriter, request *http.Request)
	StartService(writer http.ResponseWriter, request *http.Request)
	PendingService(writer http.ResponseWriter, request *http.Request)
	TransferService(writer http.ResponseWriter, request *http.Request)
	StopService(writer http.ResponseWriter, request *http.Request)
}

func NewServiceBodyshopController(service transactionbodyshopservice.ServiceBodyshopService) ServiceBodyshopController {
	return &ServiceBodyshopControllerImp{
		ServiceBodyshopService: service,
	}
}

// GetAll gets all service receipts
// @Summary Get all service receipts
// @Description Get all service receipts
// @Tags Transaction : Workshop Service Receipts
// @Accept json
// @Produce json
// @Param service_request_system_number query string false "Service receipts System Number"
// @Param service_request_id query string false "Service receipts ID"
// @Param brand_id query string false "Brand ID"
// @Param model_id query string false "Model ID"
// @Param vehicle_id query string false "Vehicle ID"
// @Param service_request_date query string false "Service receipts Date"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-receipt [get]
func (r *ServiceBodyshopControllerImp) GetAllByTechnicianWOBodyshop(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	TechnicianId, err := strconv.Atoi(chi.URLParam(request, "technician_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Technician ID", http.StatusBadRequest)
		return
	}

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryParams := map[string]string{
		"service_request_system_number": queryValues.Get("service_request_system_number"),
		"service_request_id":            queryValues.Get("service_request_id"),
		"brand_id":                      queryValues.Get("brand_id"),
		"model_id":                      queryValues.Get("model_id"),
		"vehicle_id":                    queryValues.Get("vehicle_id"),
		"service_request_date":          queryValues.Get("service_request_date"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, baseErr := r.ServiceBodyshopService.GetAllByTechnicianWOBodyshop(TechnicianId, workOrderId, criteria, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, paginatedData, "Service Retrieved Successfully", http.StatusOK)
}

// StartService starts the service
// @Summary Start the service
// @Description Start the service
// @Tags Transaction : Workshop Service Log
// @Accept json
// @Produce json
// @Param service_log_system_number query int false "Service Log System Number"
// @Param work_order_system_number query int false "Work Order System Number"
// @Param technician_allocation_system_number query int false "Allocation ID"
// @Param company_id query int false "Company ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-log/{technician_allocation_system_number}/{work_order_system_number}/{service_log_system_number}/{company_id}/start [post]
func (r *ServiceBodyshopControllerImp) StartService(writer http.ResponseWriter, request *http.Request) {
	// Extract parameters from URL
	allocId, err := strconv.Atoi(chi.URLParam(request, "technician_allocation_system_number"))
	if err != nil {

		payloads.NewHandleError(writer, "Invalid Technician Allocate ID", http.StatusBadRequest)
		return
	}

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {

		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(chi.URLParam(request, "company_id"))
	if err != nil {

		payloads.NewHandleError(writer, "Invalid company code", http.StatusBadRequest)
		return
	}

	fmt.Printf("Parameters: allocId=%d, workOrderId=%d, companyId=%d\n", allocId, workOrderId, companyId)

	// Check if ServiceBodyshopService is initialized
	if r.ServiceBodyshopService == nil {

		payloads.NewHandleError(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Call the service method
	success, baseErr := r.ServiceBodyshopService.StartService(allocId, workOrderId, companyId)
	if baseErr != nil {

		payloads.NewHandleError(writer, baseErr.Err.Error(), baseErr.StatusCode)
		return
	}

	if success {

		payloads.NewHandleSuccess(writer, success, "Service Started Successfully", http.StatusOK)
	} else {

		payloads.NewHandleError(writer, "Failed to start service", http.StatusInternalServerError)
	}
}

// PendingService pending the service
// @Summary Pending the service
// @Description Pending the service
// @Tags Transaction : Workshop Service Log
// @Accept json
// @Produce json
// @Param service_log_system_number query int false "Service Log System Number"
// @Param work_order_system_number query int false "Work Order System Number"
// @Param technician_allocation_system_number query int false "Allocation ID"
// @Param company_id query int false "Company ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-log/{technician_allocation_system_number}/{work_order_system_number}/{company_id}/pending [post]
func (r *ServiceBodyshopControllerImp) PendingService(writer http.ResponseWriter, request *http.Request) {

	alllocId, err := strconv.Atoi(chi.URLParam(request, "technician_allocation_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Technician Allocate ID", http.StatusBadRequest)
		return
	}

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(chi.URLParam(request, "company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid company code", http.StatusBadRequest)
		return
	}

	pending, baseErr := r.ServiceBodyshopService.PendingService(alllocId, workOrderId, companyId)
	if baseErr != nil {

		payloads.NewHandleError(writer, baseErr.Err.Error(), baseErr.StatusCode)
		return
	}

	if pending {

		payloads.NewHandleSuccess(writer, pending, "Service Pending Successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to pending service", http.StatusInternalServerError)
	}
}

// TransferService transfer the service
// @Summary Transfer the service
// @Description Transfer the service
// @Tags Transaction : Workshop Service Log
// @Accept json
// @Produce json
// @Param service_log_system_number query int false "Service Log System Number"
// @Param work_order_system_number query int false "Work Order System Number"
// @Param technician_allocation_system_number query int false "Allocation ID"
// @Param company_id query int false "Company ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-log/{technician_allocation_system_number}/{work_order_system_number}/{company_id}/transfer [post]

func (r *ServiceBodyshopControllerImp) TransferService(writer http.ResponseWriter, request *http.Request) {

	alllocId, err := strconv.Atoi(chi.URLParam(request, "technician_allocation_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Technician Allocate ID", http.StatusBadRequest)
		return
	}

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(chi.URLParam(request, "company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid company code", http.StatusBadRequest)
		return
	}

	transfer, baseErr := r.ServiceBodyshopService.TransferService(alllocId, workOrderId, companyId)
	if baseErr != nil {

		payloads.NewHandleError(writer, baseErr.Err.Error(), baseErr.StatusCode)
		return
	}

	if transfer {

		payloads.NewHandleSuccess(writer, transfer, "Service Transfer Successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to transfer service", http.StatusInternalServerError)
	}
}

// StopService stops the service
// @Summary Stop the service
// @Description Stop the service
// @Tags Transaction : Workshop Service Log
// @Accept json
// @Produce json
// @Param service_log_system_number query int false "Service Log System Number"
// @Param work_order_system_number query int false "Work Order System Number"
// @Param technician_allocation_system_number query int false "Allocation ID"
// @Param company_id query int false "Company ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-log/{technician_allocation_system_number}/{work_order_system_number}/{company_id}/stop [post]
func (r *ServiceBodyshopControllerImp) StopService(writer http.ResponseWriter, request *http.Request) {

	alllocId, err := strconv.Atoi(chi.URLParam(request, "technician_allocation_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Technician Allocate ID", http.StatusBadRequest)
		return
	}

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	companyId, err := strconv.Atoi(chi.URLParam(request, "company_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid company code", http.StatusBadRequest)
		return
	}

	stop, baseErr := r.ServiceBodyshopService.StopService(alllocId, workOrderId, companyId)
	if baseErr != nil {

		payloads.NewHandleError(writer, baseErr.Err.Error(), baseErr.StatusCode)
		return
	}

	if stop {

		payloads.NewHandleSuccess(writer, stop, "Service Stop Successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to stop service", http.StatusInternalServerError)
	}
}
