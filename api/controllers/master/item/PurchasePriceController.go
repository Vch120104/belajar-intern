package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/labstack/gommon/log"
	"github.com/xuri/excelize/v2"
)

type PurchasePriceController interface {
	GetAllPurchasePrice(writer http.ResponseWriter, request *http.Request)
	SavePurchasePrice(writer http.ResponseWriter, request *http.Request)
	UpdatePurchasePrice(writer http.ResponseWriter, request *http.Request)
	GetPurchasePriceById(writer http.ResponseWriter, request *http.Request)
	ChangeStatusPurchasePrice(writer http.ResponseWriter, request *http.Request)
	GetPurchasePriceDetailById(writer http.ResponseWriter, request *http.Request)
	GetPurchasePriceDetailByParam(writer http.ResponseWriter, request *http.Request)
	GetAllPurchasePriceDetail(writer http.ResponseWriter, request *http.Request)
	AddPurchasePrice(writer http.ResponseWriter, request *http.Request)
	UpdatePurchasePriceDetail(writer http.ResponseWriter, request *http.Request)
	DeletePurchasePrice(writer http.ResponseWriter, request *http.Request)
	ActivatePurchasePriceDetail(writer http.ResponseWriter, request *http.Request)
	DeactivatePurchasePriceDetail(writer http.ResponseWriter, request *http.Request)

	DownloadTemplate(writer http.ResponseWriter, request *http.Request)
	Upload(writer http.ResponseWriter, request *http.Request)
	ProcessDataUpload(writer http.ResponseWriter, request *http.Request)
	Download(writer http.ResponseWriter, request *http.Request)
}

type PurchasePriceControllerImpl struct {
	PurchasePriceService masteritemservice.PurchasePriceService
}

func NewPurchasePriceController(PurchasePriceService masteritemservice.PurchasePriceService) PurchasePriceController {
	return &PurchasePriceControllerImpl{
		PurchasePriceService: PurchasePriceService,
	}
}

// @Summary Get All Purchase Price
// @Description REST API Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_name query int false "item_name"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price [get]
func (r *PurchasePriceControllerImpl) GetAllPurchasePrice(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_purchase_price.purchase_price_id":             queryValues.Get("purchase_price_id"),
		"mtr_purchase_price.supplier_id":                   queryValues.Get("supplier_id"),
		"mtr_purchase_price.supplier_code":                 queryValues.Get("supplier_code"),
		"mtr_purchase_price.supplier_name":                 queryValues.Get("supplier_name"),
		"mtr_purchase_price.currency_id":                   queryValues.Get("currency_id"),
		"mtr_purchase_price.currency_code":                 queryValues.Get("currency_code"),
		"mtr_purchase_price.purchase_price_effective_date": queryValues.Get("purchase_price_effective_date"),
		"mtr_purchase_price.is_active":                     queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.PurchasePriceService.GetAllPurchasePrice(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Update Purchase Price
// @Description REST API Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @param reqBody body masteritempayloads.PurchasePriceRequest true "Form Request"
// @param purchase_price_id path int true "purchase_price_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/{purchase_price_id} [put]
func (r *PurchasePriceControllerImpl) UpdatePurchasePrice(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.PurchasePriceRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	PurchasePriceId, errA := strconv.Atoi(chi.URLParam(request, "purchase_price_id")) // Get Purchase Price ID from URL

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	update, err := r.PurchasePriceService.UpdatePurchasePrice(PurchasePriceId, formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	message = "Update Data Successfully!"
	payloads.NewHandleSuccess(writer, update, message, http.StatusOK)
}

// @Summary Save Purchase Price
// @Description REST API Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @param reqBody body masteritempayloads.PurchasePriceRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price [post]
func (r *PurchasePriceControllerImpl) SavePurchasePrice(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.PurchasePriceRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	create, err := r.PurchasePriceService.SavePurchasePrice(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if formRequest.PurchasePriceId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
	}

}

// @Summary Get Purchase Price By ID
// @Description REST API  Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param purchase_price_id path int true "purchase_price_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/by-id/{purchase_price_id} [get]
func (r *PurchasePriceControllerImpl) GetPurchasePriceById(writer http.ResponseWriter, request *http.Request) {

	PurchasePriceIds, errA := strconv.Atoi(chi.URLParam(request, "purchase_price_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	queryValues := request.URL.Query()

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.PurchasePriceService.GetPurchasePriceById(PurchasePriceIds, paginate)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Change Status Purchase Price
// @Description REST API Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @param purchase_price_id path int true "purchase_price_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/{purchase_price_id} [patch]
func (r *PurchasePriceControllerImpl) ChangeStatusPurchasePrice(writer http.ResponseWriter, request *http.Request) {

	PurchasePricesId, errA := strconv.Atoi(chi.URLParam(request, "purchase_price_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	entity, err := r.PurchasePriceService.ChangeStatusPurchasePrice(int(PurchasePricesId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	responseData := map[string]interface{}{
		"is_active":         entity.IsActive,
		"purchase_price_id": entity.PurchasePriceId,
	}

	payloads.NewHandleSuccess(writer, responseData, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get All Purchase Price Detail
// @Description REST API Purchase Price Detail
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param purchase_price_detail_id query string false "purchase_price_detail_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail  [get]
func (r *PurchasePriceControllerImpl) GetAllPurchasePriceDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	// Define query parameters
	queryParams := map[string]string{
		"mtr_purchase_price_detail.purchase_price_detail_id": queryValues.Get("purchase_price_detail_id"),
		"mtr_purchase_price_detail.purchase_price_id":        queryValues.Get("purchase_price_id"),
	}

	// Extract pagination parameters
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	// Build filter condition based on query parameters
	criteria := utils.BuildFilterCondition(queryParams)

	// Call service to get paginated data
	paginatedData, totalPages, totalRows, err := r.PurchasePriceService.GetAllPurchasePriceDetail(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	// Construct the response
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Purchase Price Detail By Purchase Price ID
// @Description REST API  Purchase Price Detail
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param purchase_price_id path int true "purchase_price_id"
// @Success 200 {object} payloads.ResponsePagination
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail/{purchase_price_detail_id} [get]
func (r *PurchasePriceControllerImpl) GetPurchasePriceDetailById(writer http.ResponseWriter, request *http.Request) {
	PurchasePriceIds, errA := strconv.Atoi(chi.URLParam(request, "purchase_price_detail_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.PurchasePriceService.GetPurchasePriceDetailById(PurchasePriceIds)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Update Purchase Price Detail
// @Description REST API Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @param reqBody body masteritempayloads.PurchasePriceDetailRequest true "Form Request"
// @param purchase_price_detail_id path int true "purchase_price_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail/{purchase_price_detail_id} [put]
func (r *PurchasePriceControllerImpl) UpdatePurchasePriceDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.PurchasePriceDetailRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	PurchasePriceDetailId, errA := strconv.Atoi(chi.URLParam(request, "purchase_price_detail_id")) // Get Purchase Price ID from URL
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	update, err := r.PurchasePriceService.UpdatePurchasePriceDetail(PurchasePriceDetailId, formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	message = "Update Data Successfully!"
	payloads.NewHandleSuccess(writer, update, message, http.StatusOK)
}

// @Summary Save Purchase Price Detail
// @Description REST API Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @param reqBody body masteritempayloads.PurchasePriceDetailRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail [post]
func (r *PurchasePriceControllerImpl) AddPurchasePrice(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.PurchasePriceDetailRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	create, err := r.PurchasePriceService.AddPurchasePrice(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if formRequest.PurchasePriceDetailId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Delete Purchase Price By ID
// @Description REST API  Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param purchase_price_id path int true "purchase_price_id"
// @Param multi_id path string true "Purchase Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail/{purchase_price_id}/{multi_id} [get]
func (r *PurchasePriceControllerImpl) DeletePurchasePrice(writer http.ResponseWriter, request *http.Request) {

	PurchasePriceID, err := strconv.Atoi(chi.URLParam(request, "purchase_price_id"))
	if err != nil {

		payloads.NewHandleError(writer, "Invalid Purchase Price ID", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid service request detail multi ID", http.StatusBadRequest)
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

	success, baseErr := r.PurchasePriceService.DeletePurchasePrice(PurchasePriceID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Purchase detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Purchase Detail deleted successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete Purchase detail", http.StatusInternalServerError)
	}

}

// @Summary Activate Purchase Price Detail
// @Description REST API  Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param purchase_price_id path int true "purchase_price_id"
// @Param multi_id path string true "Purchase Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail/activate/{purchase_price_id}/{multi_id} [patch]
func (r *PurchasePriceControllerImpl) ActivatePurchasePriceDetail(writer http.ResponseWriter, request *http.Request) {
	PurchasePriceID, err := strconv.Atoi(chi.URLParam(request, "purchase_price_id"))
	if err != nil {
		payloads.NewHandleError(writer, "Invalid Purchase Price ID", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if strings.TrimSpace(multiId) == "[]" || multiId == "" {
		payloads.NewHandleError(writer, "Invalid service request detail multi ID", http.StatusBadRequest)
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

	success, baseErr := r.PurchasePriceService.ActivatePurchasePriceDetail(PurchasePriceID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Purchase detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Purchase Detail activated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to activate Purchase detail", http.StatusInternalServerError)
	}
}

// @Summary Deactivate Purchase Price Detail
// @Description REST API  Purchase Price
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param purchase_price_id path int true "purchase_price_id"
// @Param multi_id path string true "Purchase Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail/deactivate/{purchase_price_id}/{multi_id} [patch]
func (r *PurchasePriceControllerImpl) DeactivatePurchasePriceDetail(writer http.ResponseWriter, request *http.Request) {

	PurchasePriceID, err := strconv.Atoi(chi.URLParam(request, "purchase_price_id"))
	if err != nil {

		payloads.NewHandleError(writer, "Invalid Purchase Price ID", http.StatusBadRequest)
		return
	}

	multiId := chi.URLParam(request, "multi_id")
	if multiId == "[]" {
		payloads.NewHandleError(writer, "Invalid service request detail multi ID", http.StatusBadRequest)
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

	success, baseErr := r.PurchasePriceService.DeactivatePurchasePriceDetail(PurchasePriceID, intIds)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Purchase detail not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	if success {
		payloads.NewHandleSuccess(writer, success, "Purchase Detail deactivated successfully", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to deactivate Purchase detail", http.StatusInternalServerError)
	}
}

// DownloadTemplate godoc
// @Summary Download Template
// @Description REST API Download Template
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/download-template [get]
func (r *PurchasePriceControllerImpl) DownloadTemplate(writer http.ResponseWriter, request *http.Request) {
	// Generate the template file
	f, err := r.PurchasePriceService.GenerateTemplateFile()
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

	downloadName := time.Now().UTC().Format("2006-01-02_15-04-05") + "_Template_Upload_PriceListSupplier.xlsx"
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

// Upload godoc
// @Summary Upload
// @Description REST API Upload
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param file formData file true "File"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/upload [post]
func (r *PurchasePriceControllerImpl) Upload(writer http.ResponseWriter, request *http.Request) {
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
		log.Printf("Error retrieving file from form data: %v", err) // Logging error
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error retrieving file from form data",
			Err:        err,
		})
		return
	}
	defer file.Close()

	// Log the filename for debugging
	log.Printf("Received file: %s", handler.Filename)

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

	rows, err := f.GetRows("purchase_price")
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving rows from sheet",
			Err:        err,
		})
		return
	}

	purchasePriceID, err := strconv.Atoi(request.FormValue("purchase_price_id"))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error parsing purchase_price_id",
			Err:        err,
		})
		return
	}

	previewData, errResponse := r.PurchasePriceService.PreviewUploadData(rows, purchasePriceID)
	if errResponse != nil {
		exceptions.NewNotFoundException(writer, request, errResponse)
		return
	}

	payloads.NewHandleSuccess(writer, previewData, "Preview Data Successfully!", http.StatusOK)
}

// ProcessDataUpload godoc
// @Summary Process Data Upload
// @Description REST API Process Data Upload
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param file formData file true "File"
// @Param data formData string true "Data"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/process [post]
func (r *PurchasePriceControllerImpl) ProcessDataUpload(writer http.ResponseWriter, request *http.Request) {
	file, _, err := request.FormFile("file")
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error retrieving file from form data",
			Err:        err,
		})
		return
	}
	defer file.Close()

	// Read and process the file
	f, err := excelize.OpenReader(file)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error reading Excel file",
			Err:        err,
		})
		return
	}

	rows, err := f.GetRows("purchase_price")
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving rows from sheet",
			Err:        err,
		})
		return
	}

	purchasePriceID, err := strconv.Atoi(request.FormValue("purchase_price_id"))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Error parsing purchase_price_id",
			Err:        err,
		})
		return
	}

	data, errResp := r.PurchasePriceService.PreviewUploadData(rows, purchasePriceID)
	if errResp != nil {
		// Set status code from errResp
		writer.WriteHeader(errResp.StatusCode)
		json.NewEncoder(writer).Encode(errResp)
		return
	}

	// Directly use the data if it matches the expected type
	formRequest := masteritempayloads.UploadRequest{
		Data: data, // Use data of type []masteritempayloads.PurchasePriceDetailResponses
	}

	// Process the upload
	success, errResp := r.PurchasePriceService.ProcessDataUpload(formRequest)
	if errResp != nil {
		// Set status code from errResp
		writer.WriteHeader(errResp.StatusCode)
		json.NewEncoder(writer).Encode(errResp)
		return
	}

	if !success {
		// Set status code for internal server error
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(&exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to process data",
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	payloads.NewHandleSuccess(writer, nil, "Upload Data Successfully!", http.StatusOK)
}

// Download godoc
// @Summary Download
// @Description REST API Download
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/download [get]
func (r *PurchasePriceControllerImpl) Download(writer http.ResponseWriter, request *http.Request) {
	// Extract purchase_price_id from query parameters
	purchasePriceIDStr := request.URL.Query().Get("purchase_price_id")
	if purchasePriceIDStr == "" {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing purchase_price_id parameter",
		})
		return
	}

	purchasePriceID, err := strconv.Atoi(purchasePriceIDStr)
	if err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid purchase_price_id parameter",
			Err:        err,
		})
		return
	}

	// Fetch the data and save to file
	filePath, errResp := r.PurchasePriceService.DownloadData(purchasePriceID)
	if errResp != nil {
		exceptions.NewAppException(writer, request, errResp)
		return
	}

	// Open the file to be served
	file, err := os.Open(filePath)
	if err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	defer file.Close() // Ensure the file is closed

	// Set headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Expires", "0")
	writer.Header().Set("Cache-Control", "must-revalidate")
	writer.Header().Set("Pragma", "public")

	// Write the file content to the response
	if _, err := io.Copy(writer, file); err != nil {
		exceptions.NewAppException(writer, request, &exceptions.BaseErrorResponse{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	// Optionally, delete the temporary file after it's been served
	go func() {
		time.Sleep(1 * time.Second) // Give a small delay to ensure the file is released
		if err := os.Remove(filePath); err != nil {
			log.Errorf("Error deleting file: %v", err)
		}
	}()
}

// GetPurchasePriceDetailByParam godoc
// @Summary Get Purchase Price Detail By Param
// @Description REST API  Purchase Price Detail
// @Accept json
// @Produce json
// @Tags Master : Purchase Price
// @Param currency_id query int true "currency_id"
// @Param supplier_id query int true "supplier_id"
// @Param effective_date query string true "effective_date"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/purchase-price/detail/{currency_id}/{supplier_id}/{effective_date} [get]
func (r *PurchasePriceControllerImpl) GetPurchasePriceDetailByParam(writer http.ResponseWriter, request *http.Request) {
	currencyIDStr := chi.URLParam(request, "currency_id")
	currencyID, err := strconv.Atoi(currencyIDStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid currency ID", http.StatusBadRequest)
		return
	}

	supplierIDStr := chi.URLParam(request, "supplier_id")
	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid supplier ID", http.StatusBadRequest)
		return
	}

	effectiveDateStr := chi.URLParam(request, "effective_date")
	effectiveDate, err := time.Parse("2006-01-02T15:04:05.000Z", effectiveDateStr)
	if err != nil {
		payloads.NewHandleError(writer, "Invalid effective date", http.StatusBadRequest)
		return
	}
	effectiveDateFormatted := effectiveDate.Format("2006-01-02")

	result, baseErr := r.PurchasePriceService.GetPurchasePriceDetailByParam(currencyID, supplierID, effectiveDateFormatted)
	if baseErr != nil {
		if baseErr.StatusCode == http.StatusNotFound {
			payloads.NewHandleError(writer, "Purchase price detail data not found", http.StatusNotFound)
		} else {
			exceptions.NewAppException(writer, request, baseErr)
		}
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
