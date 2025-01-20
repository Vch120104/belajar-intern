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

type ItemTypeController interface {
	GetAllItemType(writer http.ResponseWriter, request *http.Request)
	GetItemTypeById(writer http.ResponseWriter, request *http.Request)
	GetItemTypeByCode(writer http.ResponseWriter, request *http.Request)
	CreateItemType(writer http.ResponseWriter, request *http.Request)
	SaveItemType(writer http.ResponseWriter, request *http.Request)
	ChangeStatusItemType(writer http.ResponseWriter, request *http.Request)
	GetItemTypeDropDown(writer http.ResponseWriter, request *http.Request)
}

type ItemTypeControllerImpl struct {
	ItemTypeService masteritemservice.ItemTypeService
}

func NewItemTypeController(itemTypeService masteritemservice.ItemTypeService) ItemTypeController {
	return &ItemTypeControllerImpl{
		ItemTypeService: itemTypeService,
	}
}

// GetItemTypeDropDown implements ItemTypeController.
func (r *ItemTypeControllerImpl) GetItemTypeDropDown(writer http.ResponseWriter, request *http.Request) {
	response, err := r.ItemTypeService.GetItemTypeDropDown()

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// GetItemTypeByCode implements ItemTypeController.
func (r *ItemTypeControllerImpl) GetItemTypeByCode(writer http.ResponseWriter, request *http.Request) {
	itemTypeCode := chi.URLParam(request, "item_type_code")

	response, err := r.ItemTypeService.GetItemTypeByCode(itemTypeCode)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// GetItemTypeById implements ItemTypeController.
func (r *ItemTypeControllerImpl) GetItemTypeById(writer http.ResponseWriter, request *http.Request) {
	itemTypeId, errA := strconv.Atoi(chi.URLParam(request, "item_type_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemTypeService.GetItemTypeById(itemTypeId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Get Data Successfully!", http.StatusOK)
}

// GetAllItemType implements ItemTypeController.
func (r *ItemTypeControllerImpl) GetAllItemType(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"item_type_id":   queryValues.Get("item_type_id"),
		"item_type_code": queryValues.Get("item_type_code"),
		"item_type_name": queryValues.Get("item_type_name"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result, err := r.ItemTypeService.GetAllItemType(criteria, paginate)

	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(
		writer,
		result.Rows,
		"Get Data Successfully!",
		http.StatusOK,
		result.Limit,
		result.Page,
		int64(result.TotalRows),
		result.TotalPages,
	)
}

// CreateItemType implements ItemTypeController to create new ItemType.
func (r *ItemTypeControllerImpl) CreateItemType(writer http.ResponseWriter, request *http.Request) {
	var formRequest masteritempayloads.ItemTypeRequest

	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemTypeService.CreateItemType(formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message := "Data created successfully!"
	payloads.NewHandleSuccess(writer, create, message, http.StatusCreated)
}

// SaveItemType implements ItemTypeController to update an existing ItemType.
func (r *ItemTypeControllerImpl) SaveItemType(writer http.ResponseWriter, request *http.Request) {
	itemTypeID, errA := strconv.Atoi(chi.URLParam(request, "item_type_id"))
	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	var formRequest masteritempayloads.ItemTypeRequest

	// Parse request body to formRequest
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}

	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.ItemTypeService.SaveItemType(itemTypeID, formRequest)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	message := "Data updated successfully!"
	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// ChangeStatusItemType implements ItemTypeController.
func (r *ItemTypeControllerImpl) ChangeStatusItemType(writer http.ResponseWriter, request *http.Request) {
	itemTypeId, errA := strconv.Atoi(chi.URLParam(request, "item_type_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.ItemTypeService.ChangeStatusItemType(itemTypeId)

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Change Status Data Successfully!", http.StatusOK)
}
