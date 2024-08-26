package masterwarehousecontroller

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"strconv"

	// masteritemlevelentities "after-sales/api/entities/master/item_level"

	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	// masteritemlevelrepo "after-sales/api/repositories/master/item_level"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"
)

type WarehouseLocationControllerImpl struct {
	WarehouseLocationService masterwarehouseservice.WarehouseLocationService
}

type WarehouseLocationController interface {
	GetAll(writer http.ResponseWriter, request *http.Request)
	GetById(writer http.ResponseWriter, request *http.Request)
	Save(writer http.ResponseWriter, request *http.Request)
	ChangeStatus(writer http.ResponseWriter, request *http.Request)
	DownloadTemplate(writer http.ResponseWriter, request *http.Request)
	UploadPreviewFile(writer http.ResponseWriter, request *http.Request)
	ProcessWarehouseLocationTemplate(writer http.ResponseWriter, request *http.Request)
}

func NewWarehouseLocationController(WarehouseLocationService masterwarehouseservice.WarehouseLocationService) WarehouseLocationController {
	return &WarehouseLocationControllerImpl{
		WarehouseLocationService: WarehouseLocationService,
	}
}

// ProcessWarehouseLocationTemplate implements WarehouseLocationController.
func (r *WarehouseLocationControllerImpl) ProcessWarehouseLocationTemplate(writer http.ResponseWriter, request *http.Request) {
	companyId, errA := strconv.Atoi(chi.URLParam(request, "company_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var formRequest masterwarehousepayloads.ProcessWarehouseLocationTemplate

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

	fmt.Print(formRequest)

	create, err := r.WarehouseLocationService.ProcessWarehouseLocationTemplate(formRequest, companyId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// UploadPreviewFile implements WarehouseLocationController.
func (r *WarehouseLocationControllerImpl) UploadPreviewFile(writer http.ResponseWriter, request *http.Request) {

	companyId, errA := strconv.Atoi(chi.URLParam(request, "company_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	// Parse the multipart form
	err := request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("file size max 10MB"), StatusCode: 400})
		return
	}

	// Retrieve the file from form data
	file, handler, err := request.FormFile("WarehouseLocation-File")
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("key name must be WarehouseLocation-File"), StatusCode: 400})
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
	rows, err := f.GetRows("WarehouseLocation")
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("please check the sheet name must be WarehouseLocation"), StatusCode: 400})
		return
	}

	previewData, errorPreview := r.WarehouseLocationService.UploadPreviewFile(rows, companyId)

	if errorPreview != nil {
		helper.ReturnError(writer, request, errorPreview)
		return
	}

	payloads.NewHandleSuccess(writer, previewData, "Get Data Successfully!", http.StatusOK)
}

// DownloadTemplate implements WarehouseLocationController.
func (r *WarehouseLocationControllerImpl) DownloadTemplate(writer http.ResponseWriter, request *http.Request) {
	f, errorGenerate := r.WarehouseLocationService.GenerateTemplateFile()

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

	downloadName := time.Now().UTC().Format("Template-Upload-WarehouseLocation.xlsx")

	writer.Header().Set("Content-Description", "File Transfer")

	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)

	writer.Write(b.Bytes())
}

// @Summary Get All Warehouse Location
// @Description Get All Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Success 200 {object} payloads.Response
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param is_active query bool false "Is Active"
// @Param warehouse_location_code query string false "Warehouse Location Code"
// @Param warehouse_location_name query string false "Warehouse Location Name"
// @Param company_id query string false "Company Id"
// @Param warehouse_location_detail_name query string false "Warehouse Location Detail Name"
// @Param sort_by query string false "Sort Of: {column}"
// @Param sort_of query string false "Sort By: {asc}"
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/ [get]
func (r *WarehouseLocationControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	filter := map[string]string{
		"WarehouseGroup.warehouse_group_name":            queryValues.Get("warehouse_group_name"),
		"mtr_warehouse_master.warehouse_code":            queryValues.Get("warehouse_code"),
		"mtr_warehouse_master.warehouse_name":            queryValues.Get("warehouse_name"),
		"mtr_warehouse_location.warehouse_location_code": queryValues.Get("warehouse_location_code"),
		"mtr_warehouse_location.warehouse_location_name": queryValues.Get("warehouse_location_name"),
		"mtr_warehouse_location.is_active":               queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(filter)

	result, err := r.WarehouseLocationService.GetAll(internalCriteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Get Warehouse Location By Id
// @Description Get Warehouse Location By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Param warehouse_location_id path int true "warehouse_location_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/{warehouse_location_id} [get]
func (r *WarehouseLocationControllerImpl) GetById(writer http.ResponseWriter, request *http.Request) {

	warehouseLocationId, errA := strconv.Atoi(chi.URLParam(request, "warehouse_location_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	get, err := r.WarehouseLocationService.GetById(warehouseLocationId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, get, "Get Data Successfully!", http.StatusOK)

}

// @Summary Save Warehouse Location
// @Description Save Warehouse Location
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @param reqBody body masterwarehousepayloads.GetWarehouseLocationResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/ [post]
func (r *WarehouseLocationControllerImpl) Save(writer http.ResponseWriter, request *http.Request) {
	var message string
	var formRequest masterwarehouseentities.WarehouseLocation
	helper.ReadFromRequestBody(request, &formRequest)

	save, err := r.WarehouseLocationService.Save(formRequest)

	if formRequest.WarehouseLocationId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, save, message, http.StatusOK)

}

// @Summary Change Warehouse Location Status By Id
// @Description Change Warehouse Location Status By Id
// @Accept json
// @Produce json
// @Tags Master : Warehouse Location
// @Param warehouse_location_id path int true "Warehouse Location Id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/warehouse-location/{warehouse_location_id} [patch]
func (r *WarehouseLocationControllerImpl) ChangeStatus(writer http.ResponseWriter, request *http.Request) {

	warehouseLocationId, errA := strconv.Atoi(chi.URLParam(request, "warehouse_location_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	change_status, err := r.WarehouseLocationService.ChangeStatus(warehouseLocationId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, change_status, "Updated successfully", http.StatusOK)

}
