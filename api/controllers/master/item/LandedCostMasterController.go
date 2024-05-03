package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LandedCostMasterController interface {
	GetAllLandedCostMaster(writer http.ResponseWriter, request *http.Request)
	GetByIdLandedCost(writer http.ResponseWriter, request *http.Request)
	SaveLandedCostMaster(writer http.ResponseWriter, request *http.Request)
	ActivateLandedCostMaster(writer http.ResponseWriter, request *http.Request)
	DeactivateLandedCostmaster(writer http.ResponseWriter, request *http.Request)
}

type LandedCostMasterControllerImpl struct {
	LandedCostService masteritemservice.LandedCostMasterService
}

func NewLandedCostMasterController(LandedCostService masteritemservice.LandedCostMasterService) LandedCostMasterController {
	return &LandedCostMasterControllerImpl{
		LandedCostService: LandedCostService,
	}
}

// @Summary Get All Landed Cost Master
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master : Landed Cost Master
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param company_id query string false "company_id"
// @Param supplier_id query string false "supplier_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/landed-cost/ [get]
func (r *LandedCostMasterControllerImpl) GetAllLandedCostMaster(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"company_id":  queryValues.Get("company_id"),
		"supplier_id": queryValues.Get("supplier_id"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.LandedCostService.GetAllLandedCost(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Landed Cost Master By Id
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master : Landed Cost Master
// @Param landed_cost_id path int true "landed_cost_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/landed-cost/{landed_cost_id} [get]
func (r *LandedCostMasterControllerImpl) GetByIdLandedCost(writer http.ResponseWriter, request *http.Request) {
	LandedCostIdstr := chi.URLParam(request, "landed_cost_id")

	LandedCostId, _ := strconv.Atoi(LandedCostIdstr)

	result := r.LandedCostService.GetByIdLandedCost(LandedCostId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Landed Cost Master
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master : Landed Cost Master
// @param reqBody body masteritempayloads.LandedCostMasterPayloads true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/landed-cost/ [post]
func (r *LandedCostMasterControllerImpl) SaveLandedCostMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.LandedCostMasterPayloads
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.LandedCostService.SaveLandedCost(formRequest)

	if formRequest.LandedCostId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Landed Cost Master
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master : Landed Cost Master
// @param landed_cost_id path int true "landed_cost_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/landed-cost/activate/{landed_cost_id}  [patch]
func (r *LandedCostMasterControllerImpl) ActivateLandedCostMaster(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	queryId := query.Get("landed_cost_id")
	response := r.LandedCostService.ActivateLandedCostMaster(queryId)
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Landed Cost Master
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master : Landed Cost Master
// @param landed_cost_id path int true "landed_cost_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/landed-cost/deactivate/{landed_cost_id} [patch]
func (r *LandedCostMasterControllerImpl) DeactivateLandedCostmaster(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	queryId := query.Get("landed_cost_id")
	response := r.LandedCostService.DeactivateLandedCostMaster(queryId)
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
