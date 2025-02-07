package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
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

type ItemTypeController interface {
	GetAllItemType(writer http.ResponseWriter, request *http.Request)
	GetItemTypeById(writer http.ResponseWriter, request *http.Request)
	GetItemTypeByCode(writer http.ResponseWriter, request *http.Request)
	CreateItemType(writer http.ResponseWriter, request *http.Request)
	SaveItemType(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemType(writer http.ResponseWriter, request *http.Request)
	GetItemTypeDropDown(writer http.ResponseWriter, request *http.Request)
}

type ItemTypeControllerImpl struct {
	ItemTypeService masteritemservice.ItemTypeService
}

func NewItemTypeController(itemTypeService masteritemservice.ItemTypeService) ItemTypeController {
	return &ItemTypeControllerImpl{
		ItemTypeService: itemTypeService,
	}
}

// @Summary Get All ItemType
// @Description Get All ItemType
// @Tags Master Item : Item Type
// @Accept json
// @Produce json
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-type/drop-down [get]
func (r *ItemTypeControllerImpl) GetItemTypeDropDown(writer http.ResponseWriter, request *http.Request) {
	response, err := r.ItemTypeService.GetItemTypeDropDown()

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get ItemType By Code
// @Description Get ItemType By Code
// @Tags Master Item : Item Type
// @Accept json
// @Produce json
// @Param item_type_code path string true "Item Type Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-type/code/{item_type_code} [get]
func (r *ItemTypeControllerImpl) GetItemTypeByCode(writer http.ResponseWriter, request *http.Request) {
	itemTypeCode := chi.URLParam(request, "item_type_code")

	response, err := r.ItemTypeService.GetItemTypeByCode(itemTypeCode)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get ItemType By Id
// @Description Get ItemType By Id
// @Tags Master Item : Item Type
// @Accept json
// @Produce json
// @Param item_type_id path string true "Item Type Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-type/{item_type_id} [get]
func (r *ItemTypeControllerImpl) GetItemTypeById(writer http.ResponseWriter, request *http.Request) {
	itemTypeId, errA := strconv.Atoi(chi.URLParam(request, "item_type_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemTypeService.GetItemTypeById(itemTypeId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All ItemType
// @Description Get All ItemType
// @Tags Master Item : Item Type
// @Accept json
// @Produce json
// @Param item_type_id query string false "Item Type Id"
// @Param item_type_code query string false "Item Type Code"
// @Param item_type_name query string false "Item Type Name"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-type [get]
func (r *ItemTypeControllerImpl) GetAllItemType(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"item_type_id":   queryValues.Get("item_type_id"),
		"item_type_code": queryValues.Get("item_type_code"),
		"item_type_name": queryValues.Get("item_type_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.ItemTypeService.GetAllItemType(criteria, paginate)

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

// @Summary Create ItemType
// @Description Create ItemType
// @Tags Master Item : Item Type
// @Accept json
// @Produce json
// @param reqBody body masteritempayloads.ItemTypeRequest true "Form Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-type [post]
func (r *ItemTypeControllerImpl) CreateItemType(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemTypeRequest

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemTypeService.CreateItemType(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message := "Data created successfully!"
	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}

// @Summary Save ItemType
// @Description Save ItemType
// @Tags Master Item : Item Type
// @Accept json
// @Produce json
// @Param item_type_id path string true "Item Type Id"
// @param reqBody body masteritempayloads.ItemTypeRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-type/{item_type_id} [put]
func (r *ItemTypeControllerImpl) SaveItemType(writer http.ResponseWriter, request *http.Request) {
	itemTypeID, errA := strconv.Atoi(chi.URLParam(request, "item_type_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var formRequest masteritempayloads.ItemTypeRequest

	// Parse request body to formRequest
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemTypeService.SaveItemType(itemTypeID, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message := "Data updated successfully!"
	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status ItemType
// @Description Change Status ItemType
// @Tags Master Item : Item Type
// @Accept json
// @Produce json
// @Param item_type_id path string true "Item Type Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-type/{item_type_id} [patch]
func (r *ItemTypeControllerImpl) ChangeStatusItemType(writer http.ResponseWriter, request *http.Request) {
	itemTypeId, errA := strconv.Atoi(chi.URLParam(request, "item_type_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemTypeService.ChangeStatusItemType(itemTypeId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Data Successfully!", http.StatusOK)
}
