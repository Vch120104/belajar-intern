package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"errors"
	"time"

	helper "after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
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
	UpdateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemSubstitute(writer http.ResponseWriter, request *http.Request)
	ActivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
	DeactivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request)
	GetallItemForFilter(writer http.ResponseWriter, request *http.Request)
	GetItemSubstituteDetailLastSequence(writer http.ResponseWriter, request *http.Request)
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
// @Param substitute_type_id query string false "substitute_type_id"
// @Param substitute_type_code query string false "substitute_type_code"
// @Param item_id query string false "item_id"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/ [get]
func (r *ItemSubstituteControllerImpl) GetAllItemSubstitute(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":            queryValues.Get("is_active"),
		"substitute_type_id":   queryValues.Get("substitute_type_id"),
		"substitute_type_code": queryValues.Get("substitute_type_code"),
		"Item.item_id":         queryValues.Get("item_id"),
	}

	from, _ := time.Parse("2006-01-02T15:04:05.000Z", queryValues.Get("from"))
	to, _ := time.Parse("2006-01-02T15:04:05.000Z", queryValues.Get("to"))

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.ItemSubstituteService.GetAllItemSubstitute(filterCondition, pagination, from, to)

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

// @Summary Get Item Substitute By Id
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_id path int true "item_substitute_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/header/by-id/{item_substitute_id} [get]
func (r *ItemSubstituteControllerImpl) GetByIdItemSubstitute(writer http.ResponseWriter, request *http.Request) {
	ItemSubstituteIdStr := chi.URLParam(request, "item_substitute_id")

	ItemSubstituteId, errA := strconv.Atoi(ItemSubstituteIdStr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ItemSubstituteService.GetByIdItemSubstitute(ItemSubstituteId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Item Substitute Detail
// @Summary Get All Item Substitute Detail
// @Description REST API Item Substitute Detail
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_id path int true "item_substitute_id"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param substitute_type_code query string false "substitute_type_code"
// @Param item_id query string false "item_id"
// @Param effective_date query string false "effective_date"
// @Param is_active query string false "is_active" Enums(true,false)
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/detail/all/by-id/{item_substitute_id} [get]
func (r *ItemSubstituteControllerImpl) GetAllItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	ItemSubstituteIdStr := chi.URLParam(request, "item_substitute_id")

	ItemSubstituteId, errA := strconv.Atoi(ItemSubstituteIdStr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.ItemSubstituteService.GetAllItemSubstituteDetail(pagination, ItemSubstituteId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Item Substitute Detail By Id
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_detail_id path int true "item_substitute_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/detail/by-id/{item_substitute_detail_id} [get]
func (r *ItemSubstituteControllerImpl) GetByIdItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	ItemSubstituteDetailIdStr := chi.URLParam(request, "item_substitute_detail_id")

	ItemSubstituteDetailId, errA := strconv.Atoi(ItemSubstituteDetailIdStr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ItemSubstituteService.GetByIdItemSubstituteDetail(ItemSubstituteDetailId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param reqBody body masteritempayloads.ItemSubstitutePostPayloads true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/ [post]
func (r *ItemSubstituteControllerImpl) SaveItemSubstitute(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemSubstitutePostPayloads
	var message = ""

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemSubstituteService.SaveItemSubstitute(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/detail/{item_substitute_id} [post]
func (r *ItemSubstituteControllerImpl) SaveItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemSubstituteDetailPostPayloads
	ItemSubstituteDetailIdStr := chi.URLParam(request, "item_substitute_id")

	ItemSubstituteDetailId, errA := strconv.Atoi(ItemSubstituteDetailIdStr)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	var message = ""
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemSubstituteService.SaveItemSubstituteDetail(formRequest, ItemSubstituteDetailId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.ItemSubstituteDetailId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Update Item Substitute Detail
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param reqBody body masteritempayloads.ItemSubstituteDetailUpdatePayloads true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/detail [put]
func (r *ItemSubstituteControllerImpl) UpdateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemSubstituteDetailUpdatePayloads
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemSubstituteService.UpdateItemSubstituteDetail(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param item_substitute_id path int true "item_substitute_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/{item_substitute_id} [patch]
func (r *ItemSubstituteControllerImpl) ChangeStatusItemSubstitute(writer http.ResponseWriter, request *http.Request) {

	ItemSubstituteId, errA := strconv.Atoi(chi.URLParam(request, "item_substitute_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemSubstituteService.ChangeStatusItemSubstitute(ItemSubstituteId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param item_substitute_detail_id path int true "item_substitute_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/detail/activate/by-id/{item_substitute_detail_id} [patch]
func (r *ItemSubstituteControllerImpl) ActivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {

	queryId := chi.URLParam(request, "item_substitute_detail_id")

	response, err := r.ItemSubstituteService.ActivateItemSubstituteDetail(queryId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Change Status Item Substitute
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @param item_substitute_detail_id path int true "item_substitute_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/detail/deactivate/by-id/{item_substitute_detail_id} [patch]
func (r *ItemSubstituteControllerImpl) DeactivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request) {

	queryId := chi.URLParam(request, "item_substitute_detail_id")

	response, err := r.ItemSubstituteService.DeactivateItemSubstituteDetail(queryId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get All Item For Filter
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_code query string false "item_code"
// @Param item_name query string false "item_name"
// @Param item_class query string false "item_class"
// @Param item_type_code query string false "item_type_code"
// @Param item_level_1_code query string false "item_level_1_code"
// @Param item_level_2_code query string false "item_level_2_code"
// @Param item_level_3_code query string false "item_level_3_code"
// @Param item_level_4_code query string false "item_level_4_code"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute//item-for-substitute [get]
func (r *ItemSubstituteControllerImpl) GetallItemForFilter(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"item_code":         queryValues.Get("item_code"),
		"item_name":         queryValues.Get("item_name"),
		"item_class":        queryValues.Get("item_class"),
		"item_type_code":    queryValues.Get("item_type_code"),
		"item_level_1_code": queryValues.Get("item_level_1_code"),
		"item_level_2_code": queryValues.Get("item_level_2_code"),
		"item_level_3_code": queryValues.Get("item_level_3_code"),
		"item_level_4_code": queryValues.Get("item_level_4_code"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.ItemSubstituteService.GetallItemForFilter(filterCondition, pagination)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Item Substitute Detail Last Sequence
// @Description REST API Item Substitute
// @Accept json
// @Produce json
// @Tags Master : Item Substitute
// @Param item_substitute_id path int true "item_substitute_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-substitute/detail/last-sequence/{item_substitute_id} [get]
func (r *ItemSubstituteControllerImpl) GetItemSubstituteDetailLastSequence(writer http.ResponseWriter, request *http.Request) {
	itemSubstituteId, errA := strconv.Atoi(chi.URLParam(request, "item_substitute_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read url params, please check your param input")})
		return
	}

	result, err := r.ItemSubstituteService.GetItemSubstituteDetailLastSequence(itemSubstituteId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
