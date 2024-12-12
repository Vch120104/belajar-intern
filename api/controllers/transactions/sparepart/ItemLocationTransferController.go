package transactionsparepartcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
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

type ItemLocationTransferControllerImpl struct {
	ItemLocationTransferService transactionsparepartservice.ItemLocationTransferService
}

type ItemLocationTransferController interface {
	GetAllItemLocationTransfer(writer http.ResponseWriter, request *http.Request)
	GetItemLocationTransferById(writer http.ResponseWriter, request *http.Request)
	InsertItemLocationTransfer(writer http.ResponseWriter, request *http.Request)
	UpdateItemLocationTransfer(writer http.ResponseWriter, request *http.Request)
	AcceptItemLocationTransfer(writer http.ResponseWriter, request *http.Request)
}

func NewItemLocationTransferController(
	itemLocationTransferServiceService transactionsparepartservice.ItemLocationTransferService,
) ItemLocationTransferController {
	return &ItemLocationTransferControllerImpl{
		ItemLocationTransferService: itemLocationTransferServiceService,
	}
}
func (c *ItemLocationTransferControllerImpl) GetAllItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_item_warehouse_transfer_request.company_id":                       queryValues.Get("company_id"),
		"trx_item_warehouse_transfer_request.transfer_request_document_number": queryValues.Get("transfer_request_document_number"),
		"trx_item_warehouse_transfer_request.transfer_request_date_from":       queryValues.Get("transfer_request_date_from"),
		"trx_item_warehouse_transfer_request.transfer_request_date_to":         queryValues.Get("transfer_request_date_to"),
		"RequestFromWarehouse.warehouse_group_id":                              queryValues.Get("request_from_warehouse_group_id"),
		"RequestToWarehouse.warehouse_group_id":                                queryValues.Get("request_to_warehouse_group_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	response, err := c.ItemLocationTransferService.GetAllItemLocationTransfer(criteria, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		response.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		response.Limit,
		response.Page,
		int64(response.TotalRows),
		response.TotalPages,
	)
}

func (c *ItemLocationTransferControllerImpl) GetItemLocationTransferById(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_system_number"))

	response, err := c.ItemLocationTransferService.GetItemLocationTransferById(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully", http.StatusOK)
}

func (c *ItemLocationTransferControllerImpl) InsertItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	formRequest := transactionsparepartpayloads.InsertItemLocationTransferRequest{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	response, err := c.ItemLocationTransferService.InsertItemLocationTransfer(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Create Data Successfully", http.StatusOK)
}

func (c *ItemLocationTransferControllerImpl) UpdateItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_system_number"))

	formRequest := transactionsparepartpayloads.UpdateItemLocationTransferRequest{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	response, err := c.ItemLocationTransferService.UpdateItemLocationTransfer(transferRequestSystemNumber, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully", http.StatusOK)
}

func (c *ItemLocationTransferControllerImpl) AcceptItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_system_number"))

	formRequest := transactionsparepartpayloads.AcceptItemLocationTransferRequest{}
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	response, err := c.ItemLocationTransferService.AcceptItemLocationTransfer(transferRequestSystemNumber, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully", http.StatusOK)
}
