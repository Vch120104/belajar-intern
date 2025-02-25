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

type BomController interface {
	// Parent
	GetBomList(writer http.ResponseWriter, request *http.Request)
	GetBomById(writer http.ResponseWriter, request *http.Request)
	GetBomByUn(writer http.ResponseWriter, request *http.Request)
	GetBomTotalPercentage(writer http.ResponseWriter, request *http.Request)
	ChangeStatusBomMaster(writer http.ResponseWriter, request *http.Request)
	UpdateBomMaster(writer http.ResponseWriter, request *http.Request)
	SaveBomMaster(writer http.ResponseWriter, request *http.Request)

	// Detail
	GetBomDetailByMasterId(writer http.ResponseWriter, request *http.Request)
	GetBomDetailByMasterUn(writer http.ResponseWriter, request *http.Request)
	GetBomDetailMaxSeq(writer http.ResponseWriter, request *http.Request)
	GetBomDetailById(writer http.ResponseWriter, request *http.Request)
	SaveBomDetail(writer http.ResponseWriter, request *http.Request)
	DeleteBomDetail(writer http.ResponseWriter, request *http.Request)

	// Excel
	Upload(writer http.ResponseWriter, request *http.Request)
	ProcessDataUpload(writer http.ResponseWriter, request *http.Request)
	DownloadTemplate(writer http.ResponseWriter, request *http.Request)
}

type BomControllerImpl struct {
	BomService masteritemservice.BomService
}

func NewBomController(bomService masteritemservice.BomService) BomController {
	return &BomControllerImpl{
		BomService: bomService,
	}
}

// @Summary Get All Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param bom_master_id query string false "bom_master_id"
// @Param item_name query string false "item_name"
// @Param effective_date query string false "effective_date"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/ [get]
func (r *BomControllerImpl) GetBomList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"item_code":      queryValues.Get("item_code"),
		"item_name":      queryValues.Get("item_name"),
		"effective_date": queryValues.Get("effective_date"),
		"qty":            queryValues.Get("qty"),
		"uom_code":       queryValues.Get("uom_code"),
		"is_active":      queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, err := r.BomService.GetBomList(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		paginatedData.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		paginate.Limit,
		paginate.Page,
		int64(paginatedData.TotalRows),
		paginatedData.TotalPages,
	)
}

// @Summary Get Bom Master By ID
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @Param bom_id path int true "bom_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/{bom_id} [get]
func (r *BomControllerImpl) GetBomById(writer http.ResponseWriter, request *http.Request) {
	bomId, errA := strconv.Atoi(chi.URLParam(request, "bom_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.BomService.GetBomById(bomId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Bom Master By Un
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @param item_id path int true "item_id"
// @param effective_date path string true "effective_date"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/{item_id}/{effective_date} [get]
func (r *BomControllerImpl) GetBomByUn(writer http.ResponseWriter, request *http.Request) {
	itemId, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	effectiveDate, errB := time.Parse("2006-01-02T15:04:05.000Z", chi.URLParam(request, "effective_date"))
	if errB != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.BomService.GetBomByUn(itemId, effectiveDate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Change Status Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @param bom_id path int true "bom_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/{bom_id} [patch]
func (r *BomControllerImpl) ChangeStatusBomMaster(writer http.ResponseWriter, request *http.Request) {
	bomMasterId, errA := strconv.Atoi(chi.URLParam(request, "bom_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	entity, err := r.BomService.ChangeStatusBomMaster(bomMasterId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	responseData := map[string]interface{}{
		"is_active": entity.IsActive,
		"bom_id":    entity.BomId,
	}

	payloads.NewHandleSuccess(writer, responseData, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get Bom Detail By Master ID
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_id path int true "bom_id"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/master/{bom_id} [get]
func (r *BomControllerImpl) GetBomDetailByMasterId(writer http.ResponseWriter, request *http.Request) {
	bomId, errA := strconv.Atoi(chi.URLParam(request, "bom_id"))
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

	paginatedData, err := r.BomService.GetBomDetailByMasterId(bomId, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		paginatedData.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		paginate.Limit,
		paginate.Page,
		int64(paginatedData.TotalRows),
		paginatedData.TotalPages,
	)
}

// @Summary Get Bom Detail By Master Un
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param item_id path int true "item_id"
// @Param effective_date path string true "effective_date"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/master/{item_id}/{effective_date} [get]
func (r *BomControllerImpl) GetBomDetailByMasterUn(writer http.ResponseWriter, request *http.Request) {
	itemId, errA := strconv.Atoi(chi.URLParam(request, "item_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	effectiveDate, errB := time.Parse("2006-01-02T15:04:05.000Z", chi.URLParam(request, "effective_date"))
	if errB != nil {
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

	paginatedData, err := r.BomService.GetBomDetailByMasterUn(itemId, effectiveDate, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		paginatedData.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		paginate.Limit,
		paginate.Page,
		int64(paginatedData.TotalRows),
		paginatedData.TotalPages,
	)
}

// @Summary Get Bom Detail By ID
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_detail_id path int true "bom_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/{bom_detail_id} [get]
func (r *BomControllerImpl) GetBomDetailById(writer http.ResponseWriter, request *http.Request) {
	bomId, errA := strconv.Atoi(chi.URLParam(request, "bom_detail_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.BomService.GetBomDetailById(bomId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Bom Detail Max Seq
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_id path int true "bom_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/max-seq/{bom_id} [get]
func (r *BomControllerImpl) GetBomDetailMaxSeq(writer http.ResponseWriter, request *http.Request) {
	bomId, errA := strconv.Atoi(chi.URLParam(request, "bom_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.BomService.GetBomDetailMaxSeq(bomId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @param reqBody body masteritempayloads.BomMasterNewRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/ [post]
func (r *BomControllerImpl) SaveBomMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.BomMasterNewRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	create, err := r.BomService.SaveBomMaster(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusCreated)
}

// @Summary Update Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @Param bom_id path int true "bom_id"
// @Param reqBody body masteritempayloads.BomMasterSaveRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/{bom_id} [put]
func (r *BomControllerImpl) UpdateBomMaster(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.BomMasterSaveRequest
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	if formRequest.Qty < 0 {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("quantity cannot be negative")})
		return
	}

	bomId, errA := strconv.Atoi(chi.URLParam(request, "bom_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	update, err := r.BomService.UpdateBomMaster(bomId, formRequest.Qty)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, update, "Update Data Successfully!", http.StatusOK)
}

// @Summary Save Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param reqBody body masteritempayloads.BomDetailRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail [put]
func (r *BomControllerImpl) SaveBomDetail(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.BomDetailRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	if formRequest.ItemId == formRequest.BomItemId {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("item must be different from parent")})
		return
	}

	create, err := r.BomService.SaveBomDetail(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if create.BomDetailId != 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
	}
}

// @Summary Delete Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_detail_ids path int true "bom_detail_ids"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/{bom_detail_ids} [delete]
func (r *BomControllerImpl) DeleteBomDetail(writer http.ResponseWriter, request *http.Request) {
	var err error
	bomDetailIdsString := chi.URLParam(request, "bom_detail_ids")

	// Convert bomDetailIds to []int
	bomDetailIdsSlice := strings.Split(bomDetailIdsString, ",")
	bomDetailIdsInts := make([]int, len(bomDetailIdsSlice))
	for i := 0; i < len(bomDetailIdsSlice); i++ {
		bomDetailIdsInts[i], err = strconv.Atoi(bomDetailIdsSlice[i])
		if err != nil {
			exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("invalid input")})
			return
		}
	}

	// Call the method to delete bom details by their IDs
	if deleted, err := r.BomService.DeleteByIds(bomDetailIdsInts); err != nil {
		exceptions.NewAppException(writer, request, err)
	} else if deleted {
		payloads.NewHandleSuccess(writer, nil, "Delete Data Successfully!", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Failed to delete data", http.StatusInternalServerError)
	}
}

// DownloadTemplate godoc
// @Summary Download Template
// @Description REST API Download Template
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/download-template [get]
func (r *BomControllerImpl) DownloadTemplate(writer http.ResponseWriter, request *http.Request) {
	// Generate the template file
	f, err := r.BomService.GenerateTemplateFile()
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

	downloadName := time.Now().UTC().Format("2006-01-02_15-04-05") + "_BOMMaster.xlsx"
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
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @Param file formData file true "file"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/upload [post]
func (r *BomControllerImpl) Upload(writer http.ResponseWriter, request *http.Request) {
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

	previewData, errResponse := r.BomService.PreviewUploadData(rows)
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
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @Param reqBody body masteritempayloads.BomDetailUpload true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/process-upload [post]
func (r *BomControllerImpl) ProcessDataUpload(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.BomDetailUpload
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	create, err := r.BomService.ProcessDataUpload(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create/Update Data Successfully!", http.StatusCreated)
}

// GetBomTotalPercentage godoc
// @Summary Get Bom Total Percentage
// @Description REST API Bom Total Percentage
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Security AuthorizationKeyAuth
// @Param bom_id path int true "bom_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/total-percentage/{bom_id} [get]
func (r *BomControllerImpl) GetBomTotalPercentage(writer http.ResponseWriter, request *http.Request) {
	bomId, errA := strconv.Atoi(chi.URLParam(request, "bom_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.BomService.GetBomTotalPercentage(bomId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}
