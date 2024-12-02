package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type ItemInquiryController interface {
	GetAllItemInquiry(writer http.ResponseWriter, request *http.Request)
	GetByIdItemInquiry(writer http.ResponseWriter, request *http.Request)
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

func (i *ItemInquiryControllerImpl) GetByIdItemInquiry(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	itemId, errA := strconv.Atoi(queryValues.Get("item_id"))
	if errA != nil || itemId <= 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param 'item_id', please check your param input")})
		return
	}

	companyId, errB := strconv.Atoi(queryValues.Get("company_id"))
	if errB != nil || companyId <= 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param 'company_id', please check your param input")})
		return
	}

	var warehouseId int
	if queryValues.Get("warehouse_id") != "" {
		id, errC := strconv.Atoi(queryValues.Get("warehouse_id"))
		warehouseId = id
		if errC != nil || warehouseId < 0 {
			exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param 'warehouse_id', please check your param input")})
			return
		}
	}

	var warehouseLocationId int
	if queryValues.Get("warehouse_location_id") != "" {
		id, errD := strconv.Atoi(queryValues.Get("warehouse_location_id"))
		warehouseLocationId = id
		if errD != nil || warehouseLocationId < 0 {
			exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param 'warehouse_location_id', please check your param input")})
			return
		}
	}

	fmt.Println(warehouseId, warehouseLocationId)

	brandId, errE := strconv.Atoi(queryValues.Get("brand_id"))
	if errE != nil || brandId <= 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param 'brand_id', please check your param input")})
		return
	}

	currencyId, errF := strconv.Atoi(queryValues.Get("currency_id"))
	if errF != nil || currencyId <= 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param 'currency_id', please check your param input")})
		return
	}

	companySessionId, errG := strconv.Atoi(queryValues.Get("company_session_id"))
	if errG != nil || companySessionId <= 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param 'company_session_id', please check your param input")})
		return
	}

	queryParam := transactionsparepartpayloads.ItemInquiryGetByIdFilter{
		ItemId:              itemId,
		CompanyId:           companyId,
		WarehouseId:         warehouseId,
		WarehouseLocationId: warehouseLocationId,
		BrandId:             brandId,
		CurrencyId:          currencyId,
		CompanySessionId:    companySessionId,
	}

	result, err := i.ItemInquiryService.GetByIdItemInquiry(queryParam)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
