package masteritemcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	helper_test "after-sales/api/helper_testt"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemLocationController interface {
	GetAllItemLocation(writer http.ResponseWriter, request *http.Request)
	SaveItemLocation(writer http.ResponseWriter, request *http.Request)
	GetItemLocationById(writer http.ResponseWriter, request *http.Request)
	GetAllItemLocationDetail(writer http.ResponseWriter, request *http.Request)
	PopupItemLocation(writer http.ResponseWriter, request *http.Request)
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
// @Tags Master : Item Location Popup
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /popup-location [get]
func (r *ItemLocationControllerImpl) PopupItemLocation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
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
		exceptionsss_test.NewNotFoundException(writer, request, err)
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
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /item-location [get]
func (r *ItemLocationControllerImpl) GetAllItemLocation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_item_location.warehouse_group_id": queryValues.Get("warehouse_group_id"),
		"mtr_item_location.item_id":            queryValues.Get("item_id"),
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
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Save Item Location
// @Description REST API Item Location
// @Accept json
// @Produce json
// @Tags Master :Item Location
// @param reqBody body masteritempayloads.ItemLocationResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router / [put]
func (r *ItemLocationControllerImpl) SaveItemLocation(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.ItemLocationRequest
	var message = ""
	helper.ReadFromRequestBody(request, &formRequest)

	create, err := r.ItemLocationService.SaveItemLocation(formRequest)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}
	if formRequest.ItemLocationId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Get Item Location By ID
// @Description REST API  Item Location
// @Accept json
// @Produce json
// @Tags Master :  Item Location
// @Param item_location_id path int true "item_location_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /{item_location_id} [get]
func (r *ItemLocationControllerImpl) GetItemLocationById(writer http.ResponseWriter, request *http.Request) {

	IncentiveLevelIds, _ := strconv.Atoi(chi.URLParam(request, "item_location_id"))

	result, err := r.ItemLocationService.GetItemLocationById(IncentiveLevelIds)
	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Item Location Detail
// @Description REST API Item Location Detail
// @Accept json
// @Produce json
// @Tags Master : Item Location Detail
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param item_location_detail_id query string false "item_location_detail_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /{item_location_id}/detail [get]
func (r *ItemLocationControllerImpl) GetAllItemLocationDetail(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	// Define query parameters
	queryParams := map[string]string{
		"mtr_item_location_detail.item_location_detail_id": queryValues.Get("item_location_detail_id"), // Ambil nilai bom_detail_id tanpa mtr_bom_detail.
		"mtr_item_location_detail.item_location_id":        queryValues.Get("item_location_id"),
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

	// Construct the response
	if len(paginatedData) > 0 {
		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
	} else {
		// If paginatedData is empty, return error response
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}
}
