package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
	"strconv"
	"strings"

	masteritemservice "after-sales/api/services/master/item"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ItemController interface {
	GetAllItem(writer http.ResponseWriter, request *http.Request)
	GetAllItemLookup(writer http.ResponseWriter, request *http.Request)
	GetItemWithMultiId(writer http.ResponseWriter, request *http.Request)
	GetItemByCode(writer http.ResponseWriter, request *http.Request)
	SaveItem(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItem(writer http.ResponseWriter, request *http.Request)
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
// @Param item_code query string false "item_code"
// @Param item_name query string false "item_name"
// @Param item_type query string false "item_type"
// @Param is_active query string false "is_active"
// @Param item_class_code query string false "item_class_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /aftersales-service/api/aftersales/item/ [get]
func (r *ItemControllerImpl) GetAllItem(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_item.item_code":             queryValues.Get("item_code"),
		"mtr_item.item_name":             queryValues.Get("item_name"),
		"mtr_item.item_type":             queryValues.Get("item_type"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item.is_active":             queryValues.Get("is_active"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result := r.itemservice.GetAllItem(criteria)

	payloads.NewHandleSuccess(writer, result, "success", 200)
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
// @Router /aftersales-service/api/aftersales/item/pop-up [get]
func (r *ItemControllerImpl) GetAllItemLookup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"item_code":       queryValues.Get("item_code"),
		"item_name":       queryValues.Get("item_name"),
		"item_type":       queryValues.Get("item_type"),
		"item_group_code": queryValues.Get("item_group_code"),
		"item_class_code": queryValues.Get("item_class_code"),
		"supplier_code":   queryValues.Get("supplier_code"),
		"supplier_name":   queryValues.Get("supplier_name"),
		"is_active":       queryValues.Get("is_active"),
		"sort_of":         queryValues.Get("sort_of"),
		"sort_by":         queryValues.Get("sort_by"),
		"limit":           queryValues.Get("limit"),
		"page":            queryValues.Get("page"),
	}

	result := r.itemservice.GetAllItemLookup(queryParams)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", 200, 0, 0, int64(0), 0)
}

// @Summary Get Item With MultiId
// @Description REST API Item
// @Accept json
// @Produce json
// @Tags Master : Item
// @Param item_ids path string true "item_ids"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /aftersales-service/api/aftersales/item-multi-id/{item_ids} [get]
func (r *ItemControllerImpl) GetItemWithMultiId(writer http.ResponseWriter, request *http.Request) {

	item_ids := chi.URLParam(request, "item_ids")

	sliceOfString := strings.Split(item_ids, ",")

	result := r.itemservice.GetItemWithMultiId(sliceOfString)

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
// @Router /aftersales-service/api/aftersales/item/{item_code} [get]
func (r *ItemControllerImpl) GetItemByCode(writer http.ResponseWriter, request *http.Request) {

	itemCode := chi.URLParam(request, "item_code")

	result := r.itemservice.GetItemCode(itemCode)

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
// @Router /aftersales-service/api/aftersales/item [post]
func (r *ItemControllerImpl) SaveItem(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.ItemResponse
	var message = ""

	helper.ReadFromRequestBody(request, &formRequest)

	create := r.itemservice.SaveItem(formRequest)

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
// @Router /aftersales-service/api/aftersales/item/{item_id} [patch]
func (r *ItemControllerImpl) ChangeStatusItem(writer http.ResponseWriter, request *http.Request) {

	ItemId, _ := strconv.Atoi(chi.URLParam(request, "item_id"))

	response := r.itemservice.ChangeStatusItem(int(ItemId))

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}
