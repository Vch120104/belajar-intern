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

type ItemSubstituteController interface {
	GetAllItemSubstitute(writer http.ResponseWriter, request *http.Request)
	GetByIdItemSubstitute(writer http.ResponseWriter, request *http.Request)
	GetAllItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
	GetByIdItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
	SaveItemSubstitute(writer http.ResponseWriter, request *http.Request)
	SaveItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemSubstitute(writer http.ResponseWriter, request *http.Request)
	ActivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
	DeactivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
}

type ItemSubstituteControllerImpl struct {
	ItemSubstituteService masteritemservice.ItemSubstituteService
}

func NewItemSubstituteController(itemSubstituteService masteritemservice.ItemSubstituteService) ItemSubstituteController {
	return &ItemSubstituteControllerImpl{
		ItemSubstituteService: itemSubstituteService,
	}
}

// @Summary Get All Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param substitute_type_code query string false "substitute_type_code"
// @Param item_id query string false "item_id"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router / [get]
func (r *ItemSubstituteControllerImpl) GetAllItemSubstitute(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":            queryValues.Get("is_active"),
		"substitute_type_code": queryValues.Get("substitute_type_code"),
		"effective_date":       queryValues.Get("effective_date"),
		"item_id":              queryValues.Get("item_id"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.ItemSubstituteService.GetAllItemSubstitute(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Item Substitute By Id
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_id path int true "item_substitute_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /header/by-id/{item_substitute_id} [get]
func (r *ItemSubstituteControllerImpl) GetByIdItemSubstitute(writer http.ResponseWriter, request *http.Request) {
	ItemSubstituteIdStr := chi.URLParam(request, "item_substitute_id")

	ItemSubstituteId, _ := strconv.Atoi(ItemSubstituteIdStr)

	result := r.ItemSubstituteService.GetByIdItemSubstitute(ItemSubstituteId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Item Substitute Detail
// @Description REST API Item Substitute Detail
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param substitute_type_code query string false "substitute_type_code"
// @Param item_id query string false "item_id"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /detail/all/by-id/{item_substitute_id} [get]
func (r *ItemSubstituteControllerImpl) GetAllItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	ItemSubstituteIdStr := chi.URLParam(request, "item_substitute_id")

	ItemSubstituteId, _ := strconv.Atoi(ItemSubstituteIdStr)
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result := r.ItemSubstituteService.GetAllItemSubstituteDetail(pagination, ItemSubstituteId)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Item Substitute Detail By Id
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_detail_id path int true "item_substitute_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /detail/by-id/{item_substitute_detail_id} [get]
func (r *ItemSubstituteControllerImpl) GetByIdItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	ItemSubstituteDetailIdStr := chi.URLParam(request, "item_substitute_detail_id")

	ItemSubstituteDetailId, _ := strconv.Atoi(ItemSubstituteDetailIdStr)

	result := r.ItemSubstituteService.GetByIdItemSubstituteDetail(ItemSubstituteDetailId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param reqBody body masteritempayloads.ItemSubstitutePostPayloads true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router / [post]
func (r *ItemSubstituteControllerImpl) SaveItemSubstitute(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemSubstitutePostPayloads
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.ItemSubstituteService.SaveItemSubstitute(formRequest)

	if formRequest.ItemSubstituteId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Save Item Substitute Detail
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_id path int true "item_substitute_id"
// @param reqBody body masteritempayloads.ItemSubstituteDetailPostPayloads true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /detail/{item_substitute_id} [post]
func (r *ItemSubstituteControllerImpl) SaveItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemSubstituteDetailPostPayloads
	ItemSubstituteDetailIdStr := chi.URLParam(request, "item_substitute_id")

	ItemSubstituteDetailId, _ := strconv.Atoi(ItemSubstituteDetailIdStr)
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.ItemSubstituteService.SaveItemSubstituteDetail(formRequest, ItemSubstituteDetailId)

	if formRequest.ItemSubstituteDetailId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param item_substitute_id path int true "item_substitute_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /{item_substitute_id} [patch]
func (r *ItemSubstituteControllerImpl) ChangeStatusItemSubstitute(writer http.ResponseWriter, request *http.Request) {

	ItemSubstituteId, _ := strconv.Atoi(chi.URLParam(request, "item_substitute_id"))

	response := r.ItemSubstituteService.ChangeStatusItemSubstitute(ItemSubstituteId)

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param item_substitute_detail_id path int true "item_substitute_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router/detail/activate/by-id/ [patch]
func (r *ItemSubstituteControllerImpl) ActivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	queryId := query.Get("item_substitute_detail_id")
	response := r.ItemSubstituteService.ActivateItemSubstituteDetail(queryId)
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param item_substitute_detail_id path int true "item_substitute_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /detail/deactivate/by-id/ [patch]
func (r *ItemSubstituteControllerImpl) DeactivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	queryId := query.Get("item_substitute_detail_id")
	response := r.ItemSubstituteService.DeactivateItemSubstituteDetail(queryId)
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
