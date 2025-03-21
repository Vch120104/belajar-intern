package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	utils "after-sales/api/utils"
	"after-sales/api/validation"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type WorkOrderControllerImpl struct {
	WorkOrderService transactionworkshopservice.WorkOrderService
}

type WorkOrderController interface {
	GetAllRequest(writer http.ResponseWriter, request *http.Request)
	GetRequestById(writer http.ResponseWriter, request *http.Request)
	UpdateRequest(writer http.ResponseWriter, request *http.Request)
	AddRequest(writer http.ResponseWriter, request *http.Request)
	AddRequestMultiId(writer http.ResponseWriter, request *http.Request)
	DeleteRequest(writer http.ResponseWriter, request *http.Request)
	DeleteRequestMultiId(writer http.ResponseWriter, request *http.Request)

	GetAllVehicleService(writer http.ResponseWriter, request *http.Request)
	GetVehicleServiceById(writer http.ResponseWriter, request *http.Request)
	UpdateVehicleService(writer http.ResponseWriter, request *http.Request)
	AddVehicleService(writer http.ResponseWriter, request *http.Request)
	DeleteVehicleService(writer http.ResponseWriter, request *http.Request)
	DeleteVehicleServiceMultiId(writer http.ResponseWriter, request *http.Request)

	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	New(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
	Void(writer http.ResponseWriter, request *http.Request)
	CloseOrder(writer http.ResponseWriter, request *http.Request)

	GenerateDocumentNumber(writer http.ResponseWriter, request *http.Request)
	CalculateWorkOrderTotal(writer http.ResponseWriter, request *http.Request)

	GetAllDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	GetDetailByIdWorkOrder(writer http.ResponseWriter, request *http.Request)
	AddDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	UpdateDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	DeleteDetailWorkOrder(writer http.ResponseWriter, request *http.Request)
	DeleteDetailWorkOrderMultiId(writer http.ResponseWriter, request *http.Request)

	GetAllBooking(writer http.ResponseWriter, request *http.Request)
	GetBookingById(writer http.ResponseWriter, request *http.Request)
	NewBooking(writer http.ResponseWriter, request *http.Request)
	SaveBooking(writer http.ResponseWriter, request *http.Request)

	GetAllAffiliated(writer http.ResponseWriter, request *http.Request)
	GetAffiliatedById(writer http.ResponseWriter, request *http.Request)
	NewAffiliated(writer http.ResponseWriter, request *http.Request)
	SaveAffiliated(writer http.ResponseWriter, request *http.Request)

	DeleteCampaign(writer http.ResponseWriter, request *http.Request)
	AddContractService(writer http.ResponseWriter, request *http.Request)
	AddGeneralRepairPackage(writer http.ResponseWriter, request *http.Request)
	AddFieldAction(writer http.ResponseWriter, request *http.Request)
	ChangeBillTo(writer http.ResponseWriter, request *http.Request)
	ChangePhoneNo(writer http.ResponseWriter, request *http.Request)
	ConfirmPrice(writer http.ResponseWriter, request *http.Request)
}

func NewWorkOrderController(WorkOrderService transactionworkshopservice.WorkOrderService) WorkOrderController {
	return &WorkOrderControllerImpl{
		WorkOrderService: WorkOrderService,
	}
}

// GetAllService gets all services of a work order
// @Summary Get All Services of Work Order
// @Description Retrieve all services of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/requestservice [get]
func (r *WorkOrderControllerImpl) GetAllRequest(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	excludeParams := map[string]bool{
		"page":    true,
		"limit":   true,
		"sort_of": true,
		"sort_by": true,
	}

	filterConditions := make([]utils.FilterCondition, 0)
	for key, values := range queryValues {
		if len(values) > 0 && !excludeParams[key] {
			filterConditions = append(filterConditions, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: values[0],
			})
		}
	}

	result, err := r.WorkOrderService.GetAllRequest(filterConditions, paginate)
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

// GetServiceById gets a service of a work order by ID
// @Summary Get Service of Work Order By ID
// @Description Retrieve a service of a work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
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
// @Security AuthorizationKeyAuth
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
	if validationErr := validation.ValidationForm(writer, request, &groupRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	update, err := r.WorkOrderService.UpdateRequest(int(workorderID), int(requestID), groupRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Request updated successfully", http.StatusOK)

}

// AddRequest adds a new request to a work order
// @Summary Add Request to Work Order
// @Description Add a new request to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
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
	if validationErr := validation.ValidationForm(writer, request, &groupRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.WorkOrderService.AddRequest(int(workorderID), groupRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success.WorkOrderServiceId > 0 {
		payloads.NewHandleSuccess(writer, success, "Request added successfully", http.StatusCreated)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddRequest Multi adds multiple request to a work order
// @Summary Add Multiple Request to Work Order
// @Description Add multiple request to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservicemulti [post]
func (r *WorkOrderControllerImpl) AddRequestMultiId(writer http.ResponseWriter, request *http.Request) {

	workorderID, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {

		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	var groupRequests []transactionworkshoppayloads.WorkOrderServiceRequest

	err = json.NewDecoder(request.Body).Decode(&groupRequests)
	if err != nil {
		payloads.NewHandleError(writer, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	entities, baseErr := r.WorkOrderService.AddRequestMultiId(workorderID, groupRequests)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, entities, "Requests added successfully", http.StatusCreated)
}

// DeleteRequest deletes a request from a work order
// @Summary Delete Request from Work Order
// @Description Delete a request from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_id path string true "Work Order Service ID"
// @Success 204 {object} payloads.Response
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
		payloads.NewHandleSuccess(writer, delete, "Request deleted successfully", http.StatusNoContent)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

}

// DeleteRequestMultiId deletes multiple request from a work order
// @Summary Delete Multiple Request from Work Order
// @Description Delete multiple request from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @param multi_id query string true "Multiple Request ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/requestservice/{multi_id} [delete]
func (r *WorkOrderControllerImpl) DeleteRequestMultiId(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderstrID := chi.URLParam(request, "work_order_system_number")
	workorderID, err := strconv.Atoi(workorderstrID)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
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

	success, baseErr := r.WorkOrderService.DeleteRequestMultiId(workorderID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "request detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Service Detail deleted successfully", http.StatusNoContent)
	} else {
		payloads.NewHandleError(writer, "Failed to delete service detail", http.StatusInternalServerError)
	}

}

// GetAllVehicleService gets all vehicle services of a work order
// @Summary Get All Vehicle Services of Work Order
// @Description Retrieve all vehicle services of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/vehicleservice [get]
func (r *WorkOrderControllerImpl) GetAllVehicleService(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	excludeParams := map[string]bool{
		"page":    true,
		"limit":   true,
		"sort_of": true,
		"sort_by": true,
	}

	filterConditions := make([]utils.FilterCondition, 0)
	for key, values := range queryValues {
		if len(values) > 0 && !excludeParams[key] {
			filterConditions = append(filterConditions, utils.FilterCondition{
				ColumnField: key,
				ColumnValue: values[0],
			})
		}
	}

	result, err := r.WorkOrderService.GetAllVehicleService(filterConditions, paginate)
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

// GetVehicleServiceById gets a vehicle service of a work order by ID
// @Summary Get Vehicle Service of Work Order By ID
// @Description Retrieve a vehicle service of a work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
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
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_vehicle_id path string true "Work Order Vehicle Service ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderServiceVehicleRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice/{work_order_service_vehicle_id} [put]
func (r *WorkOrderControllerImpl) UpdateVehicleService(writer http.ResponseWriter, request *http.Request) {
	// Update vehicle service of a work order
	workorderID, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	vehicleServiceID, _ := strconv.Atoi(chi.URLParam(request, "work_order_service_vehicle_id"))

	var vehicleRequest transactionworkshoppayloads.WorkOrderServiceVehicleRequest
	helper.ReadFromRequestBody(request, &vehicleRequest)
	if validationErr := validation.ValidationForm(writer, request, &vehicleRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	update, err := r.WorkOrderService.UpdateVehicleService(int(workorderID), int(vehicleServiceID), vehicleRequest)

	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Vehicle service updated successfully", http.StatusOK)
}

// AddVehicleService adds a new vehicle service to a work order
// @Summary Add Vehicle Service to Work Order
// @Description Add a new vehicle service to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
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
	if validationErr := validation.ValidationForm(writer, request, &vehicleRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.WorkOrderService.AddVehicleService(int(workorderID), vehicleRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success.WorkOrderServiceVehicleId > 0 {
		payloads.NewHandleSuccess(writer, success, "Vehicle service added successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to add vehicle service", http.StatusInternalServerError)
	}
}

// DeleteVehicleService deletes a vehicle service from a work order
// @Summary Delete Vehicle Service from Work Order
// @Description Delete a vehicle service from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_service_vehicle_id path string true "Work Order Vehicle Service ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice/{work_order_service_vehicle_id} [delete]
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
		payloads.NewHandleSuccess(writer, delete, "Vehicle service deleted successfully", http.StatusNoContent)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

}

// DeleteVehicleServiceMultiId deletes multiple a vehicle service from a work order
// @Summary Delete multiple vehicle service
// @Description  Delete multiple vehicle service
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Service Detail System ID"
// @Param multi_id path string true "Service Detail ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/vehicleservice/{multi_id} [delete]
func (r *WorkOrderControllerImpl) DeleteVehicleServiceMultiId(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderstrID := chi.URLParam(request, "work_order_system_number")
	workorderID, err := strconv.Atoi(workorderstrID)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
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

	success, baseErr := r.WorkOrderService.DeleteVehicleServiceMultiId(workorderID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "request detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Vehicle service deleted successfully", http.StatusNoContent)
	} else {
		payloads.NewHandleError(writer, "Failed to delete Vehicle detail", http.StatusInternalServerError)
	}

}

// GetAll gets all work orders
// @Summary Get All Work Orders
// @Description Retrieve all work orders with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
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
		"trx_work_order.work_order_document_number": queryValues.Get("work_order_document_number"),
		"trx_work_order.work_order_system_number":   queryValues.Get("work_order_system_number"),
		"trx_work_order.work_order_date_from":       queryValues.Get("work_order_date_from"),
		"trx_work_order.work_order_date_to":         queryValues.Get("work_order_date_to"),
		"trx_work_order.work_order_type_id":         queryValues.Get("work_order_type_id"),
		"trx_work_order.brand_id":                   queryValues.Get("brand_id"),
		"trx_work_order.model_id":                   queryValues.Get("model_id"),
		"trx_work_order.vehicle_chassis_number":     queryValues.Get("vehicle_chassis_number"),
		"trx_work_order.vehicle_tnkb":               queryValues.Get("vehicle_tnkb"),
		"trx_work_order.work_order_status_id":       queryValues.Get("work_order_status_id"),
		"trx_work_order.variant_id":                 queryValues.Get("variant_id"),
		"trx_work_order.foreman_id":                 queryValues.Get("foreman_id"),
		"trx_work_order.service_advisor_id":         queryValues.Get("service_advisor_id"),
		"trx_work_order.company_id":                 queryValues.Get("company_id"),
		"trx_work_order.name_customer":              queryValues.Get("name_customer"),
		"trx_work_order.dealer_representative_id":   queryValues.Get("dealer_representative_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.WorkOrderService.GetAll(criteria, paginate)
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

// New creates a new work order
// @Summary Create New Work Order
// @Description Create a new work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param reqBody body transactionworkshoppayloads.WorkOrderNormalRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal [post]
func (r *WorkOrderControllerImpl) New(writer http.ResponseWriter, request *http.Request) {

	var workOrderRequest transactionworkshoppayloads.WorkOrderNormalRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)
	if validationErr := validation.ValidationForm(writer, request, &workOrderRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.WorkOrderService.New(workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if success.WorkOrderSystemNumber > 0 {
		payloads.NewHandleSuccess(writer, success, "Work order created successfully", http.StatusCreated)
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
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number} [get]
func (r *WorkOrderControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	workOrder, baseErr := r.WorkOrderService.GetById(workOrderId, paginate)
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
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
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

	var workOrderRequest transactionworkshoppayloads.WorkOrderNormalSaveRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)
	if validationErr := validation.ValidationForm(writer, request, &workOrderRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, baseErr := r.WorkOrderService.Save(workOrderRequest, workOrderId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Work order saved successfully", http.StatusOK)

}

// Submit submits a new work order
// @Summary Submit Work Order
// @Description Submit a new work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path int true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/submit/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {
	// Retrieve work order ID from URL parameters
	workOrderId := chi.URLParam(request, "work_order_system_number")
	workOrderIdInt, err := strconv.Atoi(workOrderId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, newDocumentNumber, baseErr := r.WorkOrderService.Submit(workOrderIdInt)
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
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path int true "Work Order ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/void/{work_order_system_number} [delete]
func (r *WorkOrderControllerImpl) Void(writer http.ResponseWriter, request *http.Request) {
	// Void work order
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.WorkOrderService.Void(workOrderId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Work order voided successfully", http.StatusNoContent)
	} else {
		payloads.NewHandleError(writer, "Failed to void work order", http.StatusInternalServerError)
	}
}

// CloseOrder closes a work order
// @Summary Close Work Order
// @Description Close an existing work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path int true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/close/{work_order_system_number} [patch]
func (r *WorkOrderControllerImpl) CloseOrder(writer http.ResponseWriter, request *http.Request) {
	// Close work order
	workOrderId := chi.URLParam(request, "work_order_system_number")
	workOrderIdInt, err := strconv.Atoi(workOrderId)

	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, baseErr := r.WorkOrderService.CloseOrder(workOrderIdInt)
	if baseErr != nil {
		switch baseErr.Message {
		case "Work order not found":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusNotFound)
		case "Work order cannot be closed because status is draft":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "There is still DP payment that has not been settled":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Detail Work Order without Invoice No must be deleted":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Warranty Item (PTP) must be supplied":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Warranty Item (PTM)/Operation must be Invoiced":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		case "Service Mileage must be larger than Last Mileage.":
			payloads.NewHandleError(writer, baseErr.Message, http.StatusConflict)
		default:
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Work order closed successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to close work order", http.StatusInternalServerError)
	}

}

// GetWorkOrderDetail gets the detail of a work order
// @Summary Get Work Order Detail
// @Description Retrieve the detail of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal Detail
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/detail [get]
func (r *WorkOrderControllerImpl) GetAllDetailWorkOrder(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order_detail.work_order_system_number": queryValues.Get("work_order_system_number"),
		"trx_work_order_detail.work_order_detail_id":     queryValues.Get("work_order_detail_id"),
		"trx_work_order_detail.transaction_type_id":      queryValues.Get("transaction_type_id"),
		"trx_work_order_detail.supply_quantity":          queryValues.Get("supply_quantity"),
		"trx_work_order_detail.price_list_id":            queryValues.Get("price_list_id"),
		"trx_work_order_detail.line_type_id":             queryValues.Get("line_type_id"),
		"trx_work_order_detail.job_type_id":              queryValues.Get("job_type_id"),
		"trx_work_order_detail.frt_quantity":             queryValues.Get("frt_quantity"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.WorkOrderService.GetAllDetailWorkOrder(criteria, paginate)
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

// GetDetailWorkOrderById gets the detail of a work order by ID
// @Summary Get Work Order Detail By ID
// @Description Retrieve the detail of a work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_detail_id path string true "Work Order Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail/{work_order_detail_id} [get]
func (r *WorkOrderControllerImpl) GetDetailByIdWorkOrder(writer http.ResponseWriter, request *http.Request) {
	// Get the detail of a work order by ID
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}
	detailId, err := strconv.Atoi(chi.URLParam(request, "work_order_detail_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order detail ID", http.StatusBadRequest)
		return
	}

	detail, baseErr := r.WorkOrderService.GetDetailByIdWorkOrder(workOrderId, detailId)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Work order detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, detail, "Get Data Successfully", http.StatusOK)

}

// UpdateDetailWorkOrder updates the detail of a work order
// @Summary Update Work Order Detail
// @Description Update the detail of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal Detail
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
	if validationErr := validation.ValidationForm(writer, request, &detailRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	update, err := r.WorkOrderService.UpdateDetailWorkOrder(int(workOrderId), int(detailId), detailRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if update.WorkOrderSystemNumber > 0 {
		payloads.NewHandleSuccess(writer, update, "Detail updated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// AddDetailWorkOrder adds a new detail to a work order
// @Summary Add Work Order Detail
// @Description Add a new detail to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderDetailRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail [post]
func (r *WorkOrderControllerImpl) AddDetailWorkOrder(writer http.ResponseWriter, request *http.Request) {
	workOrderStrId := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var detailRequest transactionworkshoppayloads.WorkOrderDetailRequest
	helper.ReadFromRequestBody(request, &detailRequest)
	if validationErr := validation.ValidationForm(writer, request, &detailRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, serviceErr := r.WorkOrderService.AddDetailWorkOrder(workOrderId, detailRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Detail added successfully", http.StatusCreated)

}

// DeleteDetailWorkOrder deletes a detail from a work order
// @Summary Delete Work Order Detail
// @Description Delete a detail from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal Detail
// @Param work_order_system_number path string true "Work Order ID"
// @Param work_order_detail_id path string true "Work Order Detail ID"
// @Success 204 {object} payloads.Response
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
		payloads.NewHandleSuccess(writer, delete, "Detail deleted successfully", http.StatusNoContent)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}

}

// DeleteDetailWorkOrder deletes a multiple detail from a work order
// @Summary Delete multiple Work Order Detail
// @Description  Delete multiple Work Order Detail
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal Detail
// @Param work_order_system_number path string true "Work Order System Number ID"
// @Param multi_id path string true "Work Order Detail ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/{work_order_system_number}/detail/{multi_id} [delete]
func (r *WorkOrderControllerImpl) DeleteDetailWorkOrderMultiId(writer http.ResponseWriter, request *http.Request) {
	// Delete request from work order
	workorderstrID := chi.URLParam(request, "work_order_system_number")
	workorderID, err := strconv.Atoi(workorderstrID)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order system number", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
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

	success, baseErr := r.WorkOrderService.DeleteDetailWorkOrderMultiId(workorderID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "request detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Vehicle service deleted successfully", http.StatusNoContent)
	} else {
		payloads.NewHandleError(writer, "Failed to delete Vehicle detail", http.StatusInternalServerError)
	}

}

// GetAllWorkOrderBooking gets all work order bookings
// @Summary Get All Work Order Booking
// @Description Retrieve all work order bookings
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Booking
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking [get]
func (r *WorkOrderControllerImpl) GetAllBooking(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_work_order.work_order_system_number":      chi.URLParam(request, "work_order_system_number"),
		"trx_work_order.booking_system_number":         chi.URLParam(request, "booking_system_number"),
		"trx_work_order.estimation_system_number":      chi.URLParam(request, "estimation_system_number"),
		"trx_work_order.service_request_system_number": chi.URLParam(request, "service_request_system_number"),
		"trx_work_order.brand_id":                      chi.URLParam(request, "brand_id"),
		"trx_work_order.model_id":                      chi.URLParam(request, "model_id"),
		"trx_work_order.vehicle_id":                    chi.URLParam(request, "vehicle_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.WorkOrderService.GetAllBooking(criteria, paginate)
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

// GetWorkOrderBookingById gets a work
// @Summary Get Work Order Booking By ID
// @Description Retrieve a work order booking by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Booking
// @Param work_order_system_number path string true "Work Order ID"
// @Param booking_system_number path string true "Work Order Booking ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/{work_order_system_number}/{booking_system_number} [get]
func (r *WorkOrderControllerImpl) GetBookingById(writer http.ResponseWriter, request *http.Request) {
	// Get a work order booking by ID
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	workOrderBookId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	workOrder, baseErr := r.WorkOrderService.GetBookingById(workOrderId, workOrderBookId, paginate)
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

// UpdateWorkOrderBooking updates a work order booking
// @Summary Update Work Order Booking
// @Description Update a work order booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Booking
// @Param work_order_system_number path string true "Work Order ID"
// @Param booking_system_number path string true "Work Order Booking ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderBookingRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking/{work_order_system_number}/{booking_system_number} [put]
func (r *WorkOrderControllerImpl) SaveBooking(writer http.ResponseWriter, request *http.Request) {
	// Update a work order booking
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	workOrderBookId, _ := strconv.Atoi(chi.URLParam(request, "booking_system_number"))

	var workOrderRequest transactionworkshoppayloads.WorkOrderBookingRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, err := r.WorkOrderService.SaveBooking(workOrderId, workOrderBookId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	if result {
		payloads.NewHandleSuccess(writer, result, "Work order saved successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to save work order", http.StatusInternalServerError)
	}

}

// AddWorkOrderBooking adds a new work order booking
// @Summary Add Work Order Booking
// @Description Add a new work order booking
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Booking
// @Param reqBody body transactionworkshoppayloads.WorkOrderBookingRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/booking [post]
func (r *WorkOrderControllerImpl) NewBooking(writer http.ResponseWriter, request *http.Request) {

	var workOrderRequest transactionworkshoppayloads.WorkOrderBookingRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, err := r.WorkOrderService.NewBooking(workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order saved successfully", http.StatusCreated)

}

// GetAllAffiliated gets all affiliated work orders
// @Summary Get All Affiliated Work Orders
// @Description Retrieve all affiliated work orders
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Affiliated
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

	result, err := r.WorkOrderService.GetAllAffiliated(criteria, paginate)
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

// GetAffiliatedById gets an affiliated work order by ID
// @Summary Get Affiliated Work Order By ID
// @Description Retrieve an affiliated work order by ID
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Affiliated
// @Param work_order_system_number path string true "Work Order ID"
// @Param service_request_system_number path string true "Affiliated Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated/{work_order_system_number}/{service_request_system_number} [get]
func (r *WorkOrderControllerImpl) GetAffiliatedById(writer http.ResponseWriter, request *http.Request) {
	// Get affiliated work order by ID
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	affiliatedWorkOrderId, _ := strconv.Atoi(chi.URLParam(request, "service_request_system_number"))

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	workOrder, baseErr := r.WorkOrderService.GetAffiliatedById(workOrderId, affiliatedWorkOrderId, paginate)
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

// NewAffiliated creates a new affiliated work order
// @Summary Create New Affiliated Work Order
// @Description Create a new affiliated work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Affiliated
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/affiliated [post]
func (r *WorkOrderControllerImpl) NewAffiliated(writer http.ResponseWriter, request *http.Request) {
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	var workOrderRequest transactionworkshoppayloads.WorkOrderAffiliatedRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, err := r.WorkOrderService.NewAffiliated(workOrderId, workOrderRequest)
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
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Affiliated
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

	result, err := r.WorkOrderService.SaveAffiliated(workOrderId, affiliatedWorkOrderId, workOrderRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order updated successfully", http.StatusOK)
}

// GenerateWorkOrderDocumentNumber generates a new work order document number
// @Summary Generate Work Order Document Number
// @Description Generate a new work order document number
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path string true "Work Order ID"
// Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/document-number/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) GenerateDocumentNumber(writer http.ResponseWriter, request *http.Request) {
	// Generate a new work order document number
	workOrderId, _ := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))

	result, err := r.WorkOrderService.GenerateDocumentNumber(workOrderId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Document number generated successfully", http.StatusOK)

}

// ChangeBillTo changes the bill to of a work order
// @Summary Change Bill To
// @Description Change the bill to of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.ChangeBillToRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/change-bill-to/{work_order_system_number} [put]
func (r *WorkOrderControllerImpl) ChangeBillTo(writer http.ResponseWriter, request *http.Request) {
	// Change the bill to of a work order
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.ChangeBillToRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	result, baseErr := r.WorkOrderService.ChangeBillTo(workOrderId, workOrderRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Bill to changed successfully", http.StatusOK)
}

// ChangePhoneNo changes the phone number of a work order
// @Summary Change Phone Number
// @Description Change the phone number of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.ChangePhoneNoRequest true "Work Order Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/change-phone-no/{work_order_system_number} [put]
func (r *WorkOrderControllerImpl) ChangePhoneNo(writer http.ResponseWriter, request *http.Request) {
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.ChangePhoneNoRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	updatedPayload, baseErr := r.WorkOrderService.ChangePhoneNo(workOrderId, workOrderRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, updatedPayload, "Phone number changed successfully", http.StatusOK)
}

// ConfirmPriceList confirms the price list of a work order
// @Summary Confirm Price List
// @Description Confirm the price list of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth Normal
// @Param work_order_system_number path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/confirm-price/{work_order_system_number}/{multi_id} [put]
func (r *WorkOrderControllerImpl) ConfirmPrice(writer http.ResponseWriter, request *http.Request) {
	workOrderIdStr := chi.URLParam(request, "work_order_system_number")
	workOrderId, err := strconv.Atoi(workOrderIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
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

	var confirmPriceRequest transactionworkshoppayloads.WorkOrderConfirmPriceRequest
	helper.ReadFromRequestBody(request, &confirmPriceRequest)

	result, baseErr := r.WorkOrderService.ConfirmPrice(workOrderId, intIds, confirmPriceRequest)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Confirm Price list confirmed successfully", http.StatusOK)
}

// DeleteCampaign deletes a campaign from a work order
// @Summary Delete Work Order Campaign
// @Description Delete a campaign from a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order Campaign ID"
// @Success 204 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/delete-campaign/{work_order_system_number} [delete]
func (r *WorkOrderControllerImpl) DeleteCampaign(writer http.ResponseWriter, request *http.Request) {

	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	success, serviceErr := r.WorkOrderService.DeleteCampaign(workOrderId)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Campaign deleted successfully", http.StatusNoContent)

}

// AddContractService adds a new contract service to a work order
// @Summary Add Work Order Contract Service
// @Description Add a new contract service to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderContractServiceRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/contract-service/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) AddContractService(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.WorkOrderContractServiceRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, serviceErr := r.WorkOrderService.AddContractService(workOrderId, workOrderRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Contract service added successfully", http.StatusCreated)
}

// AddGeneralRepairPackage adds a new general repair package to a work order
// @Summary Add Work Order General Repair Package
// @Description Add a new general repair package to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/general-repair-package/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) AddGeneralRepairPackage(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, serviceErr := r.WorkOrderService.AddGeneralRepairPackage(workOrderId, workOrderRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "General repair package added successfully", http.StatusCreated)
}

// AddFieldAction adds a new field action to a work order
// @Summary Add Work Order Field Action
// @Description Add a new field action to a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Param reqBody body transactionworkshoppayloads.WorkOrderFieldActionRequest true "Work Order Data"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/field-action/{work_order_system_number} [post]
func (r *WorkOrderControllerImpl) AddFieldAction(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	var workOrderRequest transactionworkshoppayloads.WorkOrderFieldActionRequest
	helper.ReadFromRequestBody(request, &workOrderRequest)

	success, serviceErr := r.WorkOrderService.AddFieldAction(workOrderId, workOrderRequest)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, success, "Field action added successfully", http.StatusCreated)
}

// CalculateWorkOrderTotal calculates the total of a work order
// @Summary Calculate Work Order Total
// @Description Calculate the total of a work order
// @Accept json
// @Produce json
// @Tags Transaction : Workshop Work Order
// @Security AuthorizationKeyAuth
// @Param work_order_system_number path string true "Work Order ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order/normal/calculate-total/{work_order_system_number} [put]
func (r *WorkOrderControllerImpl) CalculateWorkOrderTotal(writer http.ResponseWriter, request *http.Request) {
	workOrderId, err := strconv.Atoi(chi.URLParam(request, "work_order_system_number"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid work order ID", http.StatusBadRequest)
		return
	}

	result, serviceErr := r.WorkOrderService.CalculateWorkOrderTotal(workOrderId)
	if serviceErr != nil {
		exceptions.NewAppException(writer, request, serviceErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Work order total calculated successfully", http.StatusOK)
}
