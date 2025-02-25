package transactionjpcbcontroller

import (
	"after-sales/api/exceptions"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type OutstandingJobAllocationController interface {
	GetAllOutstandingJobAllocation(writer http.ResponseWriter, request *http.Request)
	GetByTypeIdOutstandingJobAllocation(writer http.ResponseWriter, request *http.Request)
	SaveOutstandingJobAllocation(writer http.ResponseWriter, request *http.Request)
}

type OutstandingJobAllocationControllerImpl struct {
	OutstandingJobAllocationService transactionjpcbservice.OutstandingJobAllocationService
}

func NewOutstandingJobAllocationController(outstandingJobAllocationService transactionjpcbservice.OutstandingJobAllocationService) OutstandingJobAllocationController {
	return &OutstandingJobAllocationControllerImpl{
		OutstandingJobAllocationService: outstandingJobAllocationService,
	}
}

// @Summary Get All Outstanding Job Allocation
// @Description Get All Outstanding Job Allocation
// @Tags Transaction : JPCB Outstanding Job Allocation
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param booking_service_date query string false "Booking Service Date"
// @Param reference_document_type query string false "Reference Document Type"
// @Param reference_document_number query string false "Reference Document Number"
// @Param tnkb query string false "TNKB"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/outstanding-job-allocation [get]
func (r *OutstandingJobAllocationControllerImpl) GetAllOutstandingJobAllocation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":                queryValues.Get("company_id"),
		"booking_service_date":      queryValues.Get("booking_service_date"),
		"reference_document_type":   queryValues.Get("reference_document_type"),
		"reference_document_number": queryValues.Get("reference_document_number"),
		"tnkb":                      queryValues.Get("tnkb"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.OutstandingJobAllocationService.GetAllOutstandingJobAllocation(criteria, paginate)
	if err != nil {
		payloads.NewHandleSuccessPagination(writer, []interface{}{}, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, 0, 0)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Outstanding Job Allocation By Type ID
// @Description Get Outstanding Job Allocation By Type ID
// @Tags Transaction : JPCB Outstanding Job Allocation
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param reference_document_type path string true "Reference Document Type"
// @Param reference_system_number path string true "Reference System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/outstanding-job-allocation/{reference_document_type}/{reference_system_number} [get]
func (r *OutstandingJobAllocationControllerImpl) GetByTypeIdOutstandingJobAllocation(writer http.ResponseWriter, request *http.Request) {
	referenceDocumentType := strings.ToUpper(chi.URLParam(request, "reference_document_type"))
	referenceSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "reference_system_number"))

	result, err := r.OutstandingJobAllocationService.GetByTypeIdOutstandingJobAllocation(referenceDocumentType, referenceSystemNumber)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully", http.StatusOK)
}

// @Summary Save Outstanding Job Allocation
// @Description Save Outstanding Job Allocation
// @Tags Transaction : JPCB Outstanding Job Allocation
// @Security AuthorizationKeyAuth
// @Accept json
// @Produce json
// @Param reference_document_type path string true "Reference Document Type"
// @Param reference_system_number path string true "Reference System Number"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/outstanding-job-allocation/{reference_document_type}/{reference_system_number} [post]
func (r *OutstandingJobAllocationControllerImpl) SaveOutstandingJobAllocation(writer http.ResponseWriter, request *http.Request) {
	referenceDocumentType := strings.ToUpper(chi.URLParam(request, "reference_document_type"))
	referenceSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "reference_system_number"))
	formRequest := transactionjpcbpayloads.OutstandingJobAllocationSaveRequest{}

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

	update, err := r.OutstandingJobAllocationService.SaveOutstandingJobAllocation(referenceDocumentType, referenceSystemNumber, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Create Data Successfully!", http.StatusCreated)
}
