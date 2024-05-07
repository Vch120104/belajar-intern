package masteritemcontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemClassController interface {
	GetAllItemClassLookup(writer http.ResponseWriter, request *http.Request)
	GetAllItemClass(writer http.ResponseWriter, request *http.Request)
	SaveItemClass(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemClass(writer http.ResponseWriter, request *http.Request)
}
type ItemClassControllerImpl struct {
	ItemClassService masteritemservice.ItemClassService
}

func NewItemClassController(itemClassService masteritemservice.ItemClassService) ItemClassController {
	return &ItemClassControllerImpl{
		ItemClassService: itemClassService,
	}
}

// @Summary Get All Item Class Lookup
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param item_class_id query int false "item_class_id"
// @Param item_class_code query string false "item_class_code"
// @Param item_class_name query string false "item_class_name"
// @Param item_group_name query string false "item_group_name"
// @Param line_type_code query string false "line_type_code"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item-class/pop-up [get]
func (r *ItemClassControllerImpl) GetAllItemClassLookup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_item_class.is_active":       queryValues.Get("is_active"),
		"mtr_item_class.item_class_id":   queryValues.Get("item_class_id"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item_class.item_class_name": queryValues.Get("item_class_name"),
		"item_group_name":                queryValues.Get("item_group_name"),
		"line_type_code":                 queryValues.Get("line_type_code"),
	}

	limit := utils.NewGetQueryInt(queryValues, "limit")
	page := utils.NewGetQueryInt(queryValues, "page")
	sortOf := queryValues.Get("sort_of")
	sortBy := queryValues.Get("sort_by")

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.ItemClassService.GetAllItemClass(criteria)

	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	paginatedData, totalPages, totalRows := utils.DataFramePaginate(result, page, limit, utils.SnaketoPascalCase(sortOf), sortBy)

	payloads.NewHandleSuccessPagination(writer, utils.ModifyKeysInResponse(paginatedData), "success", 200, limit, page, int64(totalRows), totalPages)
}

// @Summary Get All Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @Param is_active query string false "is_active" Enums(true, false)
// @Param item_class_id query int false "item_class_id"
// @Param item_class_code query string false "item_class_code"
// @Param item_class_name query string false "item_class_name"
// @Param item_group_name query string false "item_group_name"
// @Param line_type_code query string false "line_type_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item-class/ [get]
func (r *ItemClassControllerImpl) GetAllItemClass(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"mtr_item_class.is_active":       queryValues.Get("is_active"),
		"mtr_item_class.item_class_id":   queryValues.Get("item_class_id"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item_class.item_class_name": queryValues.Get("item_class_name"),
		"item_group_name":                queryValues.Get("item_group_name"),
		"line_type_code":                 queryValues.Get("line_type_code"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.ItemClassService.GetAllItemClass(criteria)

	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "success", 200)
}

// @Summary Save Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @param reqBody body masteritempayloads.ItemClassResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item-class/ [post]
func (r *ItemClassControllerImpl) SaveItemClass(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.ItemClassResponse
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message = ""

	if err != nil {
		exceptionsss_test.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemClassService.SaveItemClass(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.ItemClassId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master : Item Class
// @param item_class_id path int true "item_class_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptionsss_test.BaseErrorResponse
// @Router /v1/item-class/{item_class_id} [patch]
func (r *ItemClassControllerImpl) ChangeStatusItemClass(writer http.ResponseWriter, request *http.Request) {

	itemclassGroupId, _ := strconv.Atoi(chi.URLParam(request, "item_class_id"))

	response, err := r.ItemClassService.ChangeStatusItemClass(int(itemclassGroupId))

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
