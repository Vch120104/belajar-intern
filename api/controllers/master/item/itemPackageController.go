package masteritemcontroller

import (
	"after-sales/api/helper"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemPackageController interface {
	GetAllItemPackage(writer http.ResponseWriter, request *http.Request)
	SaveItemPackage(writer http.ResponseWriter, request *http.Request)
	GetItemPackageById(writer http.ResponseWriter, request *http.Request)
}

type ItemPackageControllerImpl struct {
	ItemPackageService masteritemservice.ItemPackageService
}

func NewItemPackageController(ItemPackageService masteritemservice.ItemPackageService) ItemPackageController {
	return &ItemPackageControllerImpl{
		ItemPackageService: ItemPackageService,
	}
}

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

	paginatedData, totalPages, totalRows := r.ItemPackageService.GetAllItemPackage(internalCriteria, externalCriteria, paginate)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

func (r *ItemPackageControllerImpl) GetItemPackageById(writer http.ResponseWriter, request *http.Request) {

	itemPackageId, _ := strconv.Atoi(chi.URLParam(request, "item_package_id"))

	result := r.ItemPackageService.GetItemPackageById(itemPackageId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemPackageControllerImpl) SaveItemPackage(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.SaveItemPackageRequest
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.ItemPackageService.SaveItemPackage(formRequest)

	if formRequest.ItemPackageId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}
