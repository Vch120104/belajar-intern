package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masteritemlevelservice "after-sales/api/services/master/item"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ItemLevelController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	ChangeStatus(writer http.ResponseWriter, request *http.Request)
	GetItemLevelDropDown(writer http.ResponseWriter, request *http.Request)
	GetItemLevelLookUp(writer http.ResponseWriter, request *http.Request)
	GetItemLevelLookUpbyId(writer http.ResponseWriter, request *http.Request)
}

type ItemLevelControllerImpl struct {
	itemLevelService masteritemlevelservice.ItemLevelService
}

func NewItemLevelController(ItemLevelService masteritemlevelservice.ItemLevelService) ItemLevelController {
	return &ItemLevelControllerImpl{
		itemLevelService: ItemLevelService,
	}
}

// @Summary Get Item Level Look Up By Id
// @Description Get Item Level Look Up By Id
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security AuthorizationKeyAuth
// @Param item_level_1_id path string true "item_level_1_id"
// @Param item_level_2_id query string false "item_level_2_id"
// @Param item_level_3_id query string false "item_level_3_id"
// @Param item_level_4_id query string false "item_level_4_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-level/look-up-item-level-by-id/{item_level_1_id} [get]
func (r *ItemLevelControllerImpl) GetItemLevelLookUpbyId(writer http.ResponseWriter, request *http.Request) {
	itemLevel, errA := strconv.Atoi(chi.URLParam(request, "item_level_1_id"))

	queryValues := request.URL.Query()
	filter := map[string]string{
		"mil2.item_level_2_id": queryValues.Get("item_level_2_id"),
		"mil3.item_level_3_id": queryValues.Get("item_level_3_id"),
		"mil4.item_level_4_id": queryValues.Get("item_level_4_id"),
	}

	internalCriteria := utils.BuildFilterCondition(filter)

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	get, err := r.itemLevelService.GetItemLevelLookUpbyId(internalCriteria, itemLevel)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Item Level Look Up
// @Description Get Item Level Look Up
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security AuthorizationKeyAuth
// @Param item_class_id path string true "item_class_id"
// @Param item_level_1_code query string false "Item Level 1 Code"
// @Param item_level_1_name query string false "Item Level 1 Name"
// @Param item_level_2_code query string false "Item Level 2 Code"
// @Param item_level_2_name query string false "Item Level 2 Name"
// @Param item_level_3_code query string false "Item Level 3 Code"
// @Param item_level_3_name query string false "Item Level 3 Name"
// @Param item_level_4_code query string false "Item Level 4 Code"
// @Param item_level_4_name query string false "Item Level 4 Name"
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-level/look-up-item-level/{item_class_id} [get]
func (r *ItemLevelControllerImpl) GetItemLevelLookUp(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	itemClassId, errA := strconv.Atoi(chi.URLParam(request, "item_class_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	filter := map[string]string{
		"mtr_item_level_1.item_level_1_code": queryValues.Get("item_level_1_code"),
		"mtr_item_level_1.item_level_1_name": queryValues.Get("item_level_1_name"),
		"mil2.item_level_2_code":             queryValues.Get("item_level_2_code"),
		"mil2.item_level_2_name":             queryValues.Get("item_level_2_name"),
		"mil3.item_level_3_code":             queryValues.Get("item_level_3_code"),
		"mil3.item_level_3_name":             queryValues.Get("item_level_3_name"),
		"mil4.item_level_4_code":             queryValues.Get("item_level_4_code"),
		"mil4.item_level_4_name":             queryValues.Get("item_level_4_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(filter)

	get, err := r.itemLevelService.GetItemLevelLookUp(internalCriteria, paginate, itemClassId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)

}

// @Summary Get Item Level Drop Down
// @Description Get Item Level Drop Down
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security AuthorizationKeyAuth
// @Param item_level path string true "item_level"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-level/drop-down-item-level/{item_level} [get]
func (r *ItemLevelControllerImpl) GetItemLevelDropDown(writer http.ResponseWriter, request *http.Request) {
	itemLevelId, errA := strconv.Atoi(chi.URLParam(request, "item_level"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	get, err := r.itemLevelService.GetItemLevelDropDown(itemLevelId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Get All Item Level
// @Description Get All Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Param item_level query string false "Item Level"
// @Param item_class_code query string false "Item Class Code"
// @Param item_level_parent query string false "Item Level Parent"
// @Param item_level_code query string false "Item Level Code"
// @Param item_level_name query string false "Item Level Name"
// @Param is_active query bool false "Is Active"
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-level/ [get]
func (r *ItemLevelControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filter := map[string]string{
		"item_level":        queryValues.Get("item_level"),
		"item_class_code":   queryValues.Get("item_class_code"),
		"item_level_parent": queryValues.Get("item_level_parent"),
		"item_level_code":   queryValues.Get("item_level_code"),
		"item_level_name":   queryValues.Get("item_level_name"),
		"is_active":         queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(filter)

	get, err := r.itemLevelService.GetAll(internalCriteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, get.Rows, "Get Data Successfully!", 200, get.Limit, get.Page, get.TotalRows, get.TotalPages)
}

// @Summary Get Item Level By Id
// @Description Get Item Level By Id
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security AuthorizationKeyAuth
// @Param item_level_id path string true "item_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-level/by-id/{item_level_id} [get]
func (r *ItemLevelControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {
	itemLevel, errA := strconv.Atoi(chi.URLParam(request, "item_level"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	itemLevelId, errB := strconv.Atoi(chi.URLParam(request, "item_level_id"))
	if errB != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	get, err := r.itemLevelService.GetById(itemLevel, itemLevelId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item Level
// @Description Save Item Level
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security AuthorizationKeyAuth
// @param reqBody body masteritemlevelpayloads.SaveItemLevelRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-level/ [post]
func (r *ItemLevelControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritemlevelpayloads.SaveItemLevelRequest
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message = ""

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	if formRequest.ItemLevel > 1 && formRequest.ItemLevelParent == 0 {
		payloads.NewHandleError(writer, "item_level_parent is required", http.StatusBadRequest)
		return
	}

	create, err := r.itemLevelService.Save(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.ItemLevelId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Item Level Status By Id
// @Description Change Item Level Status By Id
// @Accept json
// @Produce json
// @Tags Master : Item Level
// @Security AuthorizationKeyAuth
// @Param item_level_id path string true "item_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-level/{item_level_id} [patch]
func (r *ItemLevelControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request) {
	itemLevel, errA := strconv.Atoi(chi.URLParam(request, "item_level"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	itemLevelId, errB := strconv.Atoi(chi.URLParam(request, "item_level_id"))
	if errB != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.itemLevelService.ChangeStatus(itemLevel, int(itemLevelId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
