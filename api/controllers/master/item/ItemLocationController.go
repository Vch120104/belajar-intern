package masteritemcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
)

type ItemLocationController interface {
	GetAllItemLocation(writer http.ResponseWriter, request *http.Request)
	SaveItemLocation(writer http.ResponseWriter, request *http.Request)
}

type ItemLocationControllerImpl struct {
	ItemLocationService masteritemservice.ItemLocationService
}

func NewItemLocationController(ItemLocationService masteritemservice.ItemLocationService) ItemLocationController {
	return &ItemLocationControllerImpl{
		ItemLocationService: ItemLocationService,
	}
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
		"mtr_item_location.item_location_id":     queryValues.Get("item_location_id"),
		"mtr_item_location.warehouse_group_id":   queryValues.Get("warehouse_group_id"),
		"mtr_item_location.warehouse_group_name": queryValues.Get("warehouse_group_name"),
		"mtr_item_location.item_id":              queryValues.Get("item_id"),
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

// @Summary Get All Bom Detail
// @Description REST API Bom Detail
// @Accept json
// @Produce json
// @Tags Master : Bom Detail
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param bom_detail_id query string false "bom_detail_id"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /item-location/all/detail [get]
// func (r *ItemLocationControllerImpl) GetAllItemLocationDetail(writer http.ResponseWriter, request *http.Request) {
// 	queryValues := request.URL.Query()

// 	// Define query parameters
// 	queryParams := map[string]string{
// 		"mtr_item_location_detail.item_location_detail_id": queryValues.Get("item_location_detail_id"), // Ambil nilai bom_detail_id tanpa mtr_bom_detail.
// 		"mtr_item_location_detail.item_location_id":        queryValues.Get("item_location_id"),
// 	}

// 	// Extract pagination parameters
// 	paginate := pagination.Pagination{
// 		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
// 		Page:   utils.NewGetQueryInt(queryValues, "page"),
// 		SortOf: queryValues.Get("sort_of"),
// 		SortBy: queryValues.Get("sort_by"),
// 	}

// 	// Build filter condition based on query parameters
// 	criteria := utils.BuildFilterCondition(queryParams)

// 	// Call service to get paginated data
// 	paginatedData, totalPages, totalRows, err := r.ItemLocationService.GetAllItemLocationDetail(criteria, paginate)

// 	// Construct the response
// 	if len(paginatedData) > 0 {
// 		payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "Get Data Successfully", http.StatusOK, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
// 	} else {
// 		// If paginatedData is empty, return error response
// 		exceptionsss_test.NewNotFoundException(writer, request, err)
// 		return
// 	}
// }
