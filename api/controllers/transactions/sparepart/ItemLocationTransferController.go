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
	"strings"

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
	RejectItemLocationTransfer(writer http.ResponseWriter, request *http.Request)
	SubmitItemLocationTransfer(writer http.ResponseWriter, request *http.Request)
	DeleteItemLocationTransfer(writer http.ResponseWriter, request *http.Request)

	GetAllItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request)
	GetItemLocationTransferDetailById(writer http.ResponseWriter, request *http.Request)
	InsertItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request)
	UpdateItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request)
	DeleteItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request)
}

func NewItemLocationTransferController(
	itemLocationTransferServiceService transactionsparepartservice.ItemLocationTransferService,
) ItemLocationTransferController {
	return &ItemLocationTransferControllerImpl{
		ItemLocationTransferService: itemLocationTransferServiceService,
	}
}

// @Summary Get All Item Location Transfer
// @Description Get All Item Location Transfer
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID"
// @Param transfer_request_document_number query string false "Transfer Request Document Number"
// @Param transfer_request_status_id query string false "Transfer Request Status ID"
// @Param transfer_request_date_from query string false "Transfer Request Date From"
// @Param transfer_request_date_to query string false "Transfer Request Date To"
// @Param request_from_warehouse_group_id query string false "Request From Warehouse Group ID"
// @Param request_to_warehouse_group_id query string false "Request To Warehouse Group ID"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer [get]
func (c *ItemLocationTransferControllerImpl) GetAllItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_item_warehouse_transfer_request.company_id":                       queryValues.Get("company_id"),
		"trx_item_warehouse_transfer_request.transfer_request_document_number": queryValues.Get("transfer_request_document_number"),
		"trx_item_warehouse_transfer_request.transfer_request_status_id":       queryValues.Get("transfer_request_status_id"),
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

// @Summary Get Item Location Transfer By ID
// @Description Get Item Location Transfer By ID
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_system_number path string true "Transfer Request System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/{transfer_request_system_number} [get]
func (c *ItemLocationTransferControllerImpl) GetItemLocationTransferById(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_system_number"))

	response, err := c.ItemLocationTransferService.GetItemLocationTransferById(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully", http.StatusOK)
}

// @Summary Insert Item Location Transfer
// @Description Insert Item Location Transfer
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param InsertItemLocationTransferRequest body transactionsparepartpayloads.InsertItemLocationTransferRequest true "Insert Item Location Transfer Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer [post]
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

// @Summary Update Item Location Transfer
// @Description Update Item Location Transfer
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_system_number path string true "Transfer Request System Number"
// @Param UpdateItemLocationTransferRequest body transactionsparepartpayloads.UpdateItemLocationTransferRequest true "Update Item Location Transfer Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/{transfer_request_system_number} [put]
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

// @Summary Accept Item Location Transfer
// @Description Accept Item Location Transfer
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_system_number path string true "Transfer Request System Number"
// @Param AcceptItemLocationTransferRequest body transactionsparepartpayloads.AcceptItemLocationTransferRequest true "Accept Item Location Transfer Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/accept/{transfer_request_system_number} [put]
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
	payloads.NewHandleSuccess(writer, response, "Accept Data Successfully", http.StatusOK)
}

// @Summary Reject Item Location Transfer
// @Description Reject Item Location Transfer
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_system_number path string true "Transfer Request System Number"
// @Param RejectItemLocationTransferRequest body transactionsparepartpayloads.RejectItemLocationTransferRequest true "Reject Item Location Transfer Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/reject/{transfer_request_system_number} [put]
func (c *ItemLocationTransferControllerImpl) RejectItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_system_number"))

	formRequest := transactionsparepartpayloads.RejectItemLocationTransferRequest{}
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

	response, err := c.ItemLocationTransferService.RejectItemLocationTransfer(transferRequestSystemNumber, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Reject Data Successfully", http.StatusOK)
}

// @Summary Submit Item Location Transfer
// @Description Submit Item Location Transfer
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_system_number path string true "Transfer Request System Number"
// @Param SubmitItemLocationTransferRequest body transactionsparepartpayloads.SubmitItemLocationTransferRequest true "Submit Item Location Transfer Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/submit/{transfer_request_system_number} [put]
func (c *ItemLocationTransferControllerImpl) SubmitItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_system_number"))

	formRequest := transactionsparepartpayloads.SubmitItemLocationTransferRequest{}
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

	response, err := c.ItemLocationTransferService.SubmitItemLocationTransfer(transferRequestSystemNumber, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Submit Data Successfully", http.StatusOK)
}

// @Summary Delete Item Location Transfer
// @Description Delete Item Location Transfer
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_system_number path string true "Transfer Request System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/{transfer_request_system_number} [delete]
func (c *ItemLocationTransferControllerImpl) DeleteItemLocationTransfer(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_system_number"))

	response, err := c.ItemLocationTransferService.DeleteItemLocationTransfer(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Delete Data Successfully", http.StatusOK)
}

// @Summary Get All Item Location Transfer Detail
// @Description Get All Item Location Transfer Detail
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_system_number query string true "Transfer Request System Number"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/detail [get]
func (c *ItemLocationTransferControllerImpl) GetAllItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"trx_item_warehouse_transfer_request_detail.transfer_request_system_number": queryValues.Get("transfer_request_system_number"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	response, err := c.ItemLocationTransferService.GetAllItemLocationTransferDetail(criteria, paginate)
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

// @Summary Get Item Location Transfer Detail By ID
// @Description Get Item Location Transfer Detail By ID
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_detail_system_number path string true "Transfer Request Detail System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/detail/{transfer_request_detail_system_number} [get]
func (c *ItemLocationTransferControllerImpl) GetItemLocationTransferDetailById(writer http.ResponseWriter, request *http.Request) {
	transferRequestDetailSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_detail_system_number"))

	response, err := c.ItemLocationTransferService.GetItemLocationTransferDetailById(transferRequestDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully", http.StatusOK)
}

// @Summary Insert Item Location Transfer Detail
// @Description Insert Item Location Transfer Detail
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param InsertItemLocationTransferDetailRequest body transactionsparepartpayloads.InsertItemLocationTransferDetailRequest true "Insert Item Location Transfer Detail Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/detail [post]
func (c *ItemLocationTransferControllerImpl) InsertItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request) {
	formRequest := transactionsparepartpayloads.InsertItemLocationTransferDetailRequest{}
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

	response, err := c.ItemLocationTransferService.InsertItemLocationTransferDetail(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Create Data Successfully", http.StatusOK)
}

// @Summary Update Item Location Transfer Detail
// @Description Update Item Location Transfer Detail
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param transfer_request_detail_system_number path string true "Transfer Request Detail System Number"
// @Param UpdateItemLocationTransferDetailRequest body transactionsparepartpayloads.UpdateItemLocationTransferDetailRequest true "Update Item Location Transfer Detail Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/detail/{transfer_request_detail_system_number} [put]
func (c *ItemLocationTransferControllerImpl) UpdateItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request) {
	transferRequesDetailtSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "transfer_request_detail_system_number"))

	formRequest := transactionsparepartpayloads.UpdateItemLocationTransferDetailRequest{}
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

	response, err := c.ItemLocationTransferService.UpdateItemLocationTransferDetail(transferRequesDetailtSystemNumber, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Update Data Successfully", http.StatusOK)
}

// @Summary Delete Item Location Transfer Detail
// @Description Delete Item Location Transfer Detail
// @Tags Transaction Sparepart: Item Location Transfer
// @Accept json
// @Produce json
// @Param multi_id path string true "Multi ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location-transfer/detail/{multi_id} [delete]
func (c *ItemLocationTransferControllerImpl) DeleteItemLocationTransferDetail(writer http.ResponseWriter, request *http.Request) {
	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	multiId = strings.Trim(multiId, "[]")
	elements := strings.Split(multiId, ",")

	var intIds []int
	for _, element := range elements {
		num, convertErr := strconv.Atoi(strings.TrimSpace(element))
		if convertErr != nil {
			payloads.NewHandleError(writer, "Error converting data to integer", http.StatusBadRequest)
			return
		}
		intIds = append(intIds, num)
	}

	response, err := c.ItemLocationTransferService.DeleteItemLocationTransferDetail(intIds)
	if err != nil {
		if err.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, err)
		}
		return
	}

	payloads.NewHandleSuccess(writer, response, "Delete Data Successfully", http.StatusOK)
}
