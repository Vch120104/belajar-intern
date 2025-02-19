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
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"
)

type ItemWarehouseTransferRequestController interface {
	InsertWhTransferRequestHeader(writer http.ResponseWriter, request *http.Request)
	InsertWhTransferRequestDetail(writer http.ResponseWriter, request *http.Request)
	UpdateWhTransferRequest(writer http.ResponseWriter, request *http.Request)
	UpdateWhTransferRequestDetail(writer http.ResponseWriter, request *http.Request)
	SubmitWhTransferRequest(writer http.ResponseWriter, request *http.Request)
	GetAllDetailTransferRequest(writer http.ResponseWriter, request *http.Request)
	GetByIdTransferRequest(writer http.ResponseWriter, request *http.Request)
	GetByIdTransferRequestDetail(writer http.ResponseWriter, request *http.Request)
	GetAllWhTransferRequest(writer http.ResponseWriter, request *http.Request)
	GetTransferRequestLookUp(writer http.ResponseWriter, request *http.Request)
	GetTransferRequestLookUpDetail(writer http.ResponseWriter, request *http.Request)
	DeleteHeaderTransferRequest(writer http.ResponseWriter, request *http.Request)
	DeleteDetail(writer http.ResponseWriter, request *http.Request)
	Upload(writer http.ResponseWriter, request *http.Request)
	ProcessUpload(writer http.ResponseWriter, request *http.Request)
	DownloadTemplate(writer http.ResponseWriter, request *http.Request)

	Accept(writer http.ResponseWriter, request *http.Request)
	Reject(writer http.ResponseWriter, request *http.Request)
	GetAllWhTransferReceipt(writer http.ResponseWriter, request *http.Request)
}

func NewItemWarehouseTransferRequestControllerImpl(itemWarehouseTransferRequestService transactionsparepartservice.ItemWarehouseTransferRequestService) ItemWarehouseTransferRequestController {
	return &ItemWarehouseTransferRequestControllerImpl{
		ItemWarehouseTransferRequestService: itemWarehouseTransferRequestService,
	}
}

type ItemWarehouseTransferRequestControllerImpl struct {
	ItemWarehouseTransferRequestService transactionsparepartservice.ItemWarehouseTransferRequestService
}

// GetTransferRequestLookUpDetail implements ItemWarehouseTransferRequestController.
func (r *ItemWarehouseTransferRequestControllerImpl) GetTransferRequestLookUpDetail(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"it.item_code":     queryValues.Get("item_code"),
		"it.item_name":     queryValues.Get("item_name"),
		"uom.uom_code":     queryValues.Get("unit_of_measurement"),
		"request_quantity": queryValues.Get("request_quantity"),
	}

	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := r.ItemWarehouseTransferRequestService.GetTransferRequestDetailLookUp(transferRequestSystemNumber, paginations, filterCondition)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// GetAllWhTransferReceipt implements ItemWarehouseTransferRequestController.
func (r *ItemWarehouseTransferRequestControllerImpl) GetAllWhTransferReceipt(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"transfer_request_status_id":                     queryValues.Get("transfer_request_status_id"),
		"transfer_request_document_number":               queryValues.Get("item_group_id"),
		"wmt.warehouse_group_id":                         queryValues.Get("transfer_request_warehouse_group_id"),
		"trx_item_warehouse_transfer_request.company_id": queryValues.Get("company_id"),
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
	res, err := r.ItemWarehouseTransferRequestService.GetAllWhTransferReceipt(paginations, filterCondition, dateParams)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// GetTransferRequestLookUp implements ItemWarehouseTransferRequestController.
func (r *ItemWarehouseTransferRequestControllerImpl) GetTransferRequestLookUp(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"b.company_id":                          queryValues.Get("company_id"),
		"b.transfer_request_document_number":    queryValues.Get("transfer_request_document_number"),
		"b.transfer_request_date":               queryValues.Get("transfer_request_date"),
		"b.transfer_request_by_id":              queryValues.Get("transfer_request_by_id"),
		"wmf.request_from_warehouse_name":       queryValues.Get("request_from_warehouse_name"),
		"wgf.request_from_warehouse_group_name": queryValues.Get("request_from_warehouse_group_name"),
		// "b.transfer_request_date": queryValues.Get("transfer_request_date"),
	}

	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)
	res, err := r.ItemWarehouseTransferRequestService.GetTransferRequestLookUp(paginations, filterCondition)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// @Summary Accept Item Warehouse Transfer Request
// @Description Accept Item Warehouse Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request System Number"
// @Param AcceptWarehouseTransferRequestRequest body transactionsparepartpayloads.AcceptWarehouseTransferRequestRequest true "Accept Warehouse Transfer Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/accept/{id} [put]
func (r *ItemWarehouseTransferRequestControllerImpl) Accept(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var transferRequest transactionsparepartpayloads.AcceptWarehouseTransferRequestRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferRequestService.AcceptTransferReceipt(transferRequestSystemNumber, transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// @Summary Reject Item Warehouse Transfer Request
// @Description Reject Item Warehouse Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request System Number"
// @Param RejectWarehouseTransferRequestRequest body transactionsparepartpayloads.RejectWarehouseTransferRequestRequest true "Reject Warehouse Transfer Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/receipt/reject/{id} [put]
func (r *ItemWarehouseTransferRequestControllerImpl) Reject(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var transferRequest transactionsparepartpayloads.RejectWarehouseTransferRequestRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferRequestService.RejectTransferReceipt(transferRequestSystemNumber, transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// @Summary Download Template
// @Description Download Template
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/download [get]
func (r *ItemWarehouseTransferRequestControllerImpl) DownloadTemplate(writer http.ResponseWriter, request *http.Request) {
	f, err := r.ItemWarehouseTransferRequestService.GenerateTemplateFile()
	if err != nil {
		// Return error response if template generation fails
		helper.ReturnError(writer, request, err)
		return
	}

	var b bytes.Buffer
	if err := f.Write(&b); err != nil {
		// Create BaseErrorResponse for file write error
		baseErr := &exceptions.BaseErrorResponse{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		exceptions.NewAppException(writer, request, baseErr)
		return
	}

	downloadName := time.Now().UTC().Format("2006-01-02_15-04-05") + "_ItemWHTransferRequest.xlsx"
	writer.Header().Set("Content-Description", "File Transfer")
	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)
	writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Expires", "0")
	writer.Header().Set("Cache-Control", "must-revalidate")
	writer.Header().Set("Pragma", "public")

	// Write the buffer to the HTTP response
	_, writeErr := writer.Write(b.Bytes())
	if writeErr != nil {
		// Create BaseErrorResponse for writer.Write error
		baseErr := &exceptions.BaseErrorResponse{
			Err:        writeErr,
			StatusCode: http.StatusInternalServerError,
		}
		// Use a generic error handling function to respond with the error
		exceptions.NewAppException(writer, request, baseErr)
		return
	}
}

// @Summary Process Upload
// @Description Process Upload
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param transfer_request_system_number query int true "Transfer Request System Number"
// @Param modified_by_id query int true "Modified By ID"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/process-upload [post]
func (r *ItemWarehouseTransferRequestControllerImpl) ProcessUpload(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	transferRequestSystemNumber, errA := strconv.Atoi(queryValues.Get("transfer_request_system_number"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	userId, errs := strconv.Atoi(queryValues.Get("modified_by_id"))
	if errs != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	if err := request.ParseMultipartForm(10 << 20); err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error parsing multipart form",
			Err:        err,
		})
		return
	}

	// Retrieve the file from the form data
	file, handler, err := request.FormFile("file")
	if err != nil {
		//log.Printf("Error retrieving file from form data: %v", err) // Logging error
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error retrieving file from form data",
			Err:        err,
		})
		return
	}
	defer file.Close()

	// Log the filename for debugging
	//log.Printf("Received file: %s", handler.Filename)

	// Check that the file is an xlsx format
	if !strings.HasSuffix(handler.Filename, ".xlsx") {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "File must be in xlsx format",
			Err:        errors.New("file must be in xlsx format"),
		})
		return
	}

	// Read the uploaded file into an excelize.File
	f, err := excelize.OpenReader(file)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error reading Excel file",
			Err:        err,
		})
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving rows from sheet",
			Err:        err,
		})
		return
	}

	previewData, errResponse := r.ItemWarehouseTransferRequestService.PreviewUploadData(rows)
	if errResponse != nil {
		exceptions.NewNotFoundException(writer, request, errResponse)
		return
	}

	var formRequest transactionsparepartpayloads.UploadProcessItemWarehouseTransferRequestPayloads

	formRequest.TransferRequestDetails = previewData
	formRequest.TransferRequestSystemNumber = transferRequestSystemNumber
	formRequest.ModifiedById = userId

	create, errProc := r.ItemWarehouseTransferRequestService.ProcessUploadData(formRequest)
	if errProc != nil {
		exceptions.NewNotFoundException(writer, request, errProc)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create/Update Data Successfully!", http.StatusCreated)
}

// @Summary Upload
// @Description Upload
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/upload [post]
func (r *ItemWarehouseTransferRequestControllerImpl) Upload(writer http.ResponseWriter, request *http.Request) {
	// Parse the multipart form with a 10 MB limit
	if err := request.ParseMultipartForm(10 << 20); err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error parsing multipart form",
			Err:        err,
		})
		return
	}

	// Retrieve the file from the form data
	file, handler, err := request.FormFile("file")
	if err != nil {
		//log.Printf("Error retrieving file from form data: %v", err) // Logging error
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error retrieving file from form data",
			Err:        err,
		})
		return
	}
	defer file.Close()

	// Log the filename for debugging
	//log.Printf("Received file: %s", handler.Filename)

	// Check that the file is an xlsx format
	if !strings.HasSuffix(handler.Filename, ".xlsx") {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "File must be in xlsx format",
			Err:        errors.New("file must be in xlsx format"),
		})
		return
	}

	// Read the uploaded file into an excelize.File
	f, err := excelize.OpenReader(file)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error reading Excel file",
			Err:        err,
		})
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving rows from sheet",
			Err:        err,
		})
		return
	}

	previewData, errResponse := r.ItemWarehouseTransferRequestService.PreviewUploadData(rows)
	if errResponse != nil {
		exceptions.NewNotFoundException(writer, request, errResponse)
		return
	}

	payloads.NewHandleSuccess(writer, previewData, "Preview Data Successfully!", http.StatusOK)
}

// @Summary Insert Item Warehouse Transfer Request Detail
// @Description Insert Item Warehouse Transfer Request Detail
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request System Number"
// @Param InsertItemWarehouseTransferDetailRequest body transactionsparepartpayloads.InsertItemWarehouseTransferDetailRequest true "Insert Item Warehouse Transfer Detail Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/detail/{id} [get]
func (r *ItemWarehouseTransferRequestControllerImpl) GetByIdTransferRequestDetail(writer http.ResponseWriter, request *http.Request) {
	transferRequestDetailSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	success, err := r.ItemWarehouseTransferRequestService.GetByIdTransferRequestDetail(transferRequestDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// @Summary Update Item Warehouse Transfer Request Detail
// @Description Update Item Warehouse Transfer Request Detail
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request Detail System Number"
// @Param UpdateItemWarehouseTransferRequestDetailRequest body transactionsparepartpayloads.UpdateItemWarehouseTransferRequestDetailRequest true "Update Item Warehouse Transfer Request Detail Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/detail/{id} [put]
func (r *ItemWarehouseTransferRequestControllerImpl) UpdateWhTransferRequestDetail(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.UpdateItemWarehouseTransferRequestDetailRequest

	transferRequestDetailSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferRequestService.UpdateWhTransferRequestDetail(transferRequest, transferRequestDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "update success", http.StatusCreated)
}

// @Summary Delete Detail
// @Description Delete Detail
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path string true "Detail Multi ID"
// @Param DeleteDetailItemWarehouseTransferRequest body transactionsparepartpayloads.DeleteDetailItemWarehouseTransferRequest true "Delete Detail Item Warehouse Transfer Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/detail/{id} [delete]
func (r *ItemWarehouseTransferRequestControllerImpl) DeleteDetail(writer http.ResponseWriter, request *http.Request) {
	multiId := chi.URLParam(request, "id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid request detail multi ID", http.StatusBadRequest)
		return
	}

	var transferRequest transactionsparepartpayloads.DeleteDetailItemWarehouseTransferRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
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
	success, err := r.ItemWarehouseTransferRequestService.DeleteDetail(intIds, transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "delete success", http.StatusCreated)
}

// @Summary Delete Header Transfer Request
// @Description Delete Header Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request System Number"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/{id} [delete]
func (r *ItemWarehouseTransferRequestControllerImpl) DeleteHeaderTransferRequest(writer http.ResponseWriter, request *http.Request) {

	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	success, err := r.ItemWarehouseTransferRequestService.DeleteHeaderTransferRequest(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "delete success", http.StatusCreated)
}

// @Summary Get All Detail Transfer Request
// @Description Get All Detail Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param transfer_request_system_number query int true "Transfer Request System Number"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/detail [get]
func (r *ItemWarehouseTransferRequestControllerImpl) GetAllDetailTransferRequest(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	transferRequestNumber, _ := strconv.Atoi(queryValues.Get("transfer_request_system_number"))
	paginations := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	res, err := r.ItemWarehouseTransferRequestService.GetAllDetailTransferRequest(transferRequestNumber, paginations)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// @Summary Get All Warehouse Transfer Request
// @Description Get All Warehouse Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param transfer_request_status_id query int false "Transfer Request Status ID"
// @Param transfer_request_document_number query string false "Transfer Request Document Number"
// @Param transfer_request_warehouse_group_id query int false "Transfer Request Warehouse Group ID"
// @Param company_id query int false "Company ID"
// @Param transfer_request_date_from query string false "Transfer Request Date From"
// @Param transfer_request_date_to query string false "Transfer Request Date To"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Param sort_of query string false "Sort Of"
// @Param sort_by query string false "Sort By"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request [get]
func (r *ItemWarehouseTransferRequestControllerImpl) GetAllWhTransferRequest(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"transfer_request_status_id":                     queryValues.Get("transfer_request_status_id"),
		"transfer_request_document_number":               queryValues.Get("transfer_request_document_number"),
		"wmt.warehouse_group_id":                         queryValues.Get("transfer_request_warehouse_group_id"),
		"trx_item_warehouse_transfer_request.company_id": queryValues.Get("company_id"),
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
	res, err := r.ItemWarehouseTransferRequestService.GetAllWhTransferRequest(paginations, filterCondition, dateParams)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// @Summary Get By ID Transfer Request
// @Description Get By ID Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request System Number"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/{id} [get]
func (r *ItemWarehouseTransferRequestControllerImpl) GetByIdTransferRequest(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	success, err := r.ItemWarehouseTransferRequestService.GetByIdTransferRequest(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// @Summary Insert Item Warehouse Transfer Request Detail
// @Description Insert Item Warehouse Transfer Request Detail
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param InsertItemWarehouseTransferRequest body transactionsparepartpayloads.InsertItemWarehouseTransferRequest true "Insert Item Warehouse Transfer Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/detail [post]
func (r *ItemWarehouseTransferRequestControllerImpl) InsertWhTransferRequestDetail(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.InsertItemWarehouseTransferDetailRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferRequestService.InsertWhTransferRequestDetail(transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusCreated)
}

// @Summary Insert Item Warehouse Transfer Request Header
// @Description Insert Item Warehouse Transfer Request Header
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param InsertItemWarehouseTransferRequest body transactionsparepartpayloads.InsertItemWarehouseTransferRequest true "Insert Item Warehouse Transfer Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request [post]
func (r *ItemWarehouseTransferRequestControllerImpl) InsertWhTransferRequestHeader(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.InsertItemWarehouseTransferRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferRequestService.InsertWhTransferRequestHeader(transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "save success", http.StatusCreated)
}

// @Summary Submit Item Warehouse Transfer Request
// @Description Submit Item Warehouse Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request System Number"
// @Param SubmitItemWarehouseTransferRequest body transactionsparepartpayloads.SubmitItemWarehouseTransferRequest true "Submit Item Warehouse Transfer Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/submit/{id} [put]
func (r *ItemWarehouseTransferRequestControllerImpl) SubmitWhTransferRequest(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, errA := strconv.Atoi(chi.URLParam(request, "id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var transferRequest transactionsparepartpayloads.SubmitItemWarehouseTransferRequest

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferRequestService.SubmitWhTransferRequest(transferRequestSystemNumber, transferRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// @Summary Update Item Warehouse Transfer Request
// @Description Update Item Warehouse Transfer Request
// @Tags Transaction : Sparepart Item Warehouse Transfer Request
// @Accept json
// @Produce json
// @Param id path int true "Transfer Request System Number"
// @Param UpdateItemWarehouseTransferRequest body transactionsparepartpayloads.UpdateItemWarehouseTransferRequest true "Update Item Warehouse Transfer Request"
// @Success 201 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-warehouse-transfer-request/{id} [put]
func (r *ItemWarehouseTransferRequestControllerImpl) UpdateWhTransferRequest(writer http.ResponseWriter, request *http.Request) {
	var transferRequest transactionsparepartpayloads.UpdateItemWarehouseTransferRequest

	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	helper.ReadFromRequestBody(request, &transferRequest)
	if validationErr := validation.ValidationForm(writer, request, &transferRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	success, err := r.ItemWarehouseTransferRequestService.UpdateWhTransferRequest(transferRequest, transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "update success", http.StatusCreated)
}
