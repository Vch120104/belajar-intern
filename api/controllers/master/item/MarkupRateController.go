package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	masteritemservice "after-sales/api/services/master/item"

	// "after-sales/api/middlewares"

	// "strconv"

	"github.com/julienschmidt/httprouter"
)

type MarkupRateController interface {
	GetAllMarkupRate(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetMarkupRateByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveMarkupRate(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusMarkupRate(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
type MarkupRateControllerImpl struct {
	MarkupRateService masteritemservice.MarkupRateService
}

func NewMarkupRateController(markupRateService masteritemservice.MarkupRateService) MarkupRateController {
	return &MarkupRateControllerImpl{
		MarkupRateService: markupRateService,
	}
}

// @Summary Get All Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param markup_master_code query string false "markup_master_code"
// @Param markup_master_description query string false "markup_master_description"
// @Param order_type_name query string false "order_type_name"
// @Param markup_rate query float64 false "markup_rate"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate [get]
func (r *MarkupRateControllerImpl) GetAllMarkupRate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_markup_master.markup_master_code":        queryValues.Get("markup_master_code"),
		"mtr_markup_master.markup_master_description": queryValues.Get("markup_master_description"),
		"order_type_name":                             queryValues.Get("order_type_name"),
		"mtr_markup_rate.markup_rate":                 queryValues.Get("markup_rate"),
		"mtr_markup_rate.is_active":                   queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows := r.MarkupRateService.GetAllMarkupRate(criteria, paginate)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Markup Rate By ID
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @Param markup_rate_id path int true "markup_rate_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate/{markup_rate_id} [get]
func (r *MarkupRateControllerImpl) GetMarkupRateByID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	markupRateId, _ := strconv.Atoi(params.ByName("markup_rate_id"))

	result := r.MarkupRateService.GetMarkupRateById(markupRateId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @param reqBody body masteritempayloads.MarkupRateRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate [post]
func (r *MarkupRateControllerImpl) SaveMarkupRate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var formRequest masteritempayloads.MarkupRateRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.MarkupRateService.SaveMarkupRate(formRequest)

	if formRequest.MarkupRateId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Markup Rate
// @Description REST API Markup Rate
// @Accept json
// @Produce json
// @Tags Master : Markup Rate
// @param markup_rate_id path int true "markup_rate_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.Error
// @Router /aftersales-service/api/aftersales/markup-rate/{markup_rate_id} [patch]
func (r *MarkupRateControllerImpl) ChangeStatusMarkupRate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	markupRateId, _ := strconv.Atoi(params.ByName("markup_rate_id"))

	response := r.MarkupRateService.ChangeStatusMarkupRate(int(markupRateId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
