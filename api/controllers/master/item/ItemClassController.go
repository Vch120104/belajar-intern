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

type ItemClassController interface {
	GetAllItemClass(writer http.ResponseWriter, request *http.Request)
	GetItemClassDropdown(writer http.ResponseWriter, request *http.Request)
	GetItemClassbyId(writer http.ResponseWriter, request *http.Request)
	GetItemClassByCode(writer http.ResponseWriter, request *http.Request)
	SaveItemClass(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemClass(writer http.ResponseWriter, request *http.Request)
	GetItemClassDropDownbyGroupId(writer http.ResponseWriter, request *http.Request)
	GetItemClassMfgDropdown(writer http.ResponseWriter, request *http.Request)
}
type ItemClassControllerImpl struct {
	ItemClassService masteritemservice.ItemClassService
}

func NewItemClassController(itemClassService masteritemservice.ItemClassService) ItemClassController {
	return &ItemClassControllerImpl{
		ItemClassService: itemClassService,
	}
}

// @Summary Get ItemClass DropDownbyGroupId
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master Item : Item Class
// @Param item_group_id path int true "item_group_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class/drop-down/by-group-id/{item_group_id} [get]
func (r *ItemClassControllerImpl) GetItemClassDropDownbyGroupId(writer http.ResponseWriter, request *http.Request) {
	itemGroupId, errA := strconv.Atoi(chi.URLParam(request, "item_group_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemClassService.GetItemClassDropDownbyGroupId(itemGroupId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Item Class By Code
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master Item : Item Class
// @Param item_class_code path string true "item_class_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class/by-code/{item_class_code} [get]
func (r *ItemClassControllerImpl) GetItemClassByCode(writer http.ResponseWriter, request *http.Request) {
	itemClassCode := chi.URLParam(request, "item_class_code")

	response, err := r.ItemClassService.GetItemClassByCode(itemClassCode)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get Item Class By ID
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master Item : Item Class
// @Param item_class_id path int true "item_class_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class/{item_class_id} [get]
func (r *ItemClassControllerImpl) GetItemClassbyId(writer http.ResponseWriter, request *http.Request) {
	itemClassId, errA := strconv.Atoi(chi.URLParam(request, "item_class_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemClassService.GetItemClassById(itemClassId)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// @Summary Get All Item Class Lookup
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master Item : Item Class
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
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class [get]
func (r *ItemClassControllerImpl) GetAllItemClass(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	internalFilter := map[string]string{
		"mtr_item_class.is_active":       queryValues.Get("is_active"),
		"mtr_item_class.item_class_id":   queryValues.Get("item_class_id"),
		"mtr_item_class.item_class_code": queryValues.Get("item_class_code"),
		"mtr_item_class.item_class_name": queryValues.Get("item_class_name"),
	}

	externalFilter := map[string]string{
		"item_group_name": queryValues.Get("item_group_name"),
		"line_type_code":  queryValues.Get("line_type_code"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	internal := utils.BuildFilterCondition(internalFilter)
	external := utils.BuildFilterCondition(externalFilter)

	result, err := r.ItemClassService.GetAllItemClass(internal, external, pagination)

	response := utils.ModifyKeysInResponse(result.Rows)

	if err != nil {
		response = []interface{}{}
		result.TotalPages = 0
		result.TotalRows = 0
	}
	payloads.NewHandleSuccessPagination(writer, response, "Get Data Successfully!", http.StatusOK, pagination.Limit, pagination.Page, result.TotalRows, result.TotalPages)

}

// @Summary Get All Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master Item : Item Class
// @Param is_active query string false "is_active" Enums(true, false)
// @Param item_class_id query int false "item_class_id"
// @Param item_class_code query string false "item_class_code"
// @Param item_class_name query string false "item_class_name"
// @Param item_group_name query string false "item_group_name"
// @Param line_type_code query string false "line_type_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class/drop-down [get]
func (r *ItemClassControllerImpl) GetItemClassDropdown(writer http.ResponseWriter, request *http.Request) {
	result, err := r.ItemClassService.GetItemClassDropDown()

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Get Item Class Mfg Dropdown
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master Item : Item Class
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class/mfg/drop-down [get]
func (r *ItemClassControllerImpl) GetItemClassMfgDropdown(writer http.ResponseWriter, request *http.Request) {
	result, err := r.ItemClassService.GetItemClassMfgDropdown()
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "success", 200)
}

// @Summary Save Item Class
// @Description REST API Item Class
// @Accept json
// @Produce json
// @Tags Master Item : Item Class
// @param reqBody body masteritempayloads.ItemClassResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class [post]
func (r *ItemClassControllerImpl) SaveItemClass(writer http.ResponseWriter, request *http.Request) {

	var formRequest masteritempayloads.ItemClassResponse
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message = ""

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemClassService.SaveItemClass(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
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
// @Tags Master Item : Item Class
// @param item_class_id path int true "item_class_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/item-class/{item_class_id} [patch]
func (r *ItemClassControllerImpl) ChangeStatusItemClass(writer http.ResponseWriter, request *http.Request) {

	itemClassId, errA := strconv.Atoi(chi.URLParam(request, "item_class_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemClassService.ChangeStatusItemClass(int(itemClassId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
