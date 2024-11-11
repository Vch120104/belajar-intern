package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
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

	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"
)

type ItemLocationController interface {
	GetAllItemLocation(writer http.ResponseWriter, request *http.Request)
	SaveItemLocation(writer http.ResponseWriter, request *http.Request)
	GetItemLocationById(writer http.ResponseWriter, request *http.Request)
	GetAllItemLocationDetail(writer http.ResponseWriter, request *http.Request)
	PopupItemLocation(writer http.ResponseWriter, request *http.Request)
	AddItemLocation(writer http.ResponseWriter, request *http.Request)
	DeleteItemLocation(writer http.ResponseWriter, request *http.Request)
	GetAllItemLoc(writer http.ResponseWriter, request *http.Request)
	GetByIdItemLoc(writer http.ResponseWriter, request *http.Request)
	SaveItemLoc(writer http.ResponseWriter, request *http.Request)
	DeleteItemLoc(writer http.ResponseWriter, request *http.Request)
	DownloadTemplate(writer http.ResponseWriter, request *http.Request)
	UploadTemplate(writer http.ResponseWriter, request *http.Request)
	ProcessUploadData(writer http.ResponseWriter, request *http.Request)
}

type ItemLocationControllerImpl struct {
	ItemLocationService masteritemservice.ItemLocationService
}

func NewItemLocationController(ItemLocationService masteritemservice.ItemLocationService) ItemLocationController {
	return &ItemLocationControllerImpl{
		ItemLocationService: ItemLocationService,
	}
}

// @Summary Get All Item Location Popup
// @Description REST API Item Location Popup
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_location_source_id query string false "item_location_source_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location/popup-location [get]
func (r *ItemLocationControllerImpl) PopupItemLocation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_item_location_source.item_location_source_id":   queryValues.Get("item_location_source_id"),
		"mtr_item_location_source.item_location_source_code": queryValues.Get("item_location_source_code"),
		"mtr_item_location_source.item_location_source_name": queryValues.Get("item_location_source_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.ItemLocationService.PopupItemLocation(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get All Item Location
// @Description REST API Item Location
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_name query int false "item_name"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location [get]
func (r *ItemLocationControllerImpl) GetAllItemLocation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_location_item.warehouse_group_id": queryValues.Get("warehouse_group_id"),
		"mtr_location_item.item_id":            queryValues.Get("item_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.ItemLocationService.GetAllItemLocation(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Save Item Location
// @Description REST API Item Location
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @param reqBody body masteritempayloads.ItemLocationRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location [post]
func (r *ItemLocationControllerImpl) SaveItemLocation(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.ItemLocationRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	create, err := r.ItemLocationService.SaveItemLocation(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if formRequest.ItemLocationId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
	}

}

// @Summary Get Item Location By ID
// @Description REST API  Item Location
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @Param item_location_id path int true "item_location_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location/{item_location_id} [get]
func (r *ItemLocationControllerImpl) GetItemLocationById(writer http.ResponseWriter, request *http.Request) {

	ItemLocationIds, errA := strconv.Atoi(chi.URLParam(request, "item_location_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ItemLocationService.GetItemLocationById(ItemLocationIds)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Item Location Detail
// @Description REST API Item Location Detail
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param item_location_detail_id query string false "item_location_detail_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location/detail [get]
func (r *ItemLocationControllerImpl) GetAllItemLocationDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	// Define query parameters
	queryParams := map[string]string{
		"mtr_item_location_detail.item_location_detail_id": queryValues.Get("item_location_detail_id"),
		"mtr_item_location_detail.item_location_id":        queryValues.Get("item_location_id"),
		"mtr_item_location_detail.item_location_source_id": queryValues.Get("item_location_source_id"),
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
	paginatedData, totalPages, totalRows, err := r.ItemLocationService.GetAllItemLocationDetail(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	// Construct the response
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Save Item Location Detail
// @Description REST API Item Location
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @Param item_location_id path int true "Item Location Detail ID"
// @param reqBody body masteritempayloads.ItemLocationDetailRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location/detail [post]
func (r *ItemLocationControllerImpl) AddItemLocation(writer http.ResponseWriter, request *http.Request) {
	itemLocID, errA := strconv.Atoi(chi.URLParam(request, "item_location_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var formRequest masteritempayloads.ItemLocationDetailRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	if err := r.ItemLocationService.AddItemLocation(int(itemLocID), formRequest); err != nil {
		exceptions.NewAppException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, nil, "Item location added successfully", http.StatusCreated)
}

// @Summary Delete Item Location By ID
// @Description REST API  Item Location
// @Accept json
// @Produce json
// @Tags Master : Item Location
// @Param item_location_detail_id path int true "item_location_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-location/detail/{item_location_detail_id} [delete]
func (r *ItemLocationControllerImpl) DeleteItemLocation(writer http.ResponseWriter, request *http.Request) {
	// Mendapatkan ID item lokasi dari URL
	itemLocationID, err := strconv.Atoi(chi.URLParam(request, "item_location_detail_id"))
	if err != nil {
		// Jika gagal mendapatkan ID dari URL, kirim respons error
		payloads.NewHandleError(writer, "Invalid item location ID", http.StatusBadRequest)
		return
	}

	// Memanggil service untuk menghapus item lokasi
	if deleteErr := r.ItemLocationService.DeleteItemLocation(itemLocationID); deleteErr != nil {
		// Jika terjadi kesalahan saat menghapus, kirim respons error
		exceptions.NewNotFoundException(writer, request, deleteErr)
		return
	}

	// Jika berhasil, kirim respons berhasil
	payloads.NewHandleSuccess(writer, nil, "Item location deleted successfully", http.StatusOK)
}

func (r *ItemLocationControllerImpl) GetAllItemLoc(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_warehouse_group.warehouse_group_name": queryValues.Get("warehouse_group_name"),
		"mtr_warehouse_group.warehouse_group_code": queryValues.Get("warehouse_group_code"),
		"mtr_warehouse_master.warehouse_id":        queryValues.Get("warehouse_id"),
		"mtr_warehouse_master.warehouse_code":      queryValues.Get("warehouse_code"),
		"mtr_warehouse_master.warehouse_name":      queryValues.Get("warehouse_name"),
		"mtr_item.item_id":                         queryValues.Get("item_id"),
	}
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	criteria := utils.BuildFilterCondition(queryParams)
	result, totalpage, totalrows, err := r.ItemLocationService.GetAllItemLoc(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(result) == 0 {
		payloads.NewHandleSuccessPagination(writer, result, "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalrows), totalpage)
		return
	}
	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalrows), totalpage)
}

func (r *ItemLocationControllerImpl) GetByIdItemLoc(writer http.ResponseWriter, request *http.Request) {
	ItemLocationIds, errA := strconv.Atoi(chi.URLParam(request, "item_location_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ItemLocationService.GetByIdItemLoc(ItemLocationIds)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemLocationControllerImpl) SaveItemLoc(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.SaveItemlocation
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	create, err := r.ItemLocationService.SaveItemLoc(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if formRequest.ItemLocationId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
	}
}

func (r *ItemLocationControllerImpl) DeleteItemLoc(writer http.ResponseWriter, request *http.Request) {
	itemlocationids := chi.URLParam(request, "item_location_id")
	itemlocationids = strings.Trim(itemlocationids, "[]")
	elements := strings.Split(itemlocationids, ",")

	itemLocIDs := []int{}
	for _, element := range elements {
		num, convErr := strconv.Atoi(strings.TrimSpace(element))
		if convErr != nil {
			payloads.NewHandleError(writer, "Failed to convert ID string", http.StatusInternalServerError)
			return
		}
		itemLocIDs = append(itemLocIDs, num)
	}
	if deleted, err := r.ItemLocationService.DeleteItemLoc(itemLocIDs); err != nil {
		exceptions.NewAppException(writer, request, err)
	} else if deleted {
		payloads.NewHandleSuccess(writer, deleted, "Delete Data Successfully!", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete data", http.StatusInternalServerError)
	}
}

func (r *ItemLocationControllerImpl) DownloadTemplate(writer http.ResponseWriter, request *http.Request) {
	f, errorGenerate := r.ItemLocationService.GenerateTemplateFile()

	if errorGenerate != nil {
		helper.ReturnError(writer, request, errorGenerate)
		return
	}

	var b bytes.Buffer
	err := f.Write(&b)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{StatusCode: 500, Err: errors.New("failed to write file to bytes")})
		return
	}

	downloadName := time.Now().UTC().Format("Template-Upload-ItemLocationMaster.xlsx")

	writer.Header().Set("Content-Description", "File Transfer")

	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)

	writer.Write(b.Bytes())
}

func (r *ItemLocationControllerImpl) UploadTemplate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "file size max 10MB",
			Err:        err,
		})
		return
	}

	file, handler, err := request.FormFile("file")
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Key name must be 'file'",
			Err:        err,
		})
		return
	}

	if !strings.HasSuffix(handler.Filename, ".xlsx") {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "File must be in xlsx format",
		})
		return
	}

	f, err := excelize.OpenReader(file)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error reading Excel file",
			Err:        err,
		})
		return
	}

	rows, err := f.GetRows("ItemLocationMaster")
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Pease check the sheet name must be ItemLocationMaster",
			Err:        err,
		})
		return
	}
	defer file.Close()

	previewData, errorPreview := r.ItemLocationService.UploadPreviewFile(rows)
	if errorPreview != nil {
		helper.ReturnError(writer, request, errorPreview)
		return
	}

	payloads.NewHandleSuccess(writer, previewData, "Preview Data Successfully!", http.StatusOK)
}

func (r *ItemLocationControllerImpl) ProcessUploadData(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "File size max 10MB",
			Err:        err,
		})
		return
	}

	file, handler, err := request.FormFile("file")
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Key name must be 'file'",
			Err:        err,
		})
		return
	}

	if !strings.HasSuffix(handler.Filename, ".xlsx") {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "File must be in xlsx format",
		})
		return
	}

	f, err := excelize.OpenReader(file)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error reading Excel file",
			Err:        err,
		})
		return
	}

	rows, err := f.GetRows("ItemLocationMaster")
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Pease check the sheet name must be ItemLocationMaster",
			Err:        err,
		})
		return
	}
	defer file.Close()

	previewData, errorPreview := r.ItemLocationService.UploadPreviewFile(rows)
	if errorPreview != nil {
		helper.ReturnError(writer, request, errorPreview)
		return
	}

	result, resultErr := r.ItemLocationService.UploadProcessFile(previewData)
	if resultErr != nil {
		exceptions.NewBadRequestException(writer, request, resultErr)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Process Data Successfully!", http.StatusOK)
}
