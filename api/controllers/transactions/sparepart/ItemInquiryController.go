package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"net/http"
)

type ItemInquiryController interface {
	GetAllItemInquiry(writer http.ResponseWriter, request *http.Request)
}

type ItemInquiryControllerImpl struct {
	ItemInquiryService transactionsparepartservice.ItemInquiryService
}

func NewItemInquiryController(itemInquiryService transactionsparepartservice.ItemInquiryService) ItemInquiryController {
	return &ItemInquiryControllerImpl{
		ItemInquiryService: itemInquiryService,
	}
}

func (i *ItemInquiryControllerImpl) GetAllItemInquiry(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"company_id":                queryValues.Get("company_id"),
		"company_session_id":        queryValues.Get("company_session_id"),
		"mi.item_id":                queryValues.Get("item_id"),
		"mtr_item_detail.brand_id":  queryValues.Get("brand_id"),
		"mtr_item_detail.model_id":  queryValues.Get("model_id"),
		"mi.item_code":              queryValues.Get("item_code"),
		"mi.item_name":              queryValues.Get("item_name"),
		"mi.item_class_id":          queryValues.Get("item_class_id"),
		"available_quantity_from":   queryValues.Get("available_quantity_from"),
		"available_quantity_to":     queryValues.Get("available_quantity_to"),
		"sales_price_from":          queryValues.Get("sales_price_from"),
		"sales_price_to":            queryValues.Get("sales_price_to"),
		"mwg.warehouse_group_id":    queryValues.Get("warehouse_group_id"),
		"mwm.warehouse_id":          queryValues.Get("warehouse_id"),
		"mwl.warehouse_location_id": queryValues.Get("warehouse_location_id"),
	}

	if queryParams["company_id"] == "" {
		payloads.NewHandleError(writer, "company_id is required", http.StatusBadRequest)
		return
	}

	if queryParams["company_session_id"] == "" {
		payloads.NewHandleError(writer, "company_session_id is required", http.StatusBadRequest)
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := i.ItemInquiryService.GetAllItemInquiry(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result.Rows), "Get Data Successfully", http.StatusOK, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}
