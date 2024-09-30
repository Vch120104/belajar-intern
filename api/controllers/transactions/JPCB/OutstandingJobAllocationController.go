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

	payloads.NewHandleSuccess(writer, update, "Update Data Successfully!", http.StatusCreated)
}
