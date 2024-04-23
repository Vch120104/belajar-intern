package masteritemcontroller

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemImportController interface {
	GetAllItemImport(writer http.ResponseWriter, request *http.Request)
	GetItemImportbyId(writer http.ResponseWriter, request *http.Request)
	SaveItemImport(writer http.ResponseWriter, request *http.Request)
	UpdateItemImport(writer http.ResponseWriter, request *http.Request)
}

type ItemImportControllerImpl struct {
	ItemImportService masteritemservice.ItemImportService
}

// GetItemImportbyId implements ItemImportController.
func (r *ItemImportControllerImpl) GetItemImportbyId(writer http.ResponseWriter, request *http.Request) {

	itemPackageId, _ := strconv.Atoi(chi.URLParam(request, "item_import_id"))

	result, err := r.ItemImportService.GetItemImportbyId(itemPackageId)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

// GetAllItemImport implements ItemImportController.
func (r *ItemImportControllerImpl) GetAllItemImport(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	internalFilterCondition := map[string]string{
		"mtr_item.item_code": queryValues.Get("item_code"),
		"mtr_item.item_name": queryValues.Get("item_name"),
	}
	externalFilterCondition := map[string]string{

		"mtr_supplier.supplier_code": queryValues.Get("supplier_code"),
		"mtr_supplier.supplier_name": queryValues.Get("supplier_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internalCriteria := utils.BuildFilterCondition(internalFilterCondition)
	externalCriteria := utils.BuildFilterCondition(externalFilterCondition)

	paginatedData, totalPages, totalRows, err := r.ItemImportService.GetAllItemImport(internalCriteria, externalCriteria, paginate)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, paginate.Limit, paginate.Page, int64(totalRows), totalPages)
}

// SaveItemImport implements ItemImportController.
func (r *ItemImportControllerImpl) SaveItemImport(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritementities.ItemImport
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemImportService.SaveItemImport(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

// UpdateItemImport implements ItemImportController.
func (r *ItemImportControllerImpl) UpdateItemImport(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritementities.ItemImport
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemImportService.UpdateItemImport(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}

func NewItemImportController(ItemImportService masteritemservice.ItemImportService) ItemImportController {
	return &ItemImportControllerImpl{
		ItemImportService: ItemImportService,
	}
}
