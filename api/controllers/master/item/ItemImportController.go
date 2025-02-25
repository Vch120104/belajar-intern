package masteritemcontroller

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/go-chi/chi/v5"
)

type ItemImportController interface {
	GetAllItemImport(writer http.ResponseWriter, request *http.Request)
	GetItemImportbyId(writer http.ResponseWriter, request *http.Request)
	SaveItemImport(writer http.ResponseWriter, request *http.Request)
	UpdateItemImport(writer http.ResponseWriter, request *http.Request)
	GetItemImportbyItemIdandSupplierId(writer http.ResponseWriter, request *http.Request)
	DownloadTemplate(writer http.ResponseWriter, request *http.Request)
	UploadTemplate(writer http.ResponseWriter, request *http.Request)
	ProcessDataUpload(writer http.ResponseWriter, request *http.Request)
}

type ItemImportControllerImpl struct {
	ItemImportService masteritemservice.ItemImportService
}

// @Summary Process Data Upload
// @Description Process data upload
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Param reqBody body masteritempayloads.ItemImportUploadRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/process-upload [post]
func (r *ItemImportControllerImpl) ProcessDataUpload(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemImportUploadRequest

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemImportService.ProcessDataUpload(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// @Summary Upload Template
// @Description Upload template
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Param ItemImportMaster-File formData file true "File to upload"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/upload-template [post]
func (r *ItemImportControllerImpl) UploadTemplate(writer http.ResponseWriter, request *http.Request) {
	// Parse the multipart form
	err := request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("file size max 10MB"), StatusCode: 500})
		return
	}

	// Retrieve the file from form data
	file, handler, err := request.FormFile("ItemImportMaster-File")
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("key name must be ItemImportMaster-File"), StatusCode: 401})
		return
	}
	defer file.Close()

	//Check file is XML
	if !strings.Contains(handler.Header.Get("Content-Type"), "xml") {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("make sure to upload xml file"), StatusCode: 400})
		return
	}
	// Read the uploaded file into an excelize.File
	f, err := excelize.OpenReader(file)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500})
		return
	}

	// Get all the rows in the ItemImportMaster.
	rows, err := f.GetRows("ItemImportMaster")
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("please check the sheet name must be ItemImportMaster"), StatusCode: 400})
		return
	}

	previewData, errorPreview := r.ItemImportService.UploadPreviewFile(rows)

	if errorPreview != nil {
		helper.ReturnError(writer, request, errorPreview)
		return
	}

	payloads.NewHandleSuccess(writer, previewData, "Get Data Successfully!", http.StatusOK)

}

// @Summary Download Template
// @Description Download template
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/download-template [get]
func (r *ItemImportControllerImpl) DownloadTemplate(writer http.ResponseWriter, request *http.Request) {

	f, errorGenerate := r.ItemImportService.GenerateTemplateFile()

	if errorGenerate != nil {
		helper.ReturnError(writer, request, errorGenerate)
		return
	}

	// Write the Excel file to a buffer
	var b bytes.Buffer
	err := f.Write(&b)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{StatusCode: 500, Err: errors.New("failed to write file to bytes")})
		return
	}

	downloadName := time.Now().UTC().Format("Template-Upload-ItemImportMaster.xlsx")

	writer.Header().Set("Content-Description", "File Transfer")

	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)

	writer.Write(b.Bytes())

}

// @Summary Get Item Import By Item ID and Supplier ID
// @Description Retrieve an item import by its item ID and supplier ID
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Param item_id path int true "Item ID"
// @Param supplier_id path int true "Supplier ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/{item_id}/{supplier_id} [get]
func (r *ItemImportControllerImpl) GetItemImportbyItemIdandSupplierId(writer http.ResponseWriter, request *http.Request) {
	itemId, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	supplierId, errA := strconv.Atoi(chi.URLParam(request, "supplier_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ItemImportService.GetItemImportbyItemIdandSupplierId(itemId, supplierId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// GetItemImportbyId implements ItemImportController.
// @Summary Get Item Import By ID
// @Description Retrieve an item import by its ID
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Param item_import_id path int true "Item Import ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/{item_import_id} [get]
func (r *ItemImportControllerImpl) GetItemImportbyId(writer http.ResponseWriter, request *http.Request) {

	itemPackageId, errA := strconv.Atoi(chi.URLParam(request, "item_import_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ItemImportService.GetItemImportbyId(itemPackageId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// GetAllItemImport implements ItemImportController.
// @Summary Get All Item Imports
// @Description Retrieve all item imports with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param item_code query string false "Item code"
// @Param item_name query string false "Item name"
// @Param supplier_code query string false "Supplier code"
// @Param supplier_name query string false "Supplier name"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import [get]
func (r *ItemImportControllerImpl) GetAllItemImport(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	internalFilterCondition := map[string]string{
		"mtr_item_import.item_import_id": queryValues.Get("item_import_id"),
		"mtr_item_import.item_id":        queryValues.Get("item_id"),
		"mtr_item_import.supplier_id":    queryValues.Get("supplier_id"),
		"Item.item_code":                 queryValues.Get("item_code"),
		"Item.item_name":                 queryValues.Get("item_name"),
	}
	externalFilterCondition := map[string]string{
		"supplier_code": queryValues.Get("supplier_code"),
		"supplier_name": queryValues.Get("supplier_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginatedData, err := r.ItemImportService.GetAllItemImport(internalCriteria, externalCriteria, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		utils.ModifyKeysInResponse(paginatedData.Rows),
		"Get Data Successfully!",
		http.StatusOK,
		paginate.Limit,
		paginate.Page,
		int64(paginatedData.TotalRows),
		paginatedData.TotalPages,
	)
}

// SaveItemImport implements ItemImportController.
// @Summary Save Item Import
// @Description Create a new item import
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Param reqBody body masteritempayloads.ItemImportResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/save [post]
func (r *ItemImportControllerImpl) SaveItemImport(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritementities.ItemImport
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemImportService.SaveItemImport(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// UpdateItemImport implements ItemImportController.
// @Summary Update Item Import
// @Description Update an existing item import
// @Accept json
// @Produce json
// @Tags Master : Item Import
// @Security AuthorizationKeyAuth
// @Param reqBody body masteritempayloads.ItemImportResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/update [put]
func (r *ItemImportControllerImpl) UpdateItemImport(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritementities.ItemImport
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemImportService.UpdateItemImport(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}

func NewItemImportController(ItemImportService masteritemservice.ItemImportService) ItemImportController {
	return &ItemImportControllerImpl{
		ItemImportService: ItemImportService,
	}
}
