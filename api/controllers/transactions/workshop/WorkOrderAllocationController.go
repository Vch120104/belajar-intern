package transactionworkshopcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"log"
	"strconv"
	"time"

	"net/http"

	"github.com/go-chi/chi/v5"
)

type WorkOrderAllocationControllerImp struct {
	WorkOrderAllocationService transactionworkshopservice.WorkOrderAllocationService
}

type WorkOrderAllocationController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetWorkOrderAllocationHeaderData(writer http.ResponseWriter, request *http.Request)
	GetAllocate(writer http.ResponseWriter, request *http.Request)
	GetAllocateDetail(writer http.ResponseWriter, request *http.Request)
	SaveAllocateDetail(writer http.ResponseWriter, request *http.Request)

	GetAssignTechnician(writer http.ResponseWriter, request *http.Request)
	GetAssignTechnicianById(writer http.ResponseWriter, request *http.Request)
	NewAssignTechnician(writer http.ResponseWriter, request *http.Request)
	SaveAssignTechnician(writer http.ResponseWriter, request *http.Request)
}

func NewWorkOrderAllocationController(service transactionworkshopservice.WorkOrderAllocationService) WorkOrderAllocationController {
	return &WorkOrderAllocationControllerImp{
		WorkOrderAllocationService: service,
	}
}

// GetAll gets all datagrid workorder allocation
// @Summary Get all datagrid workorder allocation
// @Description Get all datagrid workorder allocation
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/{service_date}/{foreman_id}/{company_id} [get]
func (r *WorkOrderAllocationControllerImp) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"service_date": queryValues.Get("service_date"),
		"foreman_id":   queryValues.Get("foreman_id"),
		"company_id":   queryValues.Get("company_id"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	serviceDateStr := chi.URLParam(request, "service_date")
	if serviceDateStr == "" {
		payloads.NewHandleError(writer, "Service date is required", http.StatusBadRequest)
		return
	}

	// Attempt to parse serviceDateStr to time.Time
	serviceRequestDate, err := time.Parse("2006-01-02", serviceDateStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid date format", http.StatusBadRequest)
		return
	}

	technicianStrId := chi.URLParam(request, "foreman_id")
	technicianId, err := strconv.Atoi(technicianStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Foreman ID", http.StatusBadRequest)
		return
	}

	companyStrId := chi.URLParam(request, "company_id")
	companyId, err := strconv.Atoi(companyStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company Id", http.StatusBadRequest)
		return
	}

	// Call service to fetch data
	paginatedData, apiErr := r.WorkOrderAllocationService.GetAll(
		companyId,
		technicianId,
		serviceRequestDate,
		criteria,
	)
	if apiErr != nil {
		exceptions.NewNotFoundException(writer, request, apiErr)
		return
	}

	// Handle the response
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccess(writer, paginatedData, "Get Data Successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetAllocate gets all allocated work orders
// @Summary Get all allocated work orders
// @Description Get all allocated work orders
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param service_date query string true "Service Request Date"
// @Param brand_id query int true "Brand ID"
// @Param work_order_system_number query int true "Work Order System Number"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/allocate/{brand_id}/{work_order_system_number} [get]
func (r *WorkOrderAllocationControllerImp) GetAllocate(writer http.ResponseWriter, request *http.Request) {
	brandStrId := chi.URLParam(request, "brand_id")
	brandId, err := strconv.Atoi(brandStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Brand ID", http.StatusBadRequest)
		return
	}

	workorderStrId := chi.URLParam(request, "work_order_system_number")
	workorderId, err := strconv.Atoi(workorderStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Work Order System Number", http.StatusBadRequest)
		return
	}

	woAssign, baseErr := r.WorkOrderAllocationService.GetAllocate(brandId, workorderId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, woAssign, "Data retrieved successfully", http.StatusOK)
}

// GetAllocateDetail gets all allocated work orders detail
// @Summary Get all allocated work orders detail
// @Description Get all allocated work orders detail
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param service_date query string true "Service Request Date"
// @Param foreman_id query int true "Foreman ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/allocate-detail/{service_date}/{foreman_id} [get]
func (r *WorkOrderAllocationControllerImp) GetAllocateDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"service_date": queryValues.Get("service_date"),
		"foreman_id":   queryValues.Get("foreman_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderAllocationService.GetAllocateDetail(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, paginatedData, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetAssignTechnician gets all assigned technicians
// @Summary Get all assigned technicians
// @Description Get all assigned technicians
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param service_date query string true "Service Request Date"
// @Param foreman_id query int true "Foreman ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/assign-technician [get]
func (r *WorkOrderAllocationControllerImp) GetAssignTechnician(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"service_date": queryValues.Get("service_date"),
		"foreman_id":   queryValues.Get("foreman_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.WorkOrderAllocationService.GetAssignTechnician(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, paginatedData, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// GetAssignTechnicianById gets all assigned technicians by ID
// @Summary Get all assigned technicians by ID
// @Description Get all assigned technicians by ID
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param service_date query string true "Service Request Date"
// @Param foreman_id query int true "Foreman ID"
// @Param assign_technician_id query int true "Assign Technician ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/assign-technician/{service_date}/{foreman_id}/{assign_technician_id} [get]
func (r *WorkOrderAllocationControllerImp) GetAssignTechnicianById(writer http.ResponseWriter, request *http.Request) {
	serviceDateStr := chi.URLParam(request, "service_date")
	if serviceDateStr == "" {
		payloads.NewHandleError(writer, "Service date is required", http.StatusBadRequest)
		return
	}

	// Attempt to parse serviceDateStr to time.Time
	serviceRequestDate, err := time.Parse("2006-01-02", serviceDateStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid date format", http.StatusBadRequest)
		return
	}

	technicianStrId := chi.URLParam(request, "foreman_id")
	technicianId, err := strconv.Atoi(technicianStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Foreman ID", http.StatusBadRequest)
		return
	}

	AssignStrId := chi.URLParam(request, "assign_technician_id")
	AssignId, err := strconv.Atoi(AssignStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Assign Technician Id", http.StatusBadRequest)
		return
	}

	WoAssign, baseErr := r.WorkOrderAllocationService.GetAssignTechnicianById(serviceRequestDate, technicianId, AssignId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, WoAssign, "Get Data Successfully", http.StatusOK)

}

// NewAssignTechnician assigns a new technician
// @Summary Assign a new technician
// @Description Assign a new technician
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param service_date query string true "Service Request Date"
// @Param foreman_id query int true "Foreman ID"
// @Param request body transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest true "Request body"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/assign-technician/{service_date}/{foreman_id} [post]
func (r *WorkOrderAllocationControllerImp) NewAssignTechnician(writer http.ResponseWriter, request *http.Request) {
	serviceDateStr := chi.URLParam(request, "service_date")
	if serviceDateStr == "" {
		payloads.NewHandleError(writer, "Service date is required", http.StatusBadRequest)
		return
	}

	// Attempt to parse serviceDateStr to time.Time
	serviceRequestDate, err := time.Parse(time.RFC3339, serviceDateStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid date format", http.StatusBadRequest)
		return
	}

	technicianStrId := chi.URLParam(request, "foreman_id")
	technicianId, err := strconv.Atoi(technicianStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Foreman ID", http.StatusBadRequest)
		return
	}

	var req transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest
	helper.ReadFromRequestBody(request, &req)

	// Pass the parsed date to your service
	entity, baseErr := r.WorkOrderAllocationService.NewAssignTechnician(serviceRequestDate, technicianId, req)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, entity, "Assign Technician Successfully", http.StatusOK)
}

// SaveAssignTechnician saves assigned technician
// @Summary Save assigned technician
// @Description Save assigned technician
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param service_date query string true "Service Request Date"
// @Param technician_id query int true "Technician ID"
// @Param assign_technician_id query int true "Assign Technician ID"
// @Param request body transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest true "Request body"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/assign-technician/{service_date}/{technician_id}/{assign_technician_id} [put]
func (r *WorkOrderAllocationControllerImp) SaveAssignTechnician(writer http.ResponseWriter, request *http.Request) {
	serviceDateStr := chi.URLParam(request, "service_date")
	if serviceDateStr == "" {
		payloads.NewHandleError(writer, "Service date is required", http.StatusBadRequest)
		return
	}

	// Attempt to parse serviceDateStr to time.Time
	serviceRequestDate, err := time.Parse(time.RFC3339, serviceDateStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid date format", http.StatusBadRequest)
		return
	}

	technicianStrId := chi.URLParam(request, "foreman_id")
	technicianId, err := strconv.Atoi(technicianStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Foreman ID", http.StatusBadRequest)
		return
	}

	AssignStrId := chi.URLParam(request, "assign_technician_id")
	AssignId, err := strconv.Atoi(AssignStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Assign Technician Id", http.StatusBadRequest)
		return
	}

	var req transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest
	helper.ReadFromRequestBody(request, &req)

	entity, baseErr := r.WorkOrderAllocationService.SaveAssignTechnician(serviceRequestDate, technicianId, AssignId, req)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, entity, "Save Technician Successfully", http.StatusOK)
}

// GetWorkOrderAllocationHeaderData gets all work order allocation header data
// @Summary Get all work order allocation header data
// @Description Get all work order allocation header data
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/header-data [get]
func (r *WorkOrderAllocationControllerImp) GetWorkOrderAllocationHeaderData(writer http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	companyCode := queryParams.Get("company_id")
	foremanID := queryParams.Get("foreman_id")
	serviceDate := queryParams.Get("service_date")
	brandID := queryParams.Get("brand_id")

	log.Printf("Received Query Params - company_id: %s, foreman_id: %s, service_date: %s, brand_id: %s", companyCode, foremanID, serviceDate, brandID)

	foremanId, err := strconv.Atoi(foremanID)
	if err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			Message:    "Invalid foreman ID",
			Err:        err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	techallocStartDate, err := time.Parse("2006-01-02", serviceDate)
	if err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			Message:    "Invalid start date",
			Err:        err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	vehicleBrandId, err := strconv.Atoi(brandID)
	if err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			Message:    "Invalid vehicle brand ID",
			Err:        err,
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	data, baseErr := r.WorkOrderAllocationService.GetWorkOrderAllocationHeaderData(companyCode, foremanId, techallocStartDate, vehicleBrandId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, data, "Get Data Successfully", http.StatusOK)
}

// SaveAllocateDetail saves allocated work order detail
// @Summary Save allocated work order detail
// @Description Save allocated work order detail
// @Tags Transaction : Workshop Work Order Allocation
// @Accept json
// @Produce json
// @Param service_date query string true "Service Request Date"
// @Param technician_id query int true "Foreman ID"
// @Param request body transactionworkshoppayloads.WorkOrderAllocationDetailRequest true "Request body"
// @Success 200 {object}  payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/work-order-allocation/allocate-detail/{service_date}/{technician_id} [post]
func (r *WorkOrderAllocationControllerImp) SaveAllocateDetail(writer http.ResponseWriter, request *http.Request) {
	serviceDateStr := chi.URLParam(request, "service_date")
	if serviceDateStr == "" {
		payloads.NewHandleError(writer, "Service date is required", http.StatusBadRequest)
		return
	}

	// Attempt to parse serviceDateStr to time.Time
	serviceRequestDate, err := time.Parse("2006-01-02", serviceDateStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid date format", http.StatusBadRequest)
		return
	}

	technicianStrId := chi.URLParam(request, "technician_id")
	technicianId, err := strconv.Atoi(technicianStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Technician ID", http.StatusBadRequest)
		return
	}

	foremanStrId := chi.URLParam(request, "foreman_id")
	foremanId, err := strconv.Atoi(foremanStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Foreman ID", http.StatusBadRequest)
		return
	}

	companyStrId := chi.URLParam(request, "company_id")
	companyId, err := strconv.Atoi(companyStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	var req transactionworkshoppayloads.WorkOrderAllocationDetailRequest
	helper.ReadFromRequestBody(request, &req)

	entity, baseErr := r.WorkOrderAllocationService.SaveAllocateDetail(serviceRequestDate, technicianId, req, foremanId, companyId)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, entity, "Save Allocate Detail Successfully", http.StatusOK)
}
