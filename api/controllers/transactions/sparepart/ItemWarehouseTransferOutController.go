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
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ItemWarehouseTransferOutController interface {
	InsertHeader(writer http.ResponseWriter, request *http.Request)
	InsertDetail(writer http.ResponseWriter, request *http.Request)
	InsertDetailFromReceipt(writer http.ResponseWriter, request *http.Request)
	GetTransferOutById(writer http.ResponseWriter, request *http.Request)
	GetAllTransferOut(writer http.ResponseWriter, request *http.Request)
	GetAllTransferOutDetail(writer http.ResponseWriter, request *http.Request)
	SubmitTransferOut(writer http.ResponseWriter, request *http.Request)
	UpdateTransferOutDetail(writer http.ResponseWriter, request *http.Request)
	DeleteTransferOutDetail(writer http.ResponseWriter, request *http.Request)
	DeleteTransferOut(writer http.ResponseWriter, request *http.Request)
}

func NewItemWarehouseTransferOutControllerImpl(itemWarehouseTransferOutService transactionsparepartservice.ItemWarehouseTransferOutService) ItemWarehouseTransferOutController {
	return &ItemWarehouseTransferOutControllerImpl{
		ItemWarehouseTransferOutService: itemWarehouseTransferOutService,
	}
}

type ItemWarehouseTransferOutControllerImpl struct {
	ItemWarehouseTransferOutService transactionsparepartservice.ItemWarehouseTransferOutService
}

// DeleteTransferOut implements ItemWarehouseTransferOutController.
// @Summary Delete Header Transfer Out
// @Description Delete Header Transfer Out
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param id path int true "Transfer Out System Number"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/{id} [delete]
func (r *ItemWarehouseTransferOutControllerImpl) DeleteTransferOut(writer http.ResponseWriter, request *http.Request) {
	transferOutSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	success, err := r.ItemWarehouseTransferOutService.DeleteTransferOut(transferOutSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "delete success", http.StatusCreated)
}

// DeleteTransferOutDetail implements ItemWarehouseTransferOutController.
// @Summary Delete Detail Transfer Out
// @Description Delete Detail Transfer Out
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param id path string true "Detail Multi ID"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/detail/{id} [delete]
func (r *ItemWarehouseTransferOutControllerImpl) DeleteTransferOutDetail(writer http.ResponseWriter, request *http.Request) {
	multiId := chi.URLParam(request, "id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var intIds []int
	for _, element := range elements {
		num, err := strconv.Atoi(strings.TrimSpace(element))
		if err != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		intIds = append(intIds, num)
	}
	success, err := r.ItemWarehouseTransferOutService.DeleteTransferOutDetail(intIds)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "delete success", http.StatusCreated)
}

// GetAllTransferOut implements ItemWarehouseTransferOutController.
// @Summary Get All Warehouse Transfer Out
// @Description Get All Warehouse Transfer Out
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param transfer_out_status_id query int false "Transfer Out Status ID"
// @Param transfer_request_document_number query string false "Transfer Request Document Number"
// @Param transfer_out_document_number query string false "Transfer Out Document Number"
// @Param transfer_out_warehouse_group_id query int false "Transfer out Warehouse Group ID"
// @Param company_id query int false "Company ID"
// @Param transfer_out_date_from query string false "Transfer Out Date From"
// @Param transfer_out_date_to query string false "Transfer Out Date To"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out [get]
func (r *ItemWarehouseTransferOutControllerImpl) GetAllTransferOut(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"transfer_out_status_id":                     queryValues.Get("transfer_out_status_id"),
		"transfer_request_document_number":           queryValues.Get("transfer_request_document_number"),
		"transfer_out_document_number":               queryValues.Get("transfer_out_document_number"),
		"wmt.warehouse_group_id":                     queryValues.Get("transfer_out_warehouse_group_id"),
		"trx_item_warehouse_transfer_out.company_id": queryValues.Get("company_id"),
	}

	dateParams := map[string]string{
		"transfer_request_date_from": queryValues.Get("transfer_request_date_from"),
		"transfer_request_date_to":   queryValues.Get("transfer_request_date_to"),
	}

	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := r.ItemWarehouseTransferOutService.GetAllTransferOut(filterCondition, dateParams, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// GetAllTransferOutDetail implements ItemWarehouseTransferOutController.
// @Summary Get All Detail Transfer Out
// @Description Get All Detail Transfer Out
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param transfer_out_system_number query int true "Transfer out System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/detail [get]
func (r *ItemWarehouseTransferOutControllerImpl) GetAllTransferOutDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	transferOutNumber := utils.NewGetQueryInt(queryValues, "transfer_out_system_number")

	res, err := r.ItemWarehouseTransferOutService.GetAllTransferOutDetail(transferOutNumber, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// GetTransferOutById implements ItemWarehouseTransferOutController
// @Description Get By ID Transfer Out
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param id path int true "Transfer Out System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/{id} [get].
func (r *ItemWarehouseTransferOutControllerImpl) GetTransferOutById(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	success, err := r.ItemWarehouseTransferOutService.GetTransferOutById(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// InsertDetail implements ItemWarehouseTransferOutController.
// @Summary Insert Item Warehouse Transfer Out Detail
// @Description Insert Item Warehouse Transfer Out Detail
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param InsertItemWarehouseTransferOutDetail body transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailRequest true "Insert Item Warehouse Transfer Out Detail"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/detail [post]
func (r *ItemWarehouseTransferOutControllerImpl) InsertDetail(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailRequest
	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferOutService.InsertDetail(transferRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Insert Data Success", http.StatusCreated)
}

// InsertDetailFromReceipt implements ItemWarehouseTransferOutController.
// @Summary Insert Item Warehouse Transfer Out Detail From Receipt
// @Description Insert Item Warehouse Transfer Out Detail From Receipt
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param InsertItemWarehouseTransferOut body transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailCopyReceiptRequest true "Insert Item Warehouse Transfer Out Detail From Receipt"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/detail/copy-receipt [post]
func (r *ItemWarehouseTransferOutControllerImpl) InsertDetailFromReceipt(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailCopyReceiptRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferOutService.InsertDetailFromReceipt(transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailCopyReceiptRequest{
		TransferOutSystemNumber: transferRequest.TransferOutSystemNumber,
	})

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Insert Data Success", http.StatusCreated)
}

// InsertHeader implements ItemWarehouseTransferOutController.
// @Summary Insert Item Warehouse Transfer Out Header
// @Description Insert Item Warehouse Transfer Out Header
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param InsertItemWarehouseTransferOut body transactionsparepartpayloads.InsertItemWarehouseHeaderTransferOutRequest true "Insert Item Warehouse Transfer Out"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out [post]
func (r *ItemWarehouseTransferOutControllerImpl) InsertHeader(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.InsertItemWarehouseHeaderTransferOutRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferOutService.InsertHeader(transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusCreated)
}

// SubmitTransferOut implements ItemWarehouseTransferOutController.
// @Description Submit Item Warehouse Transfer Out
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param id path int true "Transfer Out System Number"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/submit/{id} [put]
func (r *ItemWarehouseTransferOutControllerImpl) SubmitTransferOut(writer http.ResponseWriter, request *http.Request) {
	transferOutSystemNumber, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	success, err := r.ItemWarehouseTransferOutService.SubmitTransferOut(transferOutSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Submit Data Success", http.StatusCreated)
}

// UpdateTransferOutDetail implements ItemWarehouseTransferOutController.
// @Summary Update Item Warehouse Transfer Out
// @Description Update Item Warehouse Transfer Out
// @Tags Transaction Sparepart: Item Warehouse Transfer Out
// @Accept json
// @Produce json
// @Param id path int true "Transfer Out System Number"
// @Param UpdateItemWarehouseTransferOut body transactionsparepartpayloads.UpdateItemWarehouseTransferOutDetailRequest true "Update Item Warehouse Transfer Out"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-out/detail/{id} [put]
func (r *ItemWarehouseTransferOutControllerImpl) UpdateTransferOutDetail(writer http.ResponseWriter, request *http.Request) {
	var transferOut transactionsparepartpayloads.UpdateItemWarehouseTransferOutDetailRequest

	transferOutSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	helper.ReadFromRequestBody(request, &transferOut)
	if validationErr := validation.ValidationForm(writer, request, &transferOut); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferOutService.UpdateTransferOutDetail(transferOut, transferOutSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "update success", http.StatusCreated)
}
