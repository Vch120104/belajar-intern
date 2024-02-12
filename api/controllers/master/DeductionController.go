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
}

type DeductionControllerImpl struct {
	DeductionService masterservice.DeductionService
}

func NewDeductionController(deductionService masterservice.DeductionService) DeductionController {
	return &DeductionControllerImpl{
		DeductionService: deductionService,
	}
}

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

func (r *DeductionControllerImpl) GetByIdDeductionDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	DeductionListIdstr := params.ByName("deduction_list_id")

	DeductionListId, _ := strconv.Atoi(DeductionListIdstr)

	result := r.DeductionService.GetByIdDeductionDetail(DeductionListId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

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
