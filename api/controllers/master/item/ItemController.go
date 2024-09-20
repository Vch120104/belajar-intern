package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ItemController interface {
	GetAllItem(writer http.ResponseWriter, request *http.Request)
	GetAllItemLookup(writer http.ResponseWriter, request *http.Request)
	GetItemWithMultiId(writer http.ResponseWriter, request *http.Request)
	GetItembyId(writer http.ResponseWriter, request *http.Request)
	GetItemByCode(writer http.ResponseWriter, request *http.Request)
	SaveItem(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItem(writer http.ResponseWriter, request *http.Request)
	GetUomTypeDropDown(writer http.ResponseWriter, request *http.Request)
	GetUomDropDown(writer http.ResponseWriter, request *http.Request)
	GetAllItemDetail(writer http.ResponseWriter, request *http.Request)
	GetItemDetailById(writer http.ResponseWriter, request *http.Request)
	AddItemDetail(writer http.ResponseWriter, request *http.Request)
	DeleteItemDetail(writer http.ResponseWriter, request *http.Request)
	UpdateItem(writer http.ResponseWriter, request *http.Request)
	UpdateItemDetail(writer http.ResponseWriter, request *http.Request)
	GetPrincipleBrandParent(writer http.ResponseWriter, request *http.Request)
	GetPrincipleBrandDropdown(writer http.ResponseWriter, request *http.Request)
	AddItemDetailByBrand(writer http.ResponseWriter, request *http.Request)
	GetAllItemSearch(writer http.ResponseWriter, request *http.Request)
	GetCatalogCode(writer http.ResponseWriter, request *http.Request)
	GetAllItemListTransLookup(writer http.ResponseWriter, request *http.Request)
}

type ItemControllerImpl struct {
	itemservice masteritemservice.ItemService
}

func NewItemController(ItemService masteritemservice.ItemService) ItemController {
	return &ItemControllerImpl{
		itemservice: ItemService,
	}
}

// GetAllItemSearch
func (r *ItemControllerImpl) GetAllItemSearch(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_item.item_code":             queryValues.Get("item_code"),
		"mtr_item.item_name":             queryValues.Get("item_name"),
		"mtr_item.item_type":             queryValues.Get("item_type"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item.is_active":             queryValues.Get("is_active"),
		"mtr_item_group.item_group_code": queryValues.Get("item_group_code"),
	}

	// Handle multi_id and supplier_id as multiple parameters
	itemIDs := strings.Split(queryValues.Get("item_id"), ",")
	supplierIDs := strings.Split(queryValues.Get("supplier_id"), ",")

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	data, totalPages, totalRows, err := r.itemservice.GetAllItemSearch(criteria, itemIDs, supplierIDs, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(data), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// GetItembyId implements ItemController.
func (r *ItemControllerImpl) GetItembyId(writer http.ResponseWriter, request *http.Request) {
	itemId, errA := strconv.Atoi(chi.URLParam(request, "item_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.itemservice.GetItemById(itemId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// GetUomDropDown implements ItemController.
func (r *ItemControllerImpl) GetUomDropDown(writer http.ResponseWriter, request *http.Request) {

	uomTypeId, errA := strconv.Atoi(chi.URLParam(request, "uom_type_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.itemservice.GetUomDropDown(uomTypeId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// GetUomTypeDropDown implements ItemController.
func (r *ItemControllerImpl) GetUomTypeDropDown(writer http.ResponseWriter, request *http.Request) {

	result, err := r.itemservice.GetUomTypeDropDown()

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)

}

// @Summary Get All Item
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_code query string false "item_code"
// @Param item_name query string false "item_name"
// @Param item_type query string false "item_type"
// @Param is_active query string false "is_active"
// @Param item_class_code query string false "item_class_code"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item [get]
func (r *ItemControllerImpl) GetAllItem(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_item.item_id":               queryValues.Get("item_id"),
		"mtr_item.item_code":             queryValues.Get("item_code"),
		"mtr_item.item_name":             queryValues.Get("item_name"),
		"mtr_item.item_type":             queryValues.Get("item_type"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item.is_active":             queryValues.Get("is_active"),
		"mtr_item_group.item_group_code": queryValues.Get("item_group_code"),
		"mtr_supplier.supplier_code":     queryValues.Get("supplier_code"),
		"mtr_supplier.supplier_name":     queryValues.Get("supplier_name"),
		"mtr_item.supplier_id":           queryValues.Get("supplier_id"), // Add supplier_id to queryParams
	}

	// Periksa apakah parameter query ada dan tidak kosong
	for key, value := range queryParams {
		if value == "" {
			delete(queryParams, key)
		}
	}

	// Debug log for query parameters
	fmt.Printf("Query parameters: %+v\n", queryParams)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	data, totalPages, totalRows, err := r.itemservice.GetAllItem(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(data), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *ItemControllerImpl) GetAllItemListTransLookup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"item_code":     queryValues.Get("item_code"),
		"item_name":     queryValues.Get("item_name"),
		"item_class_id": queryValues.Get("item_class_id"),
		"item_type":     queryValues.Get("item_type"),
		"item_level_1":  queryValues.Get("item_level_1"),
		"item_level_2":  queryValues.Get("item_level_2"),
		"item_level_3":  queryValues.Get("item_level_3"),
		"item_level_4":  queryValues.Get("item_level_4"),
	}

	for key, value := range queryParams {
		if value == "" {
			delete(queryParams, key)
		}
	}

	fmt.Printf("Query parameters: %+v\n", queryParams)

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	data, err := r.itemservice.GetAllItemListTransLookup(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(data.Rows), "success", 200, paginate.Limit, paginate.Page, data.TotalRows, data.TotalPages)
}

// @Summary Get All Item Lookup
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_code query string false "item_code"
// @Param item_name query string false "item_name"
// @Param item_type query string false "item_type"
// @Param item_group_code query string false "item_group_code"
// @Param item_class_code query string false "item_class_code"
// @Param supplier_code query string false "supplier_code"
// @Param supplier_name query string false "supplier_name"
// @Param is_active query string false "is_active"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/lookup [get]
func (r *ItemControllerImpl) GetAllItemLookup(writer http.ResponseWriter, request *http.Request) {
	// queryValues := request.URL.Query()

	// internalFilterCondition := map[string]string{
	// 	"item_code":       queryValues.Get("item_code"),
	// 	"item_name":       queryValues.Get("item_name"),
	// 	"item_type":       queryValues.Get("item_type"),
	// 	"item_group_code": queryValues.Get("item_group_code"),
	// 	"item_class_code": queryValues.Get("item_class_code"),
	// 	"is_active":       queryValues.Get("is_active"),
	// }
	// externalFilterCondition := map[string]string{

	// 	"supplier_code": queryValues.Get("supplier_code"),
	// 	"supplier_name": queryValues.Get("supplier_name"),
	// }

	// paginate := pagination.Pagination{
	// 	Limit:  utils.NewGetQueryInt(queryValues, "limit"),
	// 	Page:   utils.NewGetQueryInt(queryValues, "page"),
	// 	SortOf: queryValues.Get("sort_of"),
	// 	SortBy: queryValues.Get("sort_by"),
	// }

	// internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	// externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	// result, totalPages, totalRows, err := r.itemservice.GetAllItemLookup(internalCriteria, externalCriteria, paginate)

	// if err != nil {
	// 	exceptions.NewNotFoundException(writer, request, err)
	// 	return
	// }
	// payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)

	panic("ON PROGRESSS")
}

// @Summary Get Item With MultiId
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_ids path string true "item_ids"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/multi-id/{item_ids} [get]
func (r *ItemControllerImpl) GetItemWithMultiId(writer http.ResponseWriter, request *http.Request) {

	item_ids := chi.URLParam(request, "item_ids")

	sliceOfString := strings.Split(item_ids, ",")

	result, err := r.itemservice.GetItemWithMultiId(sliceOfString)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Get Item By code
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_code path string true "item_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/by-code/{item_code} [get]
func (r *ItemControllerImpl) GetItemByCode(writer http.ResponseWriter, request *http.Request) {

	itemCode := chi.URLParam(request, "item_code")

	itemCodeEncode := strings.ReplaceAll(itemCode, "!", "/")

	// Melakukan URL encoding pada item_code
	// encodedItemCode := url.PathEscape(itemCode)

	result, err := r.itemservice.GetItemCode(itemCodeEncode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @param reqBody body masteritempayloads.ItemRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/save [post]
func (r *ItemControllerImpl) SaveItem(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.ItemRequest
	var message = ""

	helper.ReadFromRequestBody(request, &formRequest)

	create, err := r.itemservice.SaveItem(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	if formRequest.ItemId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Item
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @param item_id path int true "item_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/{item_id} [patch]
func (r *ItemControllerImpl) ChangeStatusItem(writer http.ResponseWriter, request *http.Request) {

	ItemId, errA := strconv.Atoi(chi.URLParam(request, "item_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.itemservice.ChangeStatusItem(int(ItemId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

// @Summary Get All Detail Items
// @Description Retrieve all detail items from an item by its ID
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_id path int true "Item ID"
// @Param page query int true "Page number"
// @Param limit query int true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/{item_id}/detail [get]
func (r *ItemControllerImpl) GetAllItemDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"item_id": queryValues.Get("item_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	data, totalPages, totalRows, err := r.itemservice.GetAllItemDetail(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(data), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Detail Item By Id
// @Description Retrieve a detail item from an item by its ID
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_id path int true "Item ID"
// @Param item_detail_id path int true "Item Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/detail/{item_id}/{item_detail_id} [get]
func (r *ItemControllerImpl) GetItemDetailById(writer http.ResponseWriter, request *http.Request) {
	itemID, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	itemDetailID, errA := strconv.Atoi(chi.URLParam(request, "item_detail_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.itemservice.GetItemDetailById(int(itemID), int(itemDetailID))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Add Item Detail
// @Description Add a new item detail to an item by its ID
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_id path int true "Item ID"
// @Param reqBody body masteritempayloads.ItemDetailRequest true "Item Detail Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/{item_id}/detail [post]
func (r *ItemControllerImpl) AddItemDetail(writer http.ResponseWriter, request *http.Request) {
	itemID, errA := strconv.Atoi(chi.URLParam(request, "item_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var itemRequest masteritempayloads.ItemDetailRequest
	helper.ReadFromRequestBody(request, &itemRequest)

	if err := r.itemservice.AddItemDetail(int(itemID), itemRequest); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Item detail added successfully", http.StatusOK)
}

// @Summary Delete Item Detail
// @Description Delete an item detail from an item by its ID
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_id path int true "Item ID"
// @Param item_detail_id path int true "Item Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item/{item_id}/detail/{item_detail_id} [delete]
func (r *ItemControllerImpl) DeleteItemDetail(writer http.ResponseWriter, request *http.Request) {
	itemID, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	itemDetailID, errA := strconv.Atoi(chi.URLParam(request, "item_detail_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	if err := r.itemservice.DeleteItemDetail(int(itemID), int(itemDetailID)); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Item detail deleted successfully", http.StatusOK)
}

func (r *ItemControllerImpl) UpdateItem(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemUpdateRequest

	helper.ReadFromRequestBody(request, &formRequest)
	item_id, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	_, err := r.itemservice.UpdateItem(item_id, formRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, nil, "Item updated successfully", http.StatusOK)
}

func (r *ItemControllerImpl) UpdateItemDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemDetailUpdateRequest

	helper.ReadFromRequestBody(request, &formRequest)
	item_detail_id, errA := strconv.Atoi(chi.URLParam(request, "item_detail_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	_, err := r.itemservice.UpdateItemDetail(item_detail_id, formRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, nil, "Item updated successfully", http.StatusOK)
}

func (r *ItemControllerImpl) GetPrincipleBrandDropdown(writer http.ResponseWriter, request *http.Request) {
	result, err := r.itemservice.GetPrincipleBrandDropdown()
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

func (r *ItemControllerImpl) GetPrincipleBrandParent(writer http.ResponseWriter, request *http.Request) {
	principleBrandCode := chi.URLParam(request, "catalogue_code")
	result, err := r.itemservice.GetPrincipleBrandParent(principleBrandCode)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

func (r *ItemControllerImpl) AddItemDetailByBrand(writer http.ResponseWriter, request *http.Request) {
	ItemId, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	Id := chi.URLParam(request, "brand_id")
	result, err := r.itemservice.AddItemDetailByBrand(Id, ItemId)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

func (r *ItemControllerImpl) GetCatalogCode(writer http.ResponseWriter, request *http.Request) {
	result, err := r.itemservice.GetCatalogCode()
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}
