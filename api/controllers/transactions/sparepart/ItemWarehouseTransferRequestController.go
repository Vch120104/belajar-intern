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
	DeleteHeaderTransferRequest(writer http.ResponseWriter, request *http.Request)
	DeleteDetail(writer http.ResponseWriter, request *http.Request)
	Upload(writer http.ResponseWriter, request *http.Request)
	ProcessUpload(writer http.ResponseWriter, request *http.Request)
	DownloadTemplate(writer http.ResponseWriter, request *http.Request)
	Accept(writer http.ResponseWriter, request *http.Request)
	Reject(writer http.ResponseWriter, request *http.Request)
}

func NewItemWarehouseTransferRequestControllerImpl(itemWarehouseTransferRequestService transactionsparepartservice.ItemWarehouseTransferRequestService) ItemWarehouseTransferRequestController {
	return &ItemWarehouseTransferRequestControllerImpl{
		ItemWarehouseTransferRequestService: itemWarehouseTransferRequestService,
	}
}

type ItemWarehouseTransferRequestControllerImpl struct {
	ItemWarehouseTransferRequestService transactionsparepartservice.ItemWarehouseTransferRequestService
}

// Accept implements ItemWarehouseTransferRequestController.
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

// Reject implements ItemWarehouseTransferRequestController.
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

// DownloadTemplate implements ItemWarehouseTransferRequestController.
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

// ProcessUpload implements ItemWarehouseTransferRequestController.
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

// Upload implements ItemWarehouseTransferRequestController.
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

// GetByIdTransferRequestDetail implements ItemWarehouseTransferRequestController.
func (r *ItemWarehouseTransferRequestControllerImpl) GetByIdTransferRequestDetail(writer http.ResponseWriter, request *http.Request) {
	transferRequestDetailSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	success, err := r.ItemWarehouseTransferRequestService.GetByIdTransferRequestDetail(transferRequestDetailSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// UpdateWhTransferRequestDetail implements ItemWarehouseTransferRequestController.
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

// DeleteDetail implements ItemWarehouseTransferRequestController.
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

// DeleteHeaderTransferRequest implements ItemWarehouseTransferRequestController.
func (r *ItemWarehouseTransferRequestControllerImpl) DeleteHeaderTransferRequest(writer http.ResponseWriter, request *http.Request) {

	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	success, err := r.ItemWarehouseTransferRequestService.DeleteHeaderTransferRequest(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "delete success", http.StatusCreated)
}

// GetAllDetailTransferRequest implements ItemWarehouseTransferRequestController.
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

// GetAllWhTransferRequest implements ItemWarehouseTransferRequestController.
func (r *ItemWarehouseTransferRequestControllerImpl) GetAllWhTransferRequest(writer http.ResponseWriter, request *http.Request) {
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
	res, err := r.ItemWarehouseTransferRequestService.GetAllWhTransferRequest(paginations, filterCondition, dateParams)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, res.Rows, "Success Get All Data", 200, res.Limit, res.Page, res.TotalRows, res.TotalPages)
}

// GetByIdTransferRequest implements ItemWarehouseTransferRequestController.
func (r *ItemWarehouseTransferRequestControllerImpl) GetByIdTransferRequest(writer http.ResponseWriter, request *http.Request) {
	transferRequestSystemNumber, _ := strconv.Atoi(chi.URLParam(request, "id"))

	success, err := r.ItemWarehouseTransferRequestService.GetByIdTransferRequest(transferRequestSystemNumber)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, success, "Get Data Success", http.StatusCreated)
}

// InsertWhTransferRequestDetail implements ItemWarehouseTransferRequestController.
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

// InsertWhTransferRequestHeader implements ItemWarehouseTransferRequestController.
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

// SubmitWhTransferRequest implements ItemWarehouseTransferRequestController.
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

// UpdateWhTransferRequest implements ItemWarehouseTransferRequestController.
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
