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
	"time"

	"github.com/go-chi/chi/v5"
)

type BomController interface {
	// Parent
	GetBomMasterList(writer http.ResponseWriter, request *http.Request)

	GetBomMasterById(writer http.ResponseWriter, request *http.Request)
	SaveBomMaster(writer http.ResponseWriter, request *http.Request)
	UpdateBomMaster(writer http.ResponseWriter, request *http.Request)
	ChangeStatusBomMaster(writer http.ResponseWriter, request *http.Request)

	// Detail
	GetBomDetailList(writer http.ResponseWriter, request *http.Request)
	GetBomDetailById(writer http.ResponseWriter, request *http.Request)
	SaveBomDetail(writer http.ResponseWriter, request *http.Request)
	UpdateBomDetail(writer http.ResponseWriter, request *http.Request)
	DeleteBomDetail(writer http.ResponseWriter, request *http.Request)

	GetBomItemList(writer http.ResponseWriter, request *http.Request)
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
func (r *BomControllerImpl) GetBomMasterList(writer http.ResponseWriter, request *http.Request) {
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

	paginatedData, err := r.BomService.GetBomMasterList(criteria, paginate)
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
// @Param bom_master_id path int true "bom_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/{bom_master_id} [get]
func (r *BomControllerImpl) GetBomMasterById(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	bomMasterId, errA := strconv.Atoi(chi.URLParam(request, "bom_master_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	// Extract pagination parametersF
	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.BomService.GetBomMasterById(bomMasterId, paginate)
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
// @param reqBody body masteritempayloads.BomMasterRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/ [post]
func (r *BomControllerImpl) SaveBomMaster(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.BomMasterRequest
	var message = ""
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
	if formRequest.BomMasterId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
	}

}

// @Summary Update Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @Param bom_master_id path int true "bom_master_id"
// @Param reqBody body masteritempayloads.BomMasterRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/{bom_master_id} [put]
func (r *BomControllerImpl) UpdateBomMaster(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.BomMasterRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	bomMasterId, errA := strconv.Atoi(chi.URLParam(request, "bom_master_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	update, err := r.BomService.UpdateBomMaster(bomMasterId, formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if formRequest.BomMasterId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, update, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, update, message, http.StatusOK)
	}

}

// @Summary Change Status Bom Master
// @Description REST API Bom Master
// @Accept json
// @Produce json
// @Tags Master : Bom Master
// @param bom_master_id path int true "bom_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/{bom_master_id} [patch]
func (r *BomControllerImpl) ChangeStatusBomMaster(writer http.ResponseWriter, request *http.Request) {

	bomMasterId, errA := strconv.Atoi(chi.URLParam(request, "bom_master_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	entity, err := r.BomService.ChangeStatusBomMaster(int(bomMasterId))
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

// @Summary Get All Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param bom_detail_id query string false "bom_detail_id"
// @Param bom_master_id query string false "bom_master_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail [get]
func (r *BomControllerImpl) GetBomDetailList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_bom_detail.bom_detail_id":              queryValues.Get("bom_detail_id"), // Ambil nilai bom_detail_id tanpa mtr_bom_detail.
		"mtr_bom_detail.bom_master_id":              queryValues.Get("bom_master_id"),
		"mtr_bom_detail.bom_detail_qty":             queryValues.Get("bom_detail_qty"),
		"mtr_bom_detail.bom_detail_remark":          queryValues.Get("bom_detail_remark"),
		"mtr_bom_detail.bom_detail_costing_percent": queryValues.Get("bom_detail_costing_percent"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.BomService.GetBomDetailList(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
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

	bomMasterId, errA := strconv.Atoi(chi.URLParam(request, "bom_detail_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	queryParams := map[string]string{
		"bom_detail_id": chi.URLParam(request, "bom_detail_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(request.URL.Query(), "limit"),
		Page:   utils.NewGetQueryInt(request.URL.Query(), "page"),
		SortOf: request.URL.Query().Get("sort_of"),
		SortBy: request.URL.Query().Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, totalPages, totalRows, err := r.BomService.GetBomDetailById(bomMasterId, criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {

		exceptions.NewNotFoundException(writer, request, err)
	}

}

// @Summary Save Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param reqBody body masteritempayloads.BomDetailRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail [post]
func (r *BomControllerImpl) SaveBomDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.BomDetailRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}

	create, err := r.BomService.SaveBomDetail(formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	if formRequest.BomDetailId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
	}

}

// @Summary Update Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_master_id path int true "bom_master_id"
// @Param bom_detail_id path int true "bom_detail_id"
// @Param reqBody body masteritempayloads.BomDetailRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/{bom_master_id}/{bom_detail_id} [put]
func (r *BomControllerImpl) UpdateBomDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.BomDetailRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)
	if validationErr := validation.ValidationForm(writer, request, &formRequest); validationErr != nil {
		exceptions.NewBadRequestException(writer, request, validationErr)
		return
	}
	bomDetailId, errA := strconv.Atoi(chi.URLParam(request, "bom_detail_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	update, err := r.BomService.UpdateBomDetail(bomDetailId, formRequest)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if formRequest.BomDetailId == 0 {
		message = "Create Data Successfully!"
		payloads.NewHandleSuccess(writer, update, message, http.StatusCreated)
	} else {
		message = "Update Data Successfully!"
		payloads.NewHandleSuccess(writer, update, message, http.StatusOK)
	}

}

// @Summary Get All Bom Item Lookup
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param item_code query string false "item_code"
// @Param item_name query string false "item_name"
// @Param item_type query string false "item_type"
// @Param item_group_code query string false "item_group_code"
// @Param item_class_code query string false "item_class_code"
// @Param uom_code query string false "uom_code"
// @Param is_active query string false "is_active"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/popup-item [get]
func (r *BomControllerImpl) GetBomItemList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"item_code":       queryValues.Get("item_code"),
		"item_name":       queryValues.Get("item_name"),
		"item_type":       queryValues.Get("item_type"),
		"item_group_code": queryValues.Get("item_group_code"),
		"item_class_code": queryValues.Get("item_class_code"),
		"uom_name":        queryValues.Get("uom_name"),
		"is_active":       queryValues.Get("is_active"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	paginatedData, err := r.BomService.GetBomItemList(criteria, paginate)
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

// @Summary Delete Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_master_id path int true "bom_master_id"
// @Param bom_detail_id path int true "bom_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/{bom_master_id}/{bom_detail_id} [delete]
func (r *BomControllerImpl) DeleteBomDetail(writer http.ResponseWriter, request *http.Request) {

	bomDetailID := chi.URLParam(request, "bom_detail_id")

	// Ubah bomDetailID menjadi integer
	bomDetailIDInt, err := strconv.Atoi(bomDetailID)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{
			Err: errors.New("invalid bom_detail_id"),
		})
		return
	}

	// Call the method to delete bom details by their IDs
	if deleted, err := r.BomService.DeleteByIds([]int{bomDetailIDInt}); err != nil {
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
