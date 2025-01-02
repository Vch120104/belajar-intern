package masteritemcontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemGroupController interface {
	GetAllItemGroupWithPagination(writer http.ResponseWriter, request *http.Request)
	GetAllItemGroup(writer http.ResponseWriter, request *http.Request)
	GetItemGroupById(writer http.ResponseWriter, request *http.Request)
	GetItemGroupByCode(writer http.ResponseWriter, request *http.Request)
	DeleteItemGroupById(writer http.ResponseWriter, request *http.Request)
	UpdateItemGroupById(writer http.ResponseWriter, request *http.Request)
	UpdateStatusItemGroupById(writer http.ResponseWriter, request *http.Request)
	GetItemGroupByMultiId(writer http.ResponseWriter, request *http.Request)
	NewItemGroup(writer http.ResponseWriter, request *http.Request)
}

type ItemGroupControllerImpl struct {
	service masteritemservice.ItemGroupService
}

func (i *ItemGroupControllerImpl) GetAllItemGroupWithPagination(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	queryParams := map[string]string{
		"item_group_code":   queryValues.Get("item_group_code"),
		"item_group_name":   queryValues.Get("item_group_name"),
		"is_active":         queryValues.Get("is_active"),
		"is_item_sparepart": queryValues.Get("is_item_sparepart"),
	}
	pages := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}
	filter := utils.BuildFilterCondition(queryParams)
	result, err := i.service.GetAllItemGroupWithPagination(filter, pages)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfull", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)

}

func (i *ItemGroupControllerImpl) GetAllItemGroup(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	ItemGroupCode := queryValues.Get("item_group_code")
	res, err := i.service.GetAllItemGroup(ItemGroupCode)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Get Data Successfull", 200)
}

func (i *ItemGroupControllerImpl) GetItemGroupById(writer http.ResponseWriter, request *http.Request) {
	Id := chi.URLParam(request, "item_group_id")
	itemGroupId, errs := strconv.Atoi(Id)
	if errs != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errs,
			Message:    "failed to convert id to int please check input",
		})
		return

	}
	res, err := i.service.GetItemGroupById(itemGroupId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Get Data Successfull", 200)
}

func (i *ItemGroupControllerImpl) DeleteItemGroupById(writer http.ResponseWriter, request *http.Request) {
	Id := chi.URLParam(request, "item_group_id")
	itemGroupId, errs := strconv.Atoi(Id)
	if errs != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errs,
			Message:    "failed to convert id to int please check input",
		})
		return

	}
	res, err := i.service.DeleteItemGroupById(itemGroupId)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}
	payloads.NewHandleSuccess(writer, res, "delete data Successfull", 200)
}

func (i *ItemGroupControllerImpl) UpdateItemGroupById(writer http.ResponseWriter, request *http.Request) {
	payload := masteritempayloads.ItemGroupUpdatePayload{}
	id := chi.URLParam(request, "item_group_id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
			Message:    "failed to convert id to int please check input",
		})
	}

	errCheck := jsonchecker.ReadFromRequestBody(request, &payload)
	if errCheck != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        errCheck.Err,
			Message:    "failed to read request body please check input",
		})
	}
	res, errData := i.service.UpdateItemGroupById(payload, ids)
	if errData != nil {
		helper.ReturnError(writer, request, errData)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Update Data Successfull", 200)

}

func (i *ItemGroupControllerImpl) UpdateStatusItemGroupById(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "item_group_id")
	ids, err := strconv.Atoi(id)
	if err != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        err,
			Message:    "failed to convert id to int please check input",
		})
		return

	}
	res, errData := i.service.UpdateStatusItemGroupById(ids)
	if errData != nil {
		helper.ReturnError(writer, request, errData)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Update Data Successfully", 200)
}

func (i *ItemGroupControllerImpl) GetItemGroupByMultiId(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "item_group_id")
	res, errData := i.service.GetItemGroupByMultiId(id)
	if errData != nil {
		helper.ReturnError(writer, request, errData)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Update Data Successfully", 200)
}

func (i *ItemGroupControllerImpl) NewItemGroup(writer http.ResponseWriter, request *http.Request) {
	payload := masteritempayloads.NewItemGroupPayload{}

	errCheck := jsonchecker.ReadFromRequestBody(request, &payload)
	if errCheck != nil {
		helper.ReturnError(writer, request, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Err:        errCheck.Err,
			Message:    "failed to read request body please check input",
		})
	}
	res, errData := i.service.NewItemGroup(payload)
	if errData != nil {
		helper.ReturnError(writer, request, errData)
		return
	}
	payloads.NewHandleSuccess(writer, res, "create Data Successfully", 201)
}

func (i *ItemGroupControllerImpl) GetItemGroupByCode(writer http.ResponseWriter, request *http.Request) {
	code := chi.URLParam(request, "item_group_code")
	res, errData := i.service.GetItemGroupByCode(code)
	if errData != nil {
		helper.ReturnError(writer, request, errData)
		return
	}
	payloads.NewHandleSuccess(writer, res, "Get Data Successfully", 200)
}

func NewItemGroupControllerImpl(service masteritemservice.ItemGroupService) ItemGroupController {
	return &ItemGroupControllerImpl{service: service}
}
