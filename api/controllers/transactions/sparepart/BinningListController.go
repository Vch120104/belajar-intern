package transactionsparepartcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type BinningListController interface {
	GetBinningListById(writer http.ResponseWriter, request *http.Request)
	GetAllBinningListWithPagination(writer http.ResponseWriter, request *http.Request)
}
type BinningListControllerImpl struct {
	service transactionsparepartservice.BinningListService
}

func NewBinningListControllerImpl(service transactionsparepartservice.BinningListService) BinningListController {
	return &BinningListControllerImpl{service: service}
}

// GetBinningListById
//
//	@Summary		Get By Id Binning List
//	@Description	REST API Get By Id Binning List
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Binning List
//	@Param			binning_stock_system_number		path		string	false	"binning_stock_system_number"
//	@Success		200									{object}	transactionsparepartpayloads.BinningListGetByIdResponse
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/binning-list/by-id/{binning_stock_system_number} [get]
func (controller *BinningListControllerImpl) GetBinningListById(writer http.ResponseWriter, request *http.Request) {
	BinningStockSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "binning_stock_system_number"))
	result, err := controller.service.GetBinningListById(BinningStockSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Successfully Get Data Binning Stock !", http.StatusOK)
}

// GetAllBinningListWithPagination
//
//	@Summary		Get All Binning List
//	@Description	REST API Get All Binning List
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Binning List
//	@Param			page								query		string	true	"page"
//	@Param			limit								query		string	true	"limit"
//	@Param			binning_document_number				query		string	false	"binning_document_number"
//	@Param			reference_document_number			query		string	false	"reference_document_number"
//	@Param			supplier_invoice_number				query		string	false	"supplier_invoice_number"
//	@Param			warehouse_group_id					query		string	false	"warehouse_group_id"
//	@Param			warehouse_id						query		string	false	"warehouse_id"
//	@Param			supplier_case_number				query		string	false	"supplier_case_number"
//	@Param			binning_document_status_id			query		string	false	"binning_document_status_id"
//	@Param			supplier_id							query		string	false	"supplier_id"
//	@Param			supplier_delivery_order_number		query		string	false	"supplier_delivery_order_number"
//	@Param			sort_by								query		string	false	"sort_by"
//	@Param			sort_of								query		string	false	"sort_of"
//	@Success		200									{object}	[]transactionsparepartpayloads.PurchaseRequestGetAllListResponses
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/binning-list/ [get]
func (controller *BinningListControllerImpl) GetAllBinningListWithPagination(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"binning_document_number":        queryValues.Get("binning_document_number"),
		"reference_document_number":      queryValues.Get("reference_document_number"),
		"supplier_id":                    queryValues.Get("supplier_id"),
		"warehouse_group_id":             queryValues.Get("warehouse_group_id"),
		"warehouse_id":                   queryValues.Get("warehouse_id"),
		"supplier_case_number":           queryValues.Get("supplier_case_number"),
		"binning_document_status_id":     queryValues.Get("binning_document_status_id"),
		"supplier_invoice_number":        queryValues.Get("supplier_invoice_number"),
		"supplier_delivery_order_number": queryValues.Get("supplier_delivery_order_number"),
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllBinningListWithPagination(filterCondition, paginations, request.Context())
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", http.StatusOK, res.Limit, res.Page, res.TotalRows, res.TotalPages)
	return
}
