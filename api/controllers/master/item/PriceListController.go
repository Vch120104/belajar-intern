package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/utils"
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	helper "after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/validation"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"
)

type PriceListController interface {
	GetPriceListLookup(writer http.ResponseWriter, request *http.Request)
	SavePriceList(writer http.ResponseWriter, request *http.Request)
	ChangeStatusPriceList(writer http.ResponseWriter, request *http.Request)
	DeletePriceList(writer http.ResponseWriter, request *http.Request)
	GetAllPriceListNew(writer http.ResponseWriter, request *http.Request)
	ActivatePriceList(writer http.ResponseWriter, request *http.Request)
	DeactivatePriceList(writer http.ResponseWriter, request *http.Request)
	GetPriceListById(writer http.ResponseWriter, request *http.Request)
	GetPriceListByCodeId(writer http.ResponseWriter, request *http.Request)
	GenerateDownloadTemplateFile(writer http.ResponseWriter, request *http.Request)
	UploadFile(writer http.ResponseWriter, request *http.Request)
	CheckPriceListItem(writer http.ResponseWriter, request *http.Request)
	Download(writer http.ResponseWriter, request *http.Request)
	Duplicate(writer http.ResponseWriter, request *http.Request)
}

type PriceListControllerImpl struct {
	pricelistservice masteritemservice.PriceListService
}

func NewPriceListController(PriceListService masteritemservice.PriceListService) PriceListController {
	return &PriceListControllerImpl{
		pricelistservice: PriceListService,
	}
}

// Duplicate implements PriceListController.
func (r *PriceListControllerImpl) Duplicate(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	brandId, errA := strconv.Atoi(queryValues.Get("brand_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	currencyId, errA := strconv.Atoi(queryValues.Get("currency_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	date := queryValues.Get("effective_date")
	itemGroupId, errA := strconv.Atoi(queryValues.Get("item_group_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	if brandId == 0 || currencyId == 0 || date == "" || itemGroupId == 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: 400, Err: errors.New("fill required params")})
		return
	}

	response, err := r.pricelistservice.Duplicate(itemGroupId, brandId, currencyId, date)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", 200)
}

// Download implements PriceListController.
func (r *PriceListControllerImpl) Download(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.PriceListUploadDataRequest

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

	f, errorGenerate := r.pricelistservice.Download(formRequest)

	if errorGenerate != nil {
		helper.ReturnError(writer, request, errorGenerate)
		return
	}

	// Write the Excel file to a buffer
	var b bytes.Buffer
	errWrite := f.Write(&b)
	if errWrite != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{StatusCode: 500, Err: errors.New("failed to write file to bytes")})
		return
	}

	downloadName := time.Now().UTC().Format("Download-PriceList.xlsx")

	writer.Header().Set("Content-Description", "File Transfer")

	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)

	writer.Write(b.Bytes())
}

// CheckPriceListItem implements PriceListController.
func (r *PriceListControllerImpl) CheckPriceListItem(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	brandId, errA := strconv.Atoi(queryValues.Get("brand_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	currencyId, errA := strconv.Atoi(queryValues.Get("currency_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	date := queryValues.Get("effective_date")
	itemGroupId, errA := strconv.Atoi(queryValues.Get("item_group_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	if brandId == 0 || currencyId == 0 || date == "" || itemGroupId == 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: 400, Err: errors.New("fill required params")})
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	response, err := r.pricelistservice.CheckPriceListItem(itemGroupId, brandId, currencyId, date, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, response.Rows, "Get Data Successfully!", 200, response.Limit, response.Page, response.TotalRows, response.TotalPages)
}

// UploadFile implements PriceListController.
func (r *PriceListControllerImpl) UploadFile(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.PriceListUploadDataRequest

	// Get the JSON part from the form data
	jsonPart := request.FormValue("data")
	if err := json.Unmarshal([]byte(jsonPart), &formRequest); err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: 400,
			Err:        err,
		})
		return
	}

	err := validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	// Parse the multipart form
	errParse := request.ParseMultipartForm(10 << 20) // 10 MB
	if errParse != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("file size max 10MB"), StatusCode: 400})
		return
	}

	// Retrieve the file from form data
	file, handler, errGetFile := request.FormFile("PriceList-File")
	if errGetFile != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("key name must be PriceList-File"), StatusCode: 400})
		return
	}
	defer file.Close()

	//Check file is XML
	if !strings.Contains(handler.Header.Get("Content-Type"), "xml") {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("make sure to upload xml file"), StatusCode: 400})
		return
	}
	// Read the uploaded file into an excelize.File
	f, errReadFile := excelize.OpenReader(file)
	if errReadFile != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errReadFile, StatusCode: 500})
		return
	}

	// Get all the rows in the ItemImportMaster.
	rows, errGetRows := f.GetRows("PriceList")
	if errGetRows != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{Err: errors.New("please check the sheet name must be PriceList"), StatusCode: 400})
		return
	}

	previewData, errorPreview := r.pricelistservice.UploadFile(rows, formRequest)

	if errorPreview != nil {
		helper.ReturnError(writer, request, errorPreview)
		return
	}

	payloads.NewHandleSuccess(writer, previewData, "Get Data Successfully!", http.StatusOK)
}

// GenerateDownloadTemplateFile implements PriceListController.
func (r *PriceListControllerImpl) GenerateDownloadTemplateFile(writer http.ResponseWriter, request *http.Request) {
	f, errorGenerate := r.pricelistservice.GenerateDownloadTemplateFile()

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

	downloadName := time.Now().UTC().Format("Template-Upload-PriceList.xlsx")

	writer.Header().Set("Content-Description", "File Transfer")

	writer.Header().Set("Content-Disposition", "attachment; filename="+downloadName)

	writer.Write(b.Bytes())
}

// GetPriceListById implements PriceListController.
func (r *PriceListControllerImpl) GetPriceListById(writer http.ResponseWriter, request *http.Request) {
	PriceListId, errA := strconv.Atoi(chi.URLParam(request, "price_list_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.pricelistservice.GetPriceListById(PriceListId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Success!", http.StatusOK)

}

// @Summary Get All Price List Lookup
// @Description REST API Price List
// @Param price_list_code query string false "price_list_code"
// @Param company_id query int false "company_id"
// @Param brand_id query int false "brand_id"
// @Param currency_id query int false "currency_id"
// @Param effective_date query string false "effective_date"
// @Param item_group_id query int false "item_group_id"
// @Param item_class_id query int false "item_class_id"
// @Accept json
// @Produce json
// @Tags Master : Price List
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/price-list/lookup [get]
func (r *PriceListControllerImpl) GetPriceListLookup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	PriceListCodeId := queryValues.Get("price_list_code_id")
	companyId, errA := strconv.Atoi(queryValues.Get("company_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	brandId, errA := strconv.Atoi(queryValues.Get("brand_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	currencyId, errA := strconv.Atoi(queryValues.Get("currency_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	effectiveDate, errA := time.Parse("2006-01-02T15:04:05.000Z", queryValues.Get("effective_date"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	itemGroupId, errA := strconv.Atoi(queryValues.Get("item_group_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	itemClassId, errA := strconv.Atoi(queryValues.Get("item_class_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	priceListRequest := masteritempayloads.PriceListGetAllRequest{
		PriceListCode: PriceListCodeId,
		CompanyId:     companyId,
		BrandId:       brandId,
		CurrencyId:    currencyId,
		EffectiveDate: effectiveDate,
		ItemGroupId:   itemGroupId,
		ItemClassId:   itemClassId,
	}

	result, err := r.pricelistservice.GetPriceList(priceListRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Save Price List
// @Description REST API Price List
// @Accept json
// @Produce json
// @Tags Master : Price List
// @param reqBody body masteritempayloads.PriceListResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/price-list/ [post]
func (r *PriceListControllerImpl) SavePriceList(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.SavePriceListMultiple

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

	create, err := r.pricelistservice.SavePriceList(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Success", http.StatusOK)
}

// @Summary Change Status Price List
// @Description REST API Price List
// @Accept json
// @Produce json
// @Tags Master : Price List
// @param price_list_id path int true "price_list_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/price-list/{price_list_id} [patch]
func (r *PriceListControllerImpl) ChangeStatusPriceList(writer http.ResponseWriter, request *http.Request) {

	PriceListId, _ := strconv.Atoi(chi.URLParam(request, "price_list_id"))

	response, err := r.pricelistservice.ChangeStatusPriceList(int(PriceListId))

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

func (r *PriceListControllerImpl) GetAllPriceListNew(writer http.ResponseWriter, request *http.Request) {

	queryValues := request.URL.Query()

	if mandatoryParamExist := queryValues.Get("price_list_code_id"); mandatoryParamExist == "" {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: 400, Err: errors.New("must input price list code")})
		return
	}

	queryParams := map[string]string{
		"brand_id":                          queryValues.Get("brand_id"),
		"mtr_item_price_list.item_group_id": queryValues.Get("item_group_id"),
		"price_list_code_id":                queryValues.Get("price_list_code_id"),
		"mtr_item_price_list.item_class_id": queryValues.Get("item_class_id"),
		"currency_id":                       queryValues.Get("currency_id"),
		"effective_date":                    queryValues.Get("effective_date"),
		"company_id":                        queryValues.Get("company_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.pricelistservice.GetAllPriceListNew(criteria, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

func (r *PriceListControllerImpl) ActivatePriceList(writer http.ResponseWriter, request *http.Request) {
	PriceListId := chi.URLParam(request, "price_list_id")
	response, err := r.pricelistservice.ActivatePriceList(PriceListId)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Activate data successfully!", http.StatusOK)
}

func (r *PriceListControllerImpl) DeactivatePriceList(writer http.ResponseWriter, request *http.Request) {
	PriceListId := chi.URLParam(request, "price_list_id")
	response, err := r.pricelistservice.DeactivatePriceList(PriceListId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Deactivate data successfully!", http.StatusOK)
}

func (r *PriceListControllerImpl) DeletePriceList(writer http.ResponseWriter, request *http.Request) {
	priceListId := chi.URLParam(request, "price_list_id")
	response, err := r.pricelistservice.DeletePriceList(priceListId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Delete data successfully!", http.StatusOK)
}

func (r *PriceListControllerImpl) GetPriceListByCodeId(writer http.ResponseWriter, request *http.Request) {
	priceListCodeId := chi.URLParam(request, "price_list_code_id")
	if priceListCodeId == "" {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("price_list_code_id parameter is required"),
		})
		return
	}

	response, err := r.pricelistservice.GetPriceListByCodeId(priceListCodeId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, response, "Get data successfully!", http.StatusOK)
}
