package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
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

type ItemPriceCodeController interface {
	GetAllItemPriceCode(writer http.ResponseWriter, request *http.Request)
	GetItemPriceCodeById(writer http.ResponseWriter, request *http.Request)
	GetItemPriceCodeByCode(writer http.ResponseWriter, request *http.Request)
	SaveItemPriceCode(writer http.ResponseWriter, request *http.Request)
	DeleteItemPriceCode(writer http.ResponseWriter, request *http.Request)
	UpdateItemPriceCode(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemPriceCode(writer http.ResponseWriter, request *http.Request)
	GetItemPriceCodeDropDown(writer http.ResponseWriter, request *http.Request)
}

type ItemPriceCodeControllerImpl struct {
	ItemPriceCodeService masteritemservice.ItemPriceCodeService
}

func NewItemPriceCodeController(itemPriceCodeService masteritemservice.ItemPriceCodeService) ItemPriceCodeController {
	return &ItemPriceCodeControllerImpl{
		ItemPriceCodeService: itemPriceCodeService,
	}
}

// @Summary Get All Item Price Code
// @Description Get all item price code
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Param is_active query string false "Is Active"
// @Param item_price_code query string false "Item Price Code"
// @Param item_price_code_id query string false "Item Price Code ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_by query string false "Sort By"
// @Param sort_of query string false "Sort Of"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code [get]
func (r *ItemPriceCodeControllerImpl) GetAllItemPriceCode(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filterParams := map[string]string{
		"mtr_item_price_code.is_active":          queryValues.Get("is_active"),
		"mtr_item_price_code.item_price_code":    queryValues.Get("item_price_code"),
		"mtr_item_price_code.item_price_code_id": queryValues.Get("item_price_code_id"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortBy: queryValues.Get("sort_by"),
		SortOf: queryValues.Get("sort_of"),
	}

	filter := utils.BuildFilterCondition(filterParams)

	result, err := r.ItemPriceCodeService.GetAllItemPriceCode(filter, pagination)
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

// @Summary Get Item Price Code By ID
// @Description Get item price code by ID
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Param item_price_code_id path int true "Item Price Code ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code/{item_price_code_id} [get]
func (r *ItemPriceCodeControllerImpl) GetItemPriceCodeById(writer http.ResponseWriter, request *http.Request) {
	itemPriceCodeId, err := strconv.Atoi(chi.URLParam(request, "item_price_code_id"))
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("invalid item_price_code_id"),
		})
		return
	}

	result, errResult := r.ItemPriceCodeService.GetByIdItemPriceCode(itemPriceCodeId)
	if errResult != nil {
		exceptions.NewNotFoundException(writer, request, errResult)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Item Price Code By Code
// @Description Get item price code by code
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Param item_price_code path string true "Item Price Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code/by-code/{item_price_code} [get]
func (r *ItemPriceCodeControllerImpl) GetItemPriceCodeByCode(writer http.ResponseWriter, request *http.Request) {
	itemPriceCode := chi.URLParam(request, "item_price_code")

	result, err := r.ItemPriceCodeService.GetByCodeItemPriceCode(itemPriceCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item Price Code
// @Description Save item price code
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Param body body masteritempayloads.SaveItemPriceCode true "Save Item Price Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code [post]
func (r *ItemPriceCodeControllerImpl) SaveItemPriceCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.SaveItemPriceCode
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	result, err := r.ItemPriceCodeService.SaveItemPriceCode(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message := "Item Price Code saved successfully!"
	payloads.NewHandleSuccess(writer, result, message, http.StatusOK)
}

// @Summary Delete Item Price Code
// @Description Delete item price code
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Param item_price_code_id path int true "Item Price Code ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code/{item_price_code_id} [delete]
func (r *ItemPriceCodeControllerImpl) DeleteItemPriceCode(writer http.ResponseWriter, request *http.Request) {
	itemPriceCodeId, err := strconv.Atoi(chi.URLParam(request, "item_price_code_id"))
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("invalid item_price_code_id"),
		})
		return
	}

	result, errResult := r.ItemPriceCodeService.DeleteItemPriceCode(itemPriceCodeId)
	if errResult != nil {
		helper.ReturnError(writer, request, errResult)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Item Price Code deleted successfully!", http.StatusOK)
}

// @Summary Update Item Price Code
// @Description Update item price code
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Param item_price_code_id path int true "Item Price Code ID"
// @Param body body masteritempayloads.UpdateItemPriceCode true "Update Item Price Code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code/{item_price_code_id} [put]
func (r *ItemPriceCodeControllerImpl) UpdateItemPriceCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.UpdateItemPriceCode
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	itemPriceCodeId, err := strconv.Atoi(chi.URLParam(request, "item_price_code_id"))
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("invalid item_price_code_id"),
		})
		return
	}

	result, errResult := r.ItemPriceCodeService.UpdateItemPriceCode(itemPriceCodeId, formRequest)
	if errResult != nil {
		helper.ReturnError(writer, request, errResult)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Item Price Code updated successfully!", http.StatusOK)
}

// @Summary Change Status Item Price Code
// @Description Change status item price code
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Param item_price_code_id path int true "Item Price Code ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code/{item_price_code_id} [patch]
func (r *ItemPriceCodeControllerImpl) ChangeStatusItemPriceCode(writer http.ResponseWriter, request *http.Request) {
	itemPriceCodeId, err := strconv.Atoi(chi.URLParam(request, "item_price_code_id"))
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("invalid item_price_code_id"),
		})
		return
	}

	result, errResult := r.ItemPriceCodeService.ChangeStatusItemPriceCode(itemPriceCodeId)
	if errResult != nil {
		helper.ReturnError(writer, request, errResult)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Item Price Code status changed successfully!", http.StatusOK)
}

// @Summary Get Item Price Code Drop Down
// @Description Get item price code drop down
// @Accept json
// @Produce json
// @Tags Master : Item Price Code
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-price-code/drop-down [get]
func (r *ItemPriceCodeControllerImpl) GetItemPriceCodeDropDown(writer http.ResponseWriter, request *http.Request) {

	result, err := r.ItemPriceCodeService.GetItemPriceCodeDropDown()

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)

}
