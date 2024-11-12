package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BinningListController interface {
	GetBinningListById(writer http.ResponseWriter, request *http.Request)
	GetAllBinningListWithPagination(writer http.ResponseWriter, request *http.Request)
	InsertBinningListHeader(writer http.ResponseWriter, request *http.Request)
	UpdateBinningListHeader(writer http.ResponseWriter, request *http.Request)
	GetBinningDetailById(writer http.ResponseWriter, request *http.Request)
	GetBinningListDetailWithPagination(writer http.ResponseWriter, request *http.Request)
	InsertBinningListDetail(writer http.ResponseWriter, request *http.Request)
	UpdateBinningListDetail(writer http.ResponseWriter, request *http.Request)
	SubmitBinningList(writer http.ResponseWriter, request *http.Request)
	DeleteBinningList(writer http.ResponseWriter, request *http.Request)
	DeleteBinningListDetailMultiId(writer http.ResponseWriter, request *http.Request)
	GetReferenceNumberTypoPOWithPagination(writer http.ResponseWriter, request *http.Request)
}

type BinningListControllerImpl struct {
	service transactionsparepartservice.BinningListService
}

func (controller *BinningListControllerImpl) GetReferenceNumberTypoPOWithPagination(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"purchase_order_status_id": queryValues.Get("purchase_order_status_id"),
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetReferenceNumberTypoPOWithPagination(filterCondition, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", http.StatusOK, res.Limit, res.Page, res.TotalRows, res.TotalPages)

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
}

// InsertBinningListHeader
//
//	@Summary		Create New Binning List
//	@Description	Create a new Binning List
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Binning List
//	@Param			reqBody					body		transactionsparepartpayloads.BinningListInsertPayloads	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/binning-list/ [post]
func (controller *BinningListControllerImpl) InsertBinningListHeader(writer http.ResponseWriter, request *http.Request) {
	var BinningHeader transactionsparepartpayloads.BinningListInsertPayloads
	helper.ReadFromRequestBody(request, &BinningHeader)
	if validationErr := validation.ValidationForm(writer, request, &BinningHeader); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	res, err := controller.service.InsertBinningListHeader(BinningHeader)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted Binning List Header", http.StatusCreated)
}

// UpdateBinningListHeader
//
//	@Summary		Update Binning List Header
//	@Description	Update Binning List Header
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Binning List
//	@Param			reqBody					body		transactionsparepartpayloads.BinningListSavePayload	true	"Purchase Request Header Data"
//	@Success		201						{object}	payloads.Response
//	@Failure		500,400,401,404,403,422	{object}	exceptions.BaseErrorResponse
//	@Router			/v1/binning-list [patch]
func (controller *BinningListControllerImpl) UpdateBinningListHeader(writer http.ResponseWriter, request *http.Request) {
	var BinningHeaderSavePayloads transactionsparepartpayloads.BinningListSavePayload
	helper.ReadFromRequestBody(request, &BinningHeaderSavePayloads)
	if validationErr := validation.ValidationForm(writer, request, &BinningHeaderSavePayloads); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	res, err := controller.service.UpdateBinningListHeader(BinningHeaderSavePayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted Binning List Header", http.StatusOK)

}

// GetBinningDetailById
//
//	@Summary		Get By Id Binning List Detail
//	@Description	REST API Get By Id Binning List Detail
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Binning List
//	@Param			binning_stock_detail_system_number		path		string	false	"binning_stock_detail_system_number"
//	@Success		200									{object}	transactionsparepartpayloads.BinningListGetByIdResponses
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/binning-list/detail/by-id/{binning_stock_detail_system_number} [get]
func (controller *BinningListControllerImpl) GetBinningDetailById(writer http.ResponseWriter, request *http.Request) {
	BinningStockDetailSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "binning_stock_detail_system_number"))
	result, err := controller.service.GetBinningListDetailById(BinningStockDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Successfully Get Data Detail Binning Stock!", http.StatusOK)

}

// GetBinningListDetailWithPagination
//
//	@Summary		Get All Binning List Detail By Id
//	@Description	REST API Get All Binning List Detail By Id
//	@Accept			json
//	@Produce		json
//	@Tags			Transaction : Binning List
//	@Param			page								query		string	true	"page"
//	@Param			limit								query		string	true	"limit"
//	@Param			binning_system_number				path		int	true	"binning_system_number"
//	@Param			sort_by								query		string	false	"sort_by"
//	@Param			sort_of								query		string	false	"sort_of"
//	@Success		200									{object}	[]transactionsparepartpayloads.BinningListGetByIdResponses
//	@Failure		500,400,401,404,403,422				{object}	exceptions.BaseErrorResponse
//	@Router			/v1/binning-list/detail/{binning_system_number} [get]
func (controller *BinningListControllerImpl) GetBinningListDetailWithPagination(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	BinningStockSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "binning_system_number"))
	queryParams := map[string]string{
		//"A.binning_system_number": queryValues.Get("binning_system_number"),
	}
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filter := utils.BuildFilterCondition(queryParams)
	res, err := controller.service.GetAllBinningListDetailWithPagination(filter, paginations, BinningStockSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Get Binning Detail Scucess", http.StatusOK, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// InsertBinningListDetail route binning-list/detail [post]
func (controller *BinningListControllerImpl) InsertBinningListDetail(writer http.ResponseWriter, request *http.Request) {
	var BinningListSavePayloads transactionsparepartpayloads.BinningListDetailPayloads
	helper.ReadFromRequestBody(request, &BinningListSavePayloads)
	if validationErr := validation.ValidationForm(writer, request, &BinningListSavePayloads); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	res, err := controller.service.InsertBinningListDetail(BinningListSavePayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Inserted Binning List Detail", http.StatusCreated)
}

// route binning-list/detail [patch]
func (controller *BinningListControllerImpl) UpdateBinningListDetail(writer http.ResponseWriter, request *http.Request) {
	var BinningListSavePayloads transactionsparepartpayloads.BinningListDetailUpdatePayloads
	helper.ReadFromRequestBody(request, &BinningListSavePayloads)
	if validationErr := validation.ValidationForm(writer, request, &BinningListSavePayloads); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	res, err := controller.service.UpdateBinningListDetail(BinningListSavePayloads)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Updated Binning List Detail", http.StatusOK)
}

// router binning-list/submit/{binning_system_number}
func (controller *BinningListControllerImpl) SubmitBinningList(writer http.ResponseWriter, request *http.Request) {
	BinningStockSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "binning_system_number"))
	res, err := controller.service.SubmitBinningList(BinningStockSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Submit Binning List", http.StatusOK)

}
func (controller *BinningListControllerImpl) DeleteBinningList(writer http.ResponseWriter, request *http.Request) {
	BinningStockSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "binning_system_number"))
	res, err := controller.service.DeleteBinningList(BinningStockSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Delete Binning List", http.StatusOK)
}

func (controller *BinningListControllerImpl) DeleteBinningListDetailMultiId(writer http.ResponseWriter, request *http.Request) {
	BinningStockDetailSystemNumber := chi.URLParam(request, "binning_detail_multi_id")
	res, err := controller.service.DeleteBinningListDetailMultiId(BinningStockDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Successfully Delete Binning List Detail Multi Id", http.StatusOK)
}
