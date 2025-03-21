package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
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

type ItemPackageDetailController interface {
	GetItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request)
	GetItemPackageDetailById(writer http.ResponseWriter, request *http.Request)
	CreateItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request)
	UpdateItemPackageDetail(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemPackageDetail(writer http.ResponseWriter, request *http.Request)
	ActivateItemPackageDetail(writer http.ResponseWriter, request *http.Request)
	DeactivateItemPackageDetail(writer http.ResponseWriter, request *http.Request)
}

type ItemPackageDetailControllerImpl struct {
	ItemPackageDetailService masteritemservice.ItemPackageDetailService
}

func NewItemPackageDetailController(ItemPackageDetailService masteritemservice.ItemPackageDetailService) ItemPackageDetailController {
	return &ItemPackageDetailControllerImpl{
		ItemPackageDetailService: ItemPackageDetailService,
	}
}

// @Summary Activate Item Package Detail
// @Description Activate an item package detail by its ID
// @Accept json
// @Produce json
// @Tags Master : Item Package Detail
// @Security AuthorizationKeyAuth
// @Param item_package_detail_id path int true "Item Package Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package-detail/activate/{item_package_detail_id} [patch]
func (r *ItemPackageDetailControllerImpl) ActivateItemPackageDetail(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "item_package_detail_id")

	response, err := r.ItemPackageDetailService.ActivateItemPackageDetail(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Activate Status Successfully!", http.StatusOK)
}

// @Summary Deactivate Item Package Detail
// @Description Deactivate an item package detail by its ID
// @Accept json
// @Produce json
// @Tags Master : Item Package Detail
// @Security AuthorizationKeyAuth
// @Param item_package_detail_id path int true "Item Package Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package-detail/deactivate/{item_package_detail_id} [patch]
func (r *ItemPackageDetailControllerImpl) DeactivateItemPackageDetail(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "item_package_detail_id")

	response, err := r.ItemPackageDetailService.DeactiveItemPackageDetail(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Deactivate Status Successfully!", http.StatusOK)
}

// @Summary Change Status Item Package Detail
// @Description Change the status of an item package detail by its ID
// @Accept json
// @Produce json
// @Tags Master : Item Package Detail
// @Security AuthorizationKeyAuth
// @Param item_package_detail_id path int true "Item Package Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package-detail/{item_package_detail_id} [patch]
func (r *ItemPackageDetailControllerImpl) ChangeStatusItemPackageDetail(writer http.ResponseWriter, request *http.Request) {
	id, errA := strconv.Atoi(chi.URLParam(request, "item_package_detail_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemPackageDetailService.ChangeStatusItemPackageDetail(id)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

// @Summary Get Item Package Detail By Item Package ID
// @Description Retrieve all item package details by their package ID with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Master : Item Package Detail
// @Security AuthorizationKeyAuth
// @Param item_package_id path int true "Item Package ID"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package-detail/by-package-id/{item_package_id} [get]
func (r *ItemPackageDetailControllerImpl) GetItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	itemPackageId, errA := strconv.Atoi(chi.URLParam(request, "item_package_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.ItemPackageDetailService.GetItemPackageDetailByItemPackageId(itemPackageId, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Item Package Detail By ID
// @Description Retrieve an item package detail by its ID
// @Accept json
// @Produce json
// @Tags Master : Item Package Detail
// @Security AuthorizationKeyAuth
// @Param item_package_detail_id path int true "Item Package Detail ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package-detail/{item_package_detail_id} [get]
func (r *ItemPackageDetailControllerImpl) GetItemPackageDetailById(writer http.ResponseWriter, request *http.Request) {

	itemPackageId, errA := strconv.Atoi(chi.URLParam(request, "item_package_detail_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.ItemPackageDetailService.GetItemPackageDetailById(itemPackageId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Create Item Package Detail By Item Package ID
// @Description Create a new item package detail under a specific package ID
// @Accept json
// @Produce json
// @Tags Master : Item Package Detail
// @Security AuthorizationKeyAuth
// @Param reqBody body masteritempayloads.SaveItemPackageDetail true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package-detail [post]
func (r *ItemPackageDetailControllerImpl) CreateItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.SaveItemPackageDetail
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

	create, err := r.ItemPackageDetailService.CreateItemPackageDetailByItemPackageId(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// @Summary Update Item Package Detail By Item Package ID
// @Description Update an existing item package detail under a specific package ID
// @Accept json
// @Produce json
// @Tags Master : Item Package Detail
// @Security AuthorizationKeyAuth
// @Param reqBody body masteritempayloads.SaveItemPackageDetail true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package-detail [put]
func (r *ItemPackageDetailControllerImpl) UpdateItemPackageDetail(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.SaveItemPackageDetail
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

	create, err := r.ItemPackageDetailService.UpdateItemPackageDetail(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}
