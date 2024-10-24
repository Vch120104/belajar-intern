package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
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

	results, totalPages, totalRows, err := r.ItemPriceCodeService.GetAllItemPriceCode(filter, pagination)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, results, "Get Data Successfully!", http.StatusOK, pagination.Limit, pagination.Page, int64(totalRows), totalPages)
}

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

func (r *ItemPriceCodeControllerImpl) GetItemPriceCodeByCode(writer http.ResponseWriter, request *http.Request) {
	itemPriceCode := chi.URLParam(request, "item_price_code")

	result, err := r.ItemPriceCodeService.GetByCodeItemPriceCode(itemPriceCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemPriceCodeControllerImpl) SaveItemPriceCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.SaveItemPriceCode
	helper.ReadFromRequestBody(request, &formRequest)

	result, err := r.ItemPriceCodeService.SaveItemPriceCode(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message := "Item Price Code saved successfully!"
	payloads.NewHandleSuccess(writer, result, message, http.StatusOK)
}

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

func (r *ItemPriceCodeControllerImpl) UpdateItemPriceCode(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.UpdateItemPriceCode
	helper.ReadFromRequestBody(request, &formRequest)

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

func (r *ItemPriceCodeControllerImpl) GetItemPriceCodeDropDown(writer http.ResponseWriter, request *http.Request) {

	result, err := r.ItemPriceCodeService.GetItemPriceCodeDropDown()

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)

}
