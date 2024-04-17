package masteritemcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
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

type ItemPackageDetailController interface {
	GetItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request)
	GetItemPackageDetailById(writer http.ResponseWriter, request *http.Request)
	CreateItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request)
	UpdateItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemPackageDetail(writer http.ResponseWriter, request *http.Request)
}

type ItemPackageDetailControllerImpl struct {
	ItemPackageDetailService masteritemservice.ItemPackageDetailService
}

func NewItemPackageDetailController(ItemPackageDetailService masteritemservice.ItemPackageDetailService) ItemPackageDetailController {
	return &ItemPackageDetailControllerImpl{
		ItemPackageDetailService: ItemPackageDetailService,
	}
}

func (r *ItemPackageDetailControllerImpl) ChangeStatusItemPackageDetail(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "item_package_detail_id"))

	response, err := r.ItemPackageDetailService.ChangeStatusItemPackageDetail(id)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Successfully!", http.StatusOK)
}

func (r *ItemPackageDetailControllerImpl) GetItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	itemPackageId, _ := strconv.Atoi(chi.URLParam(request, "item_package_id"))

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	result, err := r.ItemPackageDetailService.GetItemPackageDetailByItemPackageId(itemPackageId, paginate)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemPackageDetailControllerImpl) GetItemPackageDetailById(writer http.ResponseWriter, request *http.Request) {

	itemPackageId, _ := strconv.Atoi(chi.URLParam(request, "item_package_detail_id"))

	result, err := r.ItemPackageDetailService.GetItemPackageDetailById(itemPackageId)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemPackageDetailControllerImpl) CreateItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.SaveItemPackageDetail
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

	create, err := r.ItemPackageDetailService.CreateItemPackageDetailByItemPackageId(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Create Data Successfully!", http.StatusOK)
}

func (r *ItemPackageDetailControllerImpl) UpdateItemPackageDetailByItemPackageId(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.SaveItemPackageDetail
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

	create, err := r.ItemPackageDetailService.UpdateItemPackageDetailByItemPackageId(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "Update Data Successfully!", http.StatusOK)
}
