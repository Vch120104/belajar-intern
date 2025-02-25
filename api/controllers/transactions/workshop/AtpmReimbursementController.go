package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AtpmReimbursementControllerImpl struct {
	AtpmReimbursementService transactionworkshopservice.AtpmReimbursementService
}

type AtpmReimbursementController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	New(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	Submit(writer http.ResponseWriter, request *http.Request)
}

func NewAtpmReimbursementController(AtpmReimbursementService transactionworkshopservice.AtpmReimbursementService) AtpmReimbursementController {
	return &AtpmReimbursementControllerImpl{
		AtpmReimbursementService: AtpmReimbursementService,
	}
}

// GetAll gets all atpm reimburesement
// @Summary Get All ATPM Reimbursement
// @Description Retrieve all atpm Reimbursement with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop ATPM Reimbursement
// @Security AuthorizationKeyAuth
// @Param atpm_reimbursement_number query string false "ATPM Reimbursement Number"
// @Param atpm_reimbursement_document_number query string false "ATPM Reimbursement Document Number"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/atpm-reimbursement [get]
func (r *AtpmReimbursementControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_atpm_claim_registration.atpm_claim_registration_number":          queryValues.Get("atpm_claim_registration_number"),
		"trx_atpm_claim_registration.atpm_claim_registration_document_number": queryValues.Get("atpm_claim_registration_document_number"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.AtpmReimbursementService.GetAll(criteria, paginate)
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

// New creates new atpm reimbursement
// @Summary Create New ATPM Reimbursement
// @Description Create new atpm Reimbursement
// @Accept json
// @Produce json
// @Tags Transaction : Workshop ATPM Reimbursement
// @Security AuthorizationKeyAuth
// @Param body body transactionworkshoppayloads.AtpmReimbursementRequest true "Atpm Reimbursement Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/atpm-reimbursement [post]
func (r *AtpmReimbursementControllerImpl) New(writer http.ResponseWriter, request *http.Request) {

	var req transactionworkshoppayloads.AtpmReimbursementRequest
	helper.ReadFromRequestBody(request, &req)
	if validationErr := validation.ValidationForm(writer, request, &req); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	result, err := r.AtpmReimbursementService.New(req)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Data has been created successfully!", http.StatusCreated)
}

// Save saves atpm reimbursement
// @Summary Save ATPM Reimbursement
// @Description Save atpm Reimbursement
// @Accept json
// @Produce json
// @Tags Transaction : Workshop ATPM Reimbursement
// @Security AuthorizationKeyAuth
// @Param claim_system_number path int true "ATPM Reimbursement ID"
// @Param body body transactionworkshoppayloads.AtpmReimbursementUpdate true "Atpm Reimbursement Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/atpm-reimbursement/{claim_system_number} [put]
func (r *AtpmReimbursementControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {

	ClaimSystemNumberStrId := chi.URLParam(request, "claim_system_number")
	ClaimSystemNumber, err := strconv.Atoi(ClaimSystemNumberStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid claim system number ID", http.StatusBadRequest)
		return
	}

	var req transactionworkshoppayloads.AtpmReimbursementUpdate
	helper.ReadFromRequestBody(request, &req)
	if validationErr := validation.ValidationForm(writer, request, &req); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	result, baseErr := r.AtpmReimbursementService.Save(ClaimSystemNumber, req)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Data has been updated successfully!", http.StatusOK)
}

// Submit submits atpm reimbursement
// @Summary Submit ATPM Reimbursement
// @Description Submit atpm Reimbursement
// @Accept json
// @Produce json
// @Tags Transaction : Workshop ATPM Reimbursement
// @Security AuthorizationKeyAuth
// @Param claim_system_number path int true "ATPM Reimbursement ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/atpm-reimbursement/{claim_system_number}/submit [patch]
func (r *AtpmReimbursementControllerImpl) Submit(writer http.ResponseWriter, request *http.Request) {

	ClaimSystemNumberStrId := chi.URLParam(request, "claim_system_number")
	ClaimSystemNumber, err := strconv.Atoi(ClaimSystemNumberStrId)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid claim system number ID", http.StatusBadRequest)
		return
	}

	result, baseErr := r.AtpmReimbursementService.Submit(ClaimSystemNumber)
	if baseErr != nil {
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Data has been submitted successfully!", http.StatusOK)
}
