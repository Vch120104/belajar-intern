package transactionworkshopcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
)

type AtpmReimbursementControllerImpl struct {
	AtpmReimbursementService transactionworkshopservice.AtpmReimbursementService
}

type AtpmReimbursementController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
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
