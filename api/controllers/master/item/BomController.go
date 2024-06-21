package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BomController interface {
	GetBomMasterById(writer http.ResponseWriter, request *http.Request)
	GetBomMasterList(writer http.ResponseWriter, request *http.Request)
	SaveBomMaster(writer http.ResponseWriter, request *http.Request)
	ChangeStatusBomMaster(writer http.ResponseWriter, request *http.Request)
	GetBomDetailList(writer http.ResponseWriter, request *http.Request)
	GetBomDetailById(writer http.ResponseWriter, request *http.Request)
	GetBomDetailByIds(writer http.ResponseWriter, request *http.Request)
	SaveBomDetail(writer http.ResponseWriter, request *http.Request)
	GetBomItemList(writer http.ResponseWriter, request *http.Request)
	DeleteBomDetail(writer http.ResponseWriter, request *http.Request)
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
// @Param bom_master_effective_date query string false "bom_master_effective_date"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/ [get]
func (r *BomControllerImpl) GetBomMasterList(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	// Define query parameters
	queryParams := map[string]string{
		"bom_master_id":             queryValues.Get("bom_master_id"), // Ambil nilai bom_master_id tanpa mtr_bom_master.
		"item_id":                   queryValues.Get("item_id"),
		"bom_master_effective_date": queryValues.Get("bom_master_effective_date"),
		"is_active":                 queryValues.Get("is_active"),
		"bom_master_qty":            queryValues.Get("bom_master_qty"),
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
	paginatedData, totalPages, totalRows, err := r.BomService.GetBomMasterList(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	// Construct the response
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		// If paginatedData is empty, return error response
		exceptions.NewNotFoundException(writer, request, err)
	}
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

	bomMasterId, _ := strconv.Atoi(chi.URLParam(request, "bom_master_id"))

	result, err := r.BomService.GetBomMasterById(bomMasterId)
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
// @param reqBody body masteritempayloads.BomMasterResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/ [put]
func (r *BomControllerImpl) SaveBomMaster(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.BomMasterRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

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

	bomMasterId, _ := strconv.Atoi(chi.URLParam(request, "bom_master_id"))

	entity, err := r.BomService.ChangeStatusBomMaster(int(bomMasterId))
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	responseData := map[string]interface{}{
		"is_active":     entity.IsActive,
		"bom_master_id": entity.BomMasterId,
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

	// Define query parameters
	queryParams := map[string]string{
		"mtr_bom_detail.bom_detail_id": queryValues.Get("bom_detail_id"), // Ambil nilai bom_detail_id tanpa mtr_bom_detail.
		"mtr_bom_detail.bom_master_id": queryValues.Get("bom_master_id"),
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
	paginatedData, totalPages, totalRows, err := r.BomService.GetBomDetailList(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	// Construct the response
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		// If paginatedData is empty, return error response
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
}

// @Summary Get Bom Detail By ID
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_master_id path int true "bom_master_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/{bom_master_id} [get]
func (r *BomControllerImpl) GetBomDetailById(writer http.ResponseWriter, request *http.Request) {

	bomDetailId, _ := strconv.Atoi(chi.URLParam(request, "bom_master_id"))

	result, err := r.BomService.GetBomDetailById(bomDetailId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Bom Detail By ID
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_master_id path int true "bom_master_id"
// @Param bom_detail_id path int true "bom_detail_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/{bom_master_id}/{bom_detail_id} [get]
func (r *BomControllerImpl) GetBomDetailByIds(writer http.ResponseWriter, request *http.Request) {

	bomDetailId, _ := strconv.Atoi(chi.URLParam(request, "bom_detail_id"))

	result, err := r.BomService.GetBomDetailByIds(bomDetailId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	if len(result) > 0 {
		payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
	} else {
		payloads.NewHandleError(writer, "Data not found", http.StatusNotFound)
	}
}

// @Summary Update Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param bom_master_id path int true "bom_master_id"
// @Param bom_detail_id path int true "bom_detail_id"
// @Param reqBody body masteritempayloads.BomDetailResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/bom/detail/{bom_master_id}/{bom_detail_id} [put]
func (r *BomControllerImpl) SaveBomDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.BomDetailRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

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
	paginatedData, totalPages, totalRows, err := r.BomService.GetBomItemList(criteria, paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	// Construct the response
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		// If paginatedData is empty, return error response
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
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
