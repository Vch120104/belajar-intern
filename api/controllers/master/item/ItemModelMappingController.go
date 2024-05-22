package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemModelMappingController interface {
	CreateItemModelMapping(writer http.ResponseWriter, request *http.Request)
	UpdateItemModelMapping(writer http.ResponseWriter, request *http.Request)
	GetItemModelMappingByItemId(writer http.ResponseWriter, request *http.Request)
}

type ItemModelMappingControllerImpl struct {
	ItemModelMappingService masteritemservice.ItemModelMappingService
}

// GetItemModelMappingByItemId implements ItemModelMappingController.
// @Summary Get Item Model Mapping By Item Id
// @Description REST API to get item model mapping by item id
// @Accept json
// @Produce json
// @Tags Master : Item Model Mapping
// @Param item_id path int true "Item ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-model-mapping/{item_id} [get]
func (r *ItemModelMappingControllerImpl) GetItemModelMappingByItemId(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	itemId, _ := strconv.Atoi(chi.URLParam(request, "item_id"))

	paginate := pagination.Pagination{
		Limit: utils.NewGetQueryInt(queryValues, "limit"),
		Page:  utils.NewGetQueryInt(queryValues, "page"),
	}

	paginatedData, totalPages, totalRows, err := r.ItemModelMappingService.GetItemModelMappingByItemId(itemId, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// UpdateItemModelMapping implements ItemModelMappingController.
// @Summary Update Item Model Mapping
// @Description REST API to update item model mapping
// @Accept json
// @Produce json
// @Tags Master : Item Model Mapping
// @Param reqBody body masteritempayloads.CreateItemModelMapping true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-model-mapping/ [put]
func (r *ItemModelMappingControllerImpl) UpdateItemModelMapping(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.CreateItemModelMapping
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

	create, err := r.ItemModelMappingService.UpdateItemModelMapping(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}

// CreateItemModelMapping implements ItemModelMappingController.
// @Summary Create Item Model Mapping
// @Description REST API to create item model mapping
// @Accept json
// @Produce json
// @Tags Master : Item Model Mapping
// @Param reqBody body masteritempayloads.CreateItemModelMapping true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-model-mapping/ [post]
func (r *ItemModelMappingControllerImpl) CreateItemModelMapping(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.CreateItemModelMapping

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

	create, err := r.ItemModelMappingService.CreateItemModelMapping(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

func NewItemModelMappingController(ItemModelMappingService masteritemservice.ItemModelMappingService) ItemModelMappingController {
	return &ItemModelMappingControllerImpl{
		ItemModelMappingService: ItemModelMappingService,
	}
}
