package masteritemcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"strconv"
	"strings"

	masteritemservice "after-sales/api/services/master/item"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type ItemController interface {
	GetAllItem(writer http.ResponseWriter, request *http.Request)
	GetAllItemLookup(writer http.ResponseWriter, request *http.Request)
	GetItemWithMultiId(writer http.ResponseWriter, request *http.Request)
	GetItemByCode(writer http.ResponseWriter, request *http.Request)
	SaveItem(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItem(writer http.ResponseWriter, request *http.Request)
	GetAllItemDetail(writer http.ResponseWriter, request *http.Request)
	GetItemDetailById(writer http.ResponseWriter, request *http.Request)
	AddItemDetail(writer http.ResponseWriter, request *http.Request)
	DeleteItemDetail(writer http.ResponseWriter, request *http.Request)
}

type ItemControllerImpl struct {
	itemservice masteritemservice.ItemService
}

func NewItemController(ItemService masteritemservice.ItemService) ItemController {
	return &ItemControllerImpl{
		itemservice: ItemService,
	}
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item [get]
func (r *ItemControllerImpl) GetAllItem(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_item.item_code":             queryValues.Get("item_code"),
		"mtr_item.item_name":             queryValues.Get("item_name"),
		"mtr_item.item_type":             queryValues.Get("item_type"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item.is_active":             queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.itemservice.GetAllItem(criteria, paginate)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/lookup [get]
func (r *ItemControllerImpl) GetAllItemLookup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_item.item_group_id": queryValues.Get("item_group_id"),
		"mtr_item.item_class_id": queryValues.Get("item_class_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, totalPages, totalRows, err := r.itemservice.GetAllItemLookup(criteria, paginate)

	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Item With MultiId
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_ids path string true "item_ids"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/multi-id/{item_ids} [get]
func (r *ItemControllerImpl) GetItemWithMultiId(writer http.ResponseWriter, request *http.Request) {

	item_ids := chi.URLParam(request, "item_ids")

	sliceOfString := strings.Split(item_ids, ",")

	result, err := r.itemservice.GetItemWithMultiId(sliceOfString)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/by-code/{item_code} [get]
func (r *ItemControllerImpl) GetItemByCode(writer http.ResponseWriter, request *http.Request) {

	itemCode := chi.URLParam(request, "item_code")

	// Melakukan URL encoding pada item_code
	encodedItemCode := url.PathEscape(itemCode)

	result, err := r.itemservice.GetItemCode(encodedItemCode)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/save [post]
func (r *ItemControllerImpl) SaveItem(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.ItemResponse
	var message = ""

	helper.ReadFromRequestBody(request, &formRequest)

	create, err := r.itemservice.SaveItem(formRequest)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/{item_id} [patch]
func (r *ItemControllerImpl) ChangeStatusItem(writer http.ResponseWriter, request *http.Request) {

	ItemId, _ := strconv.Atoi(chi.URLParam(request, "item_id"))

	response, err := r.itemservice.ChangeStatusItem(int(ItemId))
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/{item_id}/detail [get]
func (r *ItemControllerImpl) GetAllItemDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query() // Retrieve query parameters

	queryParams := map[string]string{
		"mtr_item.item_id":               queryValues.Get("item_id"),
		"mtr_item_detail.item_detail_id": queryValues.Get("item_detail_id"),
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
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/{item_id}/detail/{item_detail_id} [get]
func (r *ItemControllerImpl) GetItemDetailById(writer http.ResponseWriter, request *http.Request) {
	itemID, _ := strconv.Atoi(chi.URLParam(request, "item_id"))
	itemDetailID, _ := strconv.Atoi(chi.URLParam(request, "item_detail_id"))

	result, err := r.itemservice.GetItemDetailById(int(itemID), int(itemDetailID))
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/{item_id}/detail [post]
func (r *ItemControllerImpl) AddItemDetail(writer http.ResponseWriter, request *http.Request) {
	itemID, _ := strconv.Atoi(chi.URLParam(request, "item_id"))

	var itemRequest masteritempayloads.ItemDetailRequest
	helper.ReadFromRequestBody(request, &itemRequest)

	if err := r.itemservice.AddItemDetail(int(itemID), itemRequest); err != nil {
		exceptionsss_test.NewAppException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item/{item_id}/detail/{item_detail_id} [delete]
func (r *ItemControllerImpl) DeleteItemDetail(writer http.ResponseWriter, request *http.Request) {
	itemID, _ := strconv.Atoi(chi.URLParam(request, "item_id"))
	itemDetailID, _ := strconv.Atoi(chi.URLParam(request, "item_detail_id"))

	if err := r.itemservice.DeleteItemDetail(int(itemID), int(itemDetailID)); err != nil {
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Item detail deleted successfully", http.StatusOK)
}
