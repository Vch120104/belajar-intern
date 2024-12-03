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
	"log"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ServiceReceiptControllerImp struct {
	ServiceReceiptService transactionworkshopservice.ServiceReceiptService
}

type ServiceReceiptController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
}

func NewServiceReceiptController(service transactionworkshopservice.ServiceReceiptService) ServiceReceiptController {
	return &ServiceReceiptControllerImp{
		ServiceReceiptService: service,
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
func (r *ServiceReceiptControllerImp) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"service_request_system_number":   queryValues.Get("service_request_system_number"),
		"service_request_document_number": queryValues.Get("service_request_document_number"),
		"service_request_status_id":       queryValues.Get("service_request_status_id"),
		"brand_id":                        queryValues.Get("brand_id"),
		"model_id":                        queryValues.Get("model_id"),
		"vehicle_id":                      queryValues.Get("vehicle_id"),
		"service_request_date":            queryValues.Get("service_request_date"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.ServiceReceiptService.GetAll(criteria, paginate)
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

// GetById gets service receipt by id
// @Summary Get service receipt by id
// @Description Get service receipt by id
// @Tags Transaction : Workshop Service Receipts
// @Accept json
// @Produce json
// @Param id path string true "Service receipt ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-receipt/{service_request_system_number} [get]
func (r *ServiceReceiptControllerImp) GetById(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	serviceRequest, baseErr := r.ServiceReceiptService.GetById(ServiceRequestId, paginate)
	if baseErr != nil {
		// Log the error for debugging
		log.Printf("Error retrieving service receipt: %v", baseErr)

		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Service request not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, serviceRequest, "Get Data Successfully", http.StatusOK)
}

// Save saves service receipt
// @Summary Save service receipt
// @Description Save service receipt
// @Tags Transaction : Workshop Service Receipts
// @Accept json
// @Produce json
// @Param service_request_system_number path int true "Service Request ID"
// @Param request body transactionworkshoppayloads.ServiceReceiptSaveDataRequest true "Service receipt request"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-receipt/{service_request_system_number} [put]
func (r *ServiceReceiptControllerImp) Save(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	var serviceReceiptSaveRequest transactionworkshoppayloads.ServiceReceiptSaveDataRequest
	helper.ReadFromRequestBody(request, &serviceReceiptSaveRequest)
	if validationErr := validation.ValidationForm(writer, request, &serviceReceiptSaveRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	savedData, baseErr := r.ServiceReceiptService.Save(ServiceRequestId, serviceReceiptSaveRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, savedData, "Data saved successfully", http.StatusOK)
}
