package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ServiceRequestControllerImp struct {
	ServiceRequestService transactionworkshopservice.ServiceRequestService
}

type ServiceRequestController interface {
	GenerateDocumentNumberServiceRequest(writer http.ResponseWriter, request *http.Request)
	NewStatus(writer http.ResponseWriter, request *http.Request)

	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	New(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	CloseOrder(writer http.ResponseWriter, request *http.Request)

	GetAllServiceDetail(writer http.ResponseWriter, request *http.Request)
	GetServiceDetailById(writer http.ResponseWriter, request *http.Request)
	AddServiceDetail(writer http.ResponseWriter, request *http.Request)
	UpdateServiceDetail(writer http.ResponseWriter, request *http.Request)
	DeleteServiceDetail(writer http.ResponseWriter, request *http.Request)
	DeleteServiceDetailMultiId(writer http.ResponseWriter, request *http.Request)
}

func NewServiceRequestController(service transactionworkshopservice.ServiceRequestService) ServiceRequestController {
	return &ServiceRequestControllerImp{
		ServiceRequestService: service,
	}
}

// GenerateDocumentNumberServiceRequest generates document number for service request
// @Summary Generate document number for service request
// @Description Generate document number for service request
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param service_request_system_number path int true "Service Request ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/document-number/{service_request_system_number} [post]
func (r *ServiceRequestControllerImp) GenerateDocumentNumberServiceRequest(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	documentNumber, baseErr := r.ServiceRequestService.GenerateDocumentNumberServiceRequest(ServiceRequestId)
	if baseErr != nil {

		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Service request not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	switch documentNumber {
	case "":
		payloads.NewHandleError(writer, "Failed to generate document number", http.StatusInternalServerError)
	case "Document number has already been generated":
		payloads.NewHandleError(writer, documentNumber, http.StatusConflict)
	default:
		payloads.NewHandleSuccess(writer, documentNumber, "Generate Document Number Successfully", http.StatusOK)
	}
}

// GetAll gets all service request
// @Summary Get all service request
// @Description Get all service request
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param service_request_system_number query string false "Service Request System Number"
// @Param service_request_id query string false "Service Request ID"
// @Param brand_id query string false "Brand ID"
// @Param model_id query string false "Model ID"
// @Param vehicle_id query string false "Vehicle ID"
// @Param service_request_date query string false "Service Request Date"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request [get]
func (r *ServiceRequestControllerImp) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

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

	paginatedData, totalPages, totalRows, err := r.ServiceRequestService.GetAll(criteria, paginate)
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

// GetById gets service request by id
// @Summary Get service request by id
// @Description Get service request by id
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param service_request_system_number path int true "Service Request ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/{service_request_system_number} [get]
func (r *ServiceRequestControllerImp) GetById(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	serviceRequest, baseErr := r.ServiceRequestService.GetById(ServiceRequestId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Service request not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, serviceRequest, "Get Data Successfully", http.StatusOK)
}

// New creates new service request
// @Summary Create new service request
// @Description Create new service request
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param reqBody body transactionworkshoppayloads.ServiceRequestSaveRequest true "Service Request Data"
// @Success 201 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request [post]
func (r *ServiceRequestControllerImp) New(writer http.ResponseWriter, request *http.Request) {

	var ServiceRequestSaveRequest transactionworkshoppayloads.ServiceRequestSaveRequest
	helper.ReadFromRequestBody(request, &ServiceRequestSaveRequest)

	success, baseErr := r.ServiceRequestService.New(ServiceRequestSaveRequest)
	if baseErr != nil {
		errorMap := map[string]int{
			"Invalid reference document type":   http.StatusBadRequest,
			"Invalid profit center combination": http.StatusBadRequest,
		}

		statusCode, exists := errorMap[baseErr.Message]
		if exists {
			payloads.NewHandleError(writer, baseErr.Message, statusCode)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, success, "Create Data Successfully", http.StatusCreated)
}

// Save saves service request
// @Summary Save service request
// @Description Save service request
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param service_request_system_number path int true "Service Request ID"
// @Param reqBody body transactionworkshoppayloads.ServiceRequestSaveDataRequest true "Service Request Data"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/{service_request_system_number} [put]
func (r *ServiceRequestControllerImp) Save(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	var ServiceRequestSaveRequest transactionworkshoppayloads.ServiceRequestSaveDataRequest
	helper.ReadFromRequestBody(request, &ServiceRequestSaveRequest)

	success, baseErr := r.ServiceRequestService.Save(ServiceRequestId, ServiceRequestSaveRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Save Data Successfully", http.StatusOK)
}

// Submit submits service request
// @Summary Submit service request
// @Description Submit service request
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param service_request_system_number path int true "Service Request ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/submit/{service_request_system_number} [post]
func (r *ServiceRequestControllerImp) Submit(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	success, newDocumentNumber, baseErr := r.ServiceRequestService.Submit(ServiceRequestId)
	if baseErr != nil {
		switch baseErr.Message {
		case "Service request has been submitted or the document number is already generated":
			responseDataError := struct {
				ServiceRequestSystemNumber int    `json:"service_request_system_number"`
				DocumentNumber             string `json:"service_request_document_number,omitempty"`
			}{
				ServiceRequestSystemNumber: ServiceRequestId,
			}
			if newDocumentNumber != "" {
				responseDataError.DocumentNumber = newDocumentNumber
			}
			response := payloads.Response{
				StatusCode: http.StatusConflict,
				Message:    baseErr.Message,
				Data:       responseDataError,
			}
			helper.WriteToResponseBody(writer, response)
		case "Data not found":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		default:
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		responseData := transactionworkshoppayloads.SubmitServiceRequestResponse{
			DocumentNumber:             newDocumentNumber,
			ServiceRequestSystemNumber: ServiceRequestId,
		}

		payloads.NewHandleSuccess(writer, responseData, "Service request submitted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to submit service request", http.StatusInternalServerError)
	}
}

// Void voids service request
// @Summary Void service request
// @Description Void service request
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param service_request_system_number path int true "Service Request ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/void/{service_request_system_number} [delete]
func (r *ServiceRequestControllerImp) Void(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.ServiceRequestService.Void(ServiceRequestId)
	if baseErr != nil {
		if baseErr.Message == "No service request data found" {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Service Request voided successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to void service request", http.StatusInternalServerError)
	}
}

// CloseOrder closes order
// @Summary Close order
// @Description Close order
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Param service_request_system_number path int true "Service Request ID"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/close/{service_request_system_number} [patch]
func (r *ServiceRequestControllerImp) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	ServiceRequestStrId := chi.URLParam(request, "service_request_system_number")
	ServiceRequestId, err := strconv.Atoi(ServiceRequestStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service request ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.ServiceRequestService.CloseOrder(ServiceRequestId)
	if baseErr != nil {
		payloads.NewHandleError(writer, baseErr.Message, baseErr.StatusCode)
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Service Request closed successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to close service request", http.StatusInternalServerError)
	}
}

// GetAllServiceDetail gets all service detail
// @Summary Get all service detail
// @Description Get all service detail
// @Tags Transaction : Workshop Service Request Detail
// @Accept json
// @Produce json
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/detail [get]
func (r *ServiceRequestControllerImp) GetAllServiceDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	queryParams := map[string]string{
		"service_request_system_number": queryValues.Get("service_request_system_number"),
		"service_request_detail_id":     queryValues.Get("service_request_detail_id"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.ServiceRequestService.GetAllServiceDetail(criteria, paginate)
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

// GetServiceDetailById gets service detail by id
// @Summary Get service detail by id
// @Description Get service detail by id
// @Tags Transaction : Workshop Service Request Detail
// @Accept json
// @Produce json
// @Param service_request_detail_id path string true "Service Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/detail/{service_request_detail_id} [get]
func (r *ServiceRequestControllerImp) GetServiceDetailById(writer http.ResponseWriter, request *http.Request) {

	ServiceDetailSystemStrId := chi.URLParam(request, "service_request_detail_id")
	ServiceDetailSystemId, err := strconv.Atoi(ServiceDetailSystemStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service detail system ID", http.StatusBadRequest)
		return
	}

	serviceDetail, baseErr := r.ServiceRequestService.GetServiceDetailById(ServiceDetailSystemId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "service detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, serviceDetail, "Get Data Successfully", http.StatusOK)
}

// AddServiceDetail adds service detail
// @Summary Add service detail
// @Description Add service detail
// @Tags Transaction : Workshop Service Request Detail
// @Accept json
// @Produce json
// @Param service_request_system_number path string true "Service Detail System ID"
// @Param reqBody body transactionworkshoppayloads.ServiceDetailSaveRequest true "Service Detail Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/detail [post]
func (r *ServiceRequestControllerImp) AddServiceDetail(writer http.ResponseWriter, request *http.Request) {

	serviceRequestSystemNumberStr := chi.URLParam(request, "service_request_system_number")
	serviceRequestSystemNumber, err := strconv.Atoi(serviceRequestSystemNumberStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service detail system ID", http.StatusBadRequest)
		return
	}

	var serviceDetailSaveRequest transactionworkshoppayloads.ServiceDetailSaveRequest
	helper.ReadFromRequestBody(request, &serviceDetailSaveRequest)

	success, baseErr := r.ServiceRequestService.AddServiceDetail(serviceRequestSystemNumber, serviceDetailSaveRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Create Data Successfully", http.StatusCreated)
}

// UpdateServiceDetail updates service detail
// @Summary Update service detail
// @Description Update service detail
// @Tags Transaction : Workshop Service Request Detail
// @Accept json
// @Produce json
// @Param service_request_system_number path string true "Service Detail System ID"
// @Param service_request_detail_id path string true "Service Detail ID"
// @Param reqBody body transactionworkshoppayloads.ServiceDetailUpdateRequest true "Service Detail Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/detail/{service_request_system_number}/{service_request_detail_id} [put]
func (r *ServiceRequestControllerImp) UpdateServiceDetail(writer http.ResponseWriter, request *http.Request) {
	serviceDetailSystemStrId := chi.URLParam(request, "service_request_system_number")
	serviceDetailSystemId, err := strconv.Atoi(serviceDetailSystemStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service detail system ID", http.StatusBadRequest)
		return
	}

	serviceDetailStrId := chi.URLParam(request, "service_request_detail_id")
	serviceDetailId, err := strconv.Atoi(serviceDetailStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service detail ID", http.StatusBadRequest)
		return
	}

	var serviceDetailSaveRequest transactionworkshoppayloads.ServiceDetailUpdateRequest
	helper.ReadFromRequestBody(request, &serviceDetailSaveRequest)

	entity, baseErr := r.ServiceRequestService.UpdateServiceDetail(serviceDetailSystemId, serviceDetailId, serviceDetailSaveRequest)
	if baseErr != nil {
		switch baseErr.StatusCode {
		case http.StatusNotFound:
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		case http.StatusBadRequest:
			payloads.NewHandleError(writer, baseErr.Message, http.StatusBadRequest)
		case http.StatusInternalServerError:
			exceptions.NewAppException(writer, request, baseErr)
		default:
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	// If successful, return the updated entity
	payloads.NewHandleSuccess(writer, entity, "Update Data Successfully", http.StatusOK)
}

// DeleteServiceDetail deletes service detail
// @Summary Delete service detail
// @Description Delete service detail
// @Tags Transaction : Workshop Service Request Detail
// @Accept json
// @Produce json
// @Param service_request_system_number path string true "Service Detail System ID"
// @Param service_request_detail_id path string true "Service Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/detail/{service_request_system_number}/{service_request_detail_id} [delete]
func (r *ServiceRequestControllerImp) DeleteServiceDetail(writer http.ResponseWriter, request *http.Request) {

	ServiceDetailSystemStrId := chi.URLParam(request, "service_request_system_number")
	ServiceDetailSystemId, err := strconv.Atoi(ServiceDetailSystemStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service detail system ID", http.StatusBadRequest)
		return
	}

	ServiceDetailStrId := chi.URLParam(request, "service_request_detail_id")
	ServiceDetailId, err := strconv.Atoi(ServiceDetailStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service detail ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.ServiceRequestService.DeleteServiceDetail(ServiceDetailSystemId, ServiceDetailId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "service detail not found", http.StatusNotFound)
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

// DeleteServiceDetailMultiId deletes multiple service detail
// @Summary Delete multiple service detail
// @Description Delete multiple service detail
// @Tags Transaction : Workshop Service Request Detail
// @Accept json
// @Produce json
// @Param service_request_system_number path string true "Service Detail System ID"
// @Param multi_id path string true "Service Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/detail/{service_request_system_number}/{multi_id} [delete]
func (r *ServiceRequestControllerImp) DeleteServiceDetailMultiId(writer http.ResponseWriter, request *http.Request) {

	ServiceDetailSystemStrId := chi.URLParam(request, "service_request_system_number")
	ServiceDetailSystemId, err := strconv.Atoi(ServiceDetailSystemStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid service detail system ID", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid service request detail multi ID", http.StatusBadRequest)
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

	success, baseErr := r.ServiceRequestService.DeleteServiceDetailMultiId(ServiceDetailSystemId, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Service detail not found", http.StatusNotFound)
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

// NewStatus get dropdown status
// @Summary Get dropdown status
// @Description Get dropdown status
// @Tags Transaction : Workshop Service Request
// @Accept json
// @Produce json
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/service-request/status [get]
func (r *ServiceRequestControllerImp) NewStatus(writer http.ResponseWriter, request *http.Request) {

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

	statuses, err := r.ServiceRequestService.NewStatus(filters)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if len(statuses) > 0 {
		payloads.NewHandleSuccess(writer, statuses, "List of service request statuses", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}
