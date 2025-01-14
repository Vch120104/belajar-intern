package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AtpmClaimRegistrationControllerImpl struct {
	AtpmClaimRegistrationService transactionworkshopservice.AtpmClaimRegistrationService
}

type AtpmClaimRegistrationController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
}

func NewAtpmClaimRegistrationController(AtpmClaimRegistrationService transactionworkshopservice.AtpmClaimRegistrationService) AtpmClaimRegistrationController {
	return &AtpmClaimRegistrationControllerImpl{
		AtpmClaimRegistrationService: AtpmClaimRegistrationService,
	}
}

// GetAll gets all atpm claim registrations
// @Summary Get All ATPM Claim Registrations
// @Description Retrieve all atpm claim registrations with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Transaction : Workshop ATPM Claim Registration
// @Param atpm_claim_registration_number query string false "ATPM Claim Registration Number"
// @Param atpm_claim_registration_document_number query string false "ATPM Claim Registration Document Number"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Field to sort by"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/atpm-claim-registration [get]
func (r *AtpmClaimRegistrationControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {

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

	result, err := r.AtpmClaimRegistrationService.GetAll(criteria, paginate)
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

// GetById gets atpm claim registration by id
// @Summary Get ATPM Claim Registration By ID
// @Description Retrieve atpm claim registration by id
// @Accept json
// @Produce json
// @Tags Transaction : Workshop ATPM Claim Registration
// @Param claim_system_number path string true "ATPM Claim Registration ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/atpm-claim-registration/{claim_system_number} [get]
func (r *AtpmClaimRegistrationControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {

	claimSystemNumberStr := chi.URLParam(request, "claim_system_number")
	claimSystemNumber, err := strconv.Atoi(claimSystemNumberStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid claim system number", http.StatusBadRequest)
		return
	}

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, baseErr := r.AtpmClaimRegistrationService.GetById(claimSystemNumber, paginate)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Atpm Claim Registration not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
