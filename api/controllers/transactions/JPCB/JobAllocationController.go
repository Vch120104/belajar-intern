package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type JobAllocationController interface {
	GetAllJobAllocation(writer http.ResponseWriter, request *http.Request)
	GetJobAllocationById(writer http.ResponseWriter, request *http.Request)
	UpdateJobAllocation(writer http.ResponseWriter, request *http.Request)
	DeleteJobAllocation(writer http.ResponseWriter, request *http.Request)
}

type JobAllocationControllerImpl struct {
	JobAllocationService transactionjpcbservice.JobAllocationService
}

func NewJobAllocationController(jobAllocationService transactionjpcbservice.JobAllocationService) JobAllocationController {
	return &JobAllocationControllerImpl{
		JobAllocationService: jobAllocationService,
	}
}

// @Summary Get All Job Allocation
// @Description Get All Job Allocation
// @Tags Transaction : JPCB Job Allocation
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param technician_id query string false "Technician ID"
// @Param service_date query string false "Service Date"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/job-allocation [get]
func (r *JobAllocationControllerImpl) GetAllJobAllocation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":    queryValues.Get("company_id"),
		"technician_id": queryValues.Get("technician_id"),
		"service_data":  queryValues.Get("service_date"),
	}

	paginate := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "page"),
	}
	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.JobAllocationService.GetAllJobAllocation(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Job Allocation By ID
// @Description Get Job Allocation By ID
// @Tags Transaction : JPCB Job Allocation
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param technician_allocation_system_number path string true "Technician Allocation System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/job-allocation/{technician_allocation_system_number} [get]
func (r *JobAllocationControllerImpl) GetJobAllocationById(writer http.ResponseWriter, request *http.Request) {
	technicianAllocationSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "technician_allocation_system_number"))

	result, err := r.JobAllocationService.GetJobAllocationById(technicianAllocationSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// @Summary Update Job Allocation
// @Description Update Job Allocation
// @Tags Transaction : JPCB Job Allocation
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param technician_allocation_system_number path string true "Technician Allocation System Number"
// @Param body body transactionjpcbpayloads.JobAllocationUpdateRequest true "Job Allocation Update Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/job-allocation/{technician_allocation_system_number} [put]
func (r *JobAllocationControllerImpl) UpdateJobAllocation(writer http.ResponseWriter, request *http.Request) {
	technicianAllocationSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "technician_allocation_system_number"))

	formRequest := transactionjpcbpayloads.JobAllocationUpdateRequest{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	update, err := r.JobAllocationService.UpdateJobAllocation(technicianAllocationSystemNumber, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Update Data Successfully", http.StatusOK)
}

// @Summary Delete Job Allocation
// @Description Delete Job Allocation
// @Tags Transaction : JPCB Job Allocation
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param technician_allocation_system_number path string true "Technician Allocation System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/job-allocation/{technician_allocation_system_number} [delete]
func (r *JobAllocationControllerImpl) DeleteJobAllocation(writer http.ResponseWriter, request *http.Request) {
	technicianAllocationSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "technician_allocation_system_number"))

	delete, err := r.JobAllocationService.DeleteJobAllocation(technicianAllocationSystemNumber)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, delete, "Delete Data Successfully", http.StatusOK)
}
