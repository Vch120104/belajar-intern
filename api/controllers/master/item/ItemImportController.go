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
	"fmt"
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
}

type ItemImportControllerImpl struct {
	ItemImportService masteritemservice.ItemImportService
}

// UploadTemplate implements ItemImportController.
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

	response := []masteritempayloads.ItemImportUploadResponse{}

	for index, value := range rows {
		data := masteritempayloads.ItemImportUploadResponse{}
		var failedQtyParse error
		var failedOrderaParse error
		if index > 0 {
			data.ItemCode = value[0]
			data.SupplierCode = value[1]
			data.ItemAliasCode = value[2]
			data.ItemAliasName = value[3]
			data.OrderQtyMultiplier, failedQtyParse = strconv.ParseFloat(value[4], 64)

			if failedQtyParse != nil {
				helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("make sure moq value is correct"), StatusCode: 400})
				return
			}

			data.RoyaltyFlag = value[5]
			data.OrderConversion, failedOrderaParse = strconv.ParseFloat(value[6], 64)

			if failedOrderaParse != nil {
				helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("make sure order conversion value is correct"), StatusCode: 400})
				return
			}

			response = append(response, data)
		}
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)

}

// DownloadTemplate implements ItemImportController.
func (r *ItemImportControllerImpl) DownloadTemplate(writer http.ResponseWriter, request *http.Request) {

	f := excelize.NewFile()
	sheetName := "ItemImportMaster"
	defer func() {
		if err := f.Close(); err != nil {
			helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500})
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500})
		return
	}
	// Set value of a cell.USPG_GMITEM2_INSERT
	f.SetCellValue(sheetName, "A1", "Part_Number")
	f.SetCellValue(sheetName, "B1", "Supplier_Code")
	f.SetCellValue(sheetName, "C1", "Part_Number_Alias")
	f.SetCellValue(sheetName, "D1", "Part_Name_Alias")
	f.SetCellValue(sheetName, "E1", "MOQ")
	f.SetCellValue(sheetName, "F1", "Royalty")
	f.SetCellValue(sheetName, "G1", "Order_Conversion")
	f.SetColWidth(sheetName, "A", "G", 21.5)

	// Create a style with bold font and border
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	})
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: err, StatusCode: 500})
	}

	// Apply the style to the header cells
	for col := 'A'; col <= 'G'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Get data example

	id := []int{}

	internalFilterCondition := map[string]string{}
	externalFilterCondition := map[string]string{}

	paginate := pagination.Pagination{
		Limit: 3,
		Page:  0,
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginatedData, _, _, _ := r.ItemImportService.GetAllItemImport(internalCriteria, externalCriteria, paginate)

	data, _ := masteritempayloads.ConvertItemImportMapToStruct(paginatedData)

	for _, value := range data {
		id = append(id, value.ItemImportId)
	}

	for i := 0; i < len(id); i++ {

		result, _ := r.ItemImportService.GetItemImportbyId(id[i])

		fmt.Println(result)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), result.ItemCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), result.SupplierCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), result.ItemAliasCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), result.ItemAliasName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), result.OrderQtyMultiplier)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), result.RoyaltyFlag)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), result.OrderConversion)

	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Write the Excel file to a buffer
	var b bytes.Buffer
	err = f.Write(&b)
	if err != nil {
		http.Error(writer, "Failed to write file to buffer.", http.StatusInternalServerError)
		return
	}

	downloadName := time.Now().UTC().Format("Template-Upload-ItemImportMaster.xlsx")

	writer.Header().Set("Content-Description", "File Transfer")

	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)

	writer.Write(b.Bytes())

}

// GetItemImportbyItemIdandSupplierId implements ItemImportController.
func (r *ItemImportControllerImpl) GetItemImportbyItemIdandSupplierId(writer http.ResponseWriter, request *http.Request) {
	itemId, _ := strconv.Atoi(chi.URLParam(request, "item_id"))
	supplierId, _ := strconv.Atoi(chi.URLParam(request, "supplier_id"))

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
// @Param item_import_id path int true "Item Import ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-import/{item_import_id} [get]
func (r *ItemImportControllerImpl) GetItemImportbyId(writer http.ResponseWriter, request *http.Request) {

	itemPackageId, _ := strconv.Atoi(chi.URLParam(request, "item_import_id"))

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
		"mtr_item.item_code": queryValues.Get("item_code"),
		"mtr_item.item_name": queryValues.Get("item_name"),
	}
	externalFilterCondition := map[string]string{
		"mtr_supplier.supplier_code": queryValues.Get("supplier_code"),
		"mtr_supplier.supplier_name": queryValues.Get("supplier_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginatedData, totalPages, totalRows, err := r.ItemImportService.GetAllItemImport(internalCriteria, externalCriteria, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// SaveItemImport implements ItemImportController.
// @Summary Save Item Import
// @Description Create a new item import
// @Accept json
// @Produce json
// @Tags Master : Item Import
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
