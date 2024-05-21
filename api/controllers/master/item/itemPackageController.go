package masteritemcontroller

import (
	exceptions "after-sales/api/exceptions"
	"errors"

	helper "after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemPackageController interface {
	GetAllItemPackage(writer http.ResponseWriter, request *http.Request)
	SaveItemPackage(writer http.ResponseWriter, request *http.Request)
	GetItemPackageById(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemPackage(writer http.ResponseWriter, request *http.Request)
}

type ItemPackageControllerImpl struct {
	ItemPackageService masteritemservice.ItemPackageService
}

func NewItemPackageController(ItemPackageService masteritemservice.ItemPackageService) ItemPackageController {
	return &ItemPackageControllerImpl{
		ItemPackageService: ItemPackageService,
	}
}

// @Summary Get All Item Packages
// @Description Retrieve all item packages with optional filtering and pagination
// @Accept json
// @Produce json
// @Tags Master : Item Package
// @Param item_package_code query string false "Item Package Code"
// @Param item_package_name query string false "Item Package Name"
// @Param item_package_set query string false "Item Package Set"
// @Param description query string false "Description"
// @Param is_active query string false "Is Active"
// @Param item_group_code query string false "Item Group Code"
// @Param page query string true "Page number"
// @Param limit query string true "Items per page"
// @Param sort_by query string false "Field to sort by"
// @Param sort_of query string false "Sort order (asc/desc)"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package/ [get]
func (r *ItemPackageControllerImpl) GetAllItemPackage(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	internalFilterCondition := map[string]string{
		"mtr_item_package.item_package_code": queryValues.Get("item_package_code"),
		"mtr_item_package.item_package_name": queryValues.Get("item_package_name"),
		"mtr_item_package.item_package_set":  queryValues.Get("item_package_set"),
		"mtr_item_package.description":       queryValues.Get("description"),
		"mtr_item_package.is_active":         queryValues.Get("is_active"),
	}
	externalFilterCondition := map[string]string{

		"mtr_item_group.item_group_code": queryValues.Get("item_group_code"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginatedData, totalPages, totalRows, err := r.ItemPackageService.GetAllItemPackage(internalCriteria, externalCriteria, paginate)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// @Summary Get Item Package By ID
// @Description Retrieve an item package by its ID
// @Accept json
// @Produce json
// @Tags Master : Item Package
// @Param item_package_id path int true "Item Package ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package/{item_package_id} [get]
func (r *ItemPackageControllerImpl) GetItemPackageById(writer http.ResponseWriter, request *http.Request) {

	itemPackageId, _ := strconv.Atoi(chi.URLParam(request, "item_package_id"))

	result, err := r.ItemPackageService.GetItemPackageById(itemPackageId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Item Package
// @Description Create or update an item package
// @Accept json
// @Produce json
// @Tags Master : Item Package
// @Param reqBody body masteritempayloads.SaveItemPackageRequest true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package/ [post]
func (r *ItemPackageControllerImpl) SaveItemPackage(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.SaveItemPackageRequest
	var message string
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, errors.New("invalid form request"))
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, errors.New("invalid format request"))
		return
	}

	create, err := r.ItemPackageService.SaveItemPackage(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.ItemPackageId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Item Package
// @Description Change the status of an item package by its ID
// @Accept json
// @Produce json
// @Tags Master : Item Package
// @Param item_package_id path int true "Item Package ID"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-package/{item_package_id} [patch]
func (r *ItemPackageControllerImpl) ChangeStatusItemPackage(writer http.ResponseWriter, request *http.Request) {

	PriceListId, _ := strconv.Atoi(chi.URLParam(request, "item_package_id"))

	response, err := r.ItemPackageService.ChangeStatusItemPackage(PriceListId)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}
