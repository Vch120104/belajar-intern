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
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ItemController interface {
	GetAllItemLookup(writer http.ResponseWriter, request *http.Request)
	GetAllItemInventory(writer http.ResponseWriter, request *http.Request)
	GetItemInventoryByCode(writer http.ResponseWriter, request *http.Request)
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
	DeleteItemDetails(writer http.ResponseWriter, request *http.Request)
	UpdateItem(writer http.ResponseWriter, request *http.Request)
	UpdateItemDetail(writer http.ResponseWriter, request *http.Request)
	GetPrincipalBrandParent(writer http.ResponseWriter, request *http.Request)
	GetPrincipalBrandDropdown(writer http.ResponseWriter, request *http.Request)
	AddItemDetailByBrand(writer http.ResponseWriter, request *http.Request)
	GetAllItemSearch(writer http.ResponseWriter, request *http.Request)
	GetPrincipalCatalog(writer http.ResponseWriter, request *http.Request)
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
		"mtr_item_type.item_type_code":   queryValues.Get("item_type"),
		"mtr_item.item_class_id":         queryValues.Get("item_class_id"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item.is_active":             queryValues.Get("is_active"),
		"mtr_item.item_group_id":         queryValues.Get("item_group_id"),
		"mtr_item_group.item_group_code": queryValues.Get("item_group_code"),
		"dms_microservices_general_dev.dbo.mtr_supplier.supplier_code": queryValues.Get("supplier_code"),
		"dms_microservices_general_dev.dbo.mtr_supplier.supplier_name": queryValues.Get("supplier_name"),
	}

	itemTypes := strings.Split(queryValues.Get("item_type"), ",")
	var processedItemTypes []string
	for _, itemType := range itemTypes {
		switch strings.ToLower(itemType) {
		case "goods", "g", "go", "goo", "good":
			processedItemTypes = append(processedItemTypes, "G")
		case "services", "s", "se", "ser", "serv", "servi", "servic", "service":
			processedItemTypes = append(processedItemTypes, "S")
		}
	}

	if len(processedItemTypes) > 0 {
		queryParams["mtr_item_type.item_type_code"] = strings.Join(processedItemTypes, ",")
	}

	itemIDs := strings.Split(queryValues.Get("item_id"), ",")
	supplierIDs := strings.Split(queryValues.Get("supplier_id"), ",")

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	data, err := r.itemservice.GetAllItemSearch(criteria, itemIDs, supplierIDs, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, data.Rows, "success", http.StatusOK, data.Limit, data.Page, data.TotalRows, data.TotalPages)
}

func (r *ItemControllerImpl) GetAllItemInventory(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"itm.item_code":       queryValues.Get("item_code"),
		"itm.item_name":       queryValues.Get("item_name"),
		"cls.item_class_name": queryValues.Get("item_class_name"),
		"grp.item_group_code": queryValues.Get("item_group_code"),
		"cls.item_class_code": queryValues.Get("item_class_code"),
		"uom.uom_code":        queryValues.Get("uom_code"),
		"itm.is_active":       queryValues.Get("is_active"),
		"itm.item_class_id":   queryValues.Get("item_class_id"), // Use case: item class dropdown
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	data, err := r.itemservice.GetAllItemInventory(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, data.Rows, "success", http.StatusOK, data.Limit, data.Page, data.TotalRows, data.TotalPages)
}

func (r *ItemControllerImpl) GetItemInventoryByCode(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	itemCode := queryValues.Get("item_code")

	result, err := r.itemservice.GetItemInventoryByCode(itemCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "success", 200)
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

	queryValues := request.URL.Query()
	itemCodeEncode := queryValues.Get("item_code")

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

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

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
		"item_id":             queryValues.Get("item_id"),
		"is_active":           queryValues.Get("is_active"),
		"brand_id":            queryValues.Get("brand_id"),
		"brand_name":          queryValues.Get("brand_name"),
		"model_id":            queryValues.Get("model_id"),
		"model_code":          queryValues.Get("model_code"),
		"model_description":   queryValues.Get("model_description"),
		"variant_id":          queryValues.Get("variant_id"),
		"variant_code":        queryValues.Get("variant_code"),
		"variant_description": queryValues.Get("variant_description"),
		"mileage_every":       queryValues.Get("mileage_every"),
		"return_every":        queryValues.Get("return_every"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: chi.URLParam(request, "sort_of"),
		SortBy: chi.URLParam(request, "sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)
	result, err := r.itemservice.GetAllItemDetail(criteria, paginate)

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

	itemIdStr := chi.URLParam(request, "item_id")
	itemId, errA := strconv.Atoi(itemIdStr)
	if errA != nil {
		payloads.NewHandleError(writer, "Failed to read request param", http.StatusBadRequest)
		return
	}

	var itemRequest masteritempayloads.ItemDetailRequest
	helper.ReadFromRequestBody(request, &itemRequest)

	itemDetail, err := r.itemservice.AddItemDetail(itemId, itemRequest)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, itemDetail, "Item detail added successfully", http.StatusCreated)
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
func (r *ItemControllerImpl) DeleteItemDetails(writer http.ResponseWriter, request *http.Request) {
	itemID, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to read request param, please check your param input",
			Err:        errA,
		})
		return
	}

	// Parse multiple itemDetailIDs from the request body
	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var itemDetailIDs []int
	for _, element := range elements {
		num, err := strconv.Atoi(strings.TrimSpace(element))
		if err != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		itemDetailIDs = append(itemDetailIDs, num)
	}

	resp, err := r.itemservice.DeleteItemDetails(itemID, itemDetailIDs)
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, resp, "Item details deleted successfully", http.StatusOK)
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

	itemIdStr := chi.URLParam(request, "item_id")
	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		payloads.NewHandleError(writer, "Failed to read request param", http.StatusBadRequest)
		return
	}

	itemDetailStr := chi.URLParam(request, "item_detail_id")
	itemDetailId, err := strconv.Atoi(itemDetailStr)
	if err != nil {
		payloads.NewHandleError(writer, "Failed to read request param", http.StatusBadRequest)
		return
	}

	var itemDetailRequest masteritempayloads.ItemDetailUpdateRequest
	helper.ReadFromRequestBody(request, &itemDetailRequest)

	update, errResp := r.itemservice.UpdateItemDetail(itemId, itemDetailId, itemDetailRequest)
	if errResp != nil {
		exceptions.NewAppException(writer, request, errResp)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Item detail updated successfully", http.StatusOK)
}

func (r *ItemControllerImpl) GetPrincipalBrandDropdown(writer http.ResponseWriter, request *http.Request) {
	result, err := r.itemservice.GetPrincipalBrandDropdown()
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

func (r *ItemControllerImpl) GetPrincipalBrandParent(writer http.ResponseWriter, request *http.Request) {
	principalCatalogId, errA := strconv.Atoi(chi.URLParam(request, "principal_catalog_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	result, err := r.itemservice.GetPrincipalBrandParent(principalCatalogId)
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

func (r *ItemControllerImpl) GetPrincipalCatalog(writer http.ResponseWriter, request *http.Request) {
	result, err := r.itemservice.GetPrincipalCatalog()
	if err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}
