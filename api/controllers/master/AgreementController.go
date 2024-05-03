package mastercontroller

import (

	// "after-sales/api/middlewares"

	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AgreementController interface {
	GetAgreementById(writer http.ResponseWriter, request *http.Request)
	SaveAgreement(writer http.ResponseWriter, request *http.Request)
	ChangeStatusAgreement(writer http.ResponseWriter, request *http.Request)
	GetAllAgreement(writer http.ResponseWriter, request *http.Request)
}
type AgreementControllerImpl struct {
	AgreementService masterservice.AgreementService
}

func NewAgreementController(AgreementService masterservice.AgreementService) AgreementController {
	return &AgreementControllerImpl{
		AgreementService: AgreementService,
	}
}

// @Summary Get Agreement By Id
// @Description Retrieve an agreement by its ID
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/agreement/{agreement_id} [get]
func (r *AgreementControllerImpl) GetAgreementById(writer http.ResponseWriter, request *http.Request) {

	AgreementId, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))

	result, err := r.AgreementService.GetAgreementById(int(AgreementId))
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Agreement
// @Description Create or update an agreement
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param reqBody body masterpayloads.AgreementResponse true "Agreement Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/agreement/ [post]
func (r *AgreementControllerImpl) SaveAgreement(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.AgreementResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create, err := r.AgreementService.SaveAgreement(formRequest)
	if err != nil {
		exceptionsss_test.NewConflictException(writer, request, err)
		return
	}

	if formRequest.AgreementId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Agreement
// @Description Change the status of an agreement
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param agreement_id path int true "Agreement ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/agreement/{agreement_id} [patch]
func (r *AgreementControllerImpl) ChangeStatusAgreement(writer http.ResponseWriter, request *http.Request) {

	agreement_id, _ := strconv.Atoi(chi.URLParam(request, "agreement_id"))

	response, err := r.AgreementService.ChangeStatusAgreement(int(agreement_id))
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get All Agreements
// @Description Retrieve all agreements with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Master : Agreement
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param is_active query string false "Agreement status"
// @Param supplier_name query string false "Supplier name"
// @Param agreement_code query string false "Agreement code"
// @Param brand_id query string false "Brand ID"
// @Param customer_id query string false "Customer ID"
// @Param profit_center_id query string false "Profit center ID"
// @Param dealer_id query string false "Dealer ID"
// @Param top_id query string false "Top ID"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/agreement/ [get]

func (r *AgreementControllerImpl) GetAllAgreement(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"supplier_name": queryValues.Get("supplier_name"),
		"mtr_moving_code.moving_code_description":           queryValues.Get("moving_code_description"),
		"order_type_name":                                   queryValues.Get("order_type_name"),
		"mtr_forecast_master.forecast_master_lead_time":     queryValues.Get("forecast_master_lead_time"),
		"mtr_forecast_master.forecast_master_safety_factor": queryValues.Get("forecast_master_safety_factor"),
		"mtr_forecast_master.forecast_master_order_cycle":   queryValues.Get("forecast_master_order_cycle"),
		"mtr_forecast_master.is_active":                     queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}
	print(queryParams)

	criteria := utils.BuildFilterCondition(queryParams)
	paginatedData, totalPages, totalRows, err := r.AgreementService.GetAllAgreement(criteria, paginate)

	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}
