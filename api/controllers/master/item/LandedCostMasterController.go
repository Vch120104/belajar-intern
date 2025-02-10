package masteritemcontroller

import (
	"after-sales/api/exceptions"
	helper "after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
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
	UpdateLandedCostMaster(writer http.ResponseWriter, request *http.Request)
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
// @Tags Master Item : Landed Cost Master
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param company_id query string false "company_id"
// @Param supplier_id query string false "supplier_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
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

	result, err := r.LandedCostService.GetAllLandedCost(filterCondition, pagination)
	if err != nil {
		helper.ReturnError(writer, request, err)
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

// @Summary Get Landed Cost Master By Id
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master Item : Landed Cost Master
// @Param landed_cost_id path int true "landed_cost_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/landed-cost/{landed_cost_id} [get]
func (r *LandedCostMasterControllerImpl) GetByIdLandedCost(writer http.ResponseWriter, request *http.Request) {
	LandedCostIdstr := chi.URLParam(request, "landed_cost_id")

	LandedCostId, errA := strconv.Atoi(LandedCostIdstr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.LandedCostService.GetByIdLandedCost(LandedCostId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Landed Cost Master
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master Item : Landed Cost Master
// @param reqBody body masteritempayloads.LandedCostMasterPayloads true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/landed-cost/ [post]
func (r *LandedCostMasterControllerImpl) SaveLandedCostMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.LandedCostMasterRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	var message = ""

	create, err := r.LandedCostService.SaveLandedCost(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

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
// @Tags Master Item : Landed Cost Master
// @param landed_cost_id path int true "landed_cost_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/landed-cost/activate/{landed_cost_id}  [patch]
func (r *LandedCostMasterControllerImpl) ActivateLandedCostMaster(writer http.ResponseWriter, request *http.Request) {
	queryId := chi.URLParam(request, "landed_cost_id")
	response, err := r.LandedCostService.ActivateLandedCostMaster(queryId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Landed Cost Master
// @Description REST API Landed Cost Master
// @Accept json
// @Produce json
// @Tags Master Item : Landed Cost Master
// @param landed_cost_id path int true "landed_cost_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/landed-cost/deactivate/{landed_cost_id} [patch]
func (r *LandedCostMasterControllerImpl) DeactivateLandedCostmaster(writer http.ResponseWriter, request *http.Request) {
	queryId := chi.URLParam(request, "landed_cost_id")
	response, err := r.LandedCostService.DeactivateLandedCostMaster(queryId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

func (r *LandedCostMasterControllerImpl) UpdateLandedCostMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.LandedCostMasterUpdateRequest
	query, errA := strconv.Atoi(chi.URLParam(request, "landed_cost_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	helper.ReadFromRequestBody(request, &formRequest)

	update, err := r.LandedCostService.UpdateLandedCostMaster(query, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, update, "Update Data Successfully!", http.StatusOK)
}
