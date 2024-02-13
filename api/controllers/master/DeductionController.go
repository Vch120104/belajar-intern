package mastercontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DeductionController interface {
	GetAllDeductionList(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetByIdDeductionDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetByIdDeductionList(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveDeductionList(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveDeductionDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusDeduction(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type DeductionControllerImpl struct {
	DeductionService masterservice.DeductionService
}

func NewDeductionController(deductionService masterservice.DeductionService) DeductionController {
	return &DeductionControllerImpl{
		DeductionService: deductionService,
	}
}

// @Summary Get All Deduction
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param deduction_code query string false "deduction_code"
// @Param deduction_name query string false "deduction_name"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/deduction [get]
func (r *DeductionControllerImpl) GetAllDeductionList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"is_active":      queryValues.Get("is_active"),
		"deduction_name": queryValues.Get("deduction_name"),
		"effective_date": queryValues.Get("effective_date"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition((queryParams))

	result := r.DeductionService.GetAllDeduction(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Deduction Detail By Id
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @Param deduction_detail_id path int true "deduction_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/deduction/detail/by-id/{deduction_detail_id} [get]
func (r *DeductionControllerImpl) GetByIdDeductionDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	DeductionDetailIdstr := params.ByName("deduction_detail_id")

	DeductionListId, _ := strconv.Atoi(DeductionDetailIdstr)

	result := r.DeductionService.GetByIdDeductionDetail(DeductionListId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Deduction By Id
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @Param deduction_list_id path int true "deduction_list_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/deduction/header/by-id/{deduction_list_id} [get]
func (r *DeductionControllerImpl) GetByIdDeductionList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// DeductionListIdstr := params.ByName("deduction_list_id")
	// page := params.ByName("page")
	// limit := params.ByName("limit")

	// DeductionListId, _ := strconv.Atoi(DeductionListIdstr)
	// pageInt, _ := strconv.Atoi(page)
	// limitInt, _ := strconv.Atoi(limit)

	result := r.DeductionService.GetByIdDeductionList(1, 0, 0)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Deduction
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @param reqBody body masterpayloads.DeductionListResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/deduction [post]
func (r *DeductionControllerImpl) SaveDeductionList(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var formRequest masterpayloads.DeductionListResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.DeductionService.PostDeductionList(formRequest)

	if formRequest.DeductionListId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Save Deduction Detail
// @Description REST API Deduction Detail
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @param reqBody body masterpayloads.DeductionDetailResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/deduction-detail [post]
func (r *DeductionControllerImpl) SaveDeductionDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var formRequest masterpayloads.DeductionDetailResponse
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.DeductionService.PostDeductionDetail(formRequest)

	if formRequest.DeductionDetailId == 0 {
		message = "Create data Successfully!"
	} else {
		message = "Update data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Deduction
// @Description REST API Deduction
// @Accept json
// @Produce json
// @Tags Master : Deduction
// @param deduction_list_id path int true "deduction_list_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/deduction/{deduction_list_id} [patch]
func (r *DeductionControllerImpl) ChangeStatusDeduction(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	deductionId, _ := strconv.Atoi(params.ByName("deduction_list_id"))

	response := r.DeductionService.ChangeStatusDeduction(int(deductionId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
