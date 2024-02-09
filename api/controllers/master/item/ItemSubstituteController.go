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

	"github.com/julienschmidt/httprouter"
)

type ItemSubstituteController interface {
	GetAllItemSubstitute(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetByIdItemSubstitute(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllItemSubstituteDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetByIdItemSubstituteDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveItemSubstitute(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SaveItemSubstituteDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangeStatusOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}


type ItemSubstituteControllerImpl struct {
	ItemSubstituteService masteritemservice.ItemSubstituteService
}

func NewItemSubstituteController(itemSubstituteService masteritemservice.ItemSubstituteService) ItemSubstituteController {
	return &ItemSubstituteControllerImpl{
		ItemSubstituteService: itemSubstituteService,
	}
}

func (r *ItemSubstituteControllerImpl) GetAllItemSubstitute(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"is_active":            queryValues.Get("is_active"),
		"substitute_type_code": queryValues.Get("substitute_type_code"),
		"effective_date":       queryValues.Get("effective_date"),
		"item_id":              queryValues.Get("item_id"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.ItemSubstituteService.GetAllItemSubstitute(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *ItemSubstituteControllerImpl) GetByIdItemSubstitute(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ItemSubstituteIdStr := params.ByName("item_substitute_id")

	ItemSubstituteId, _ := strconv.Atoi(ItemSubstituteIdStr)

	result := r.ItemSubstituteService.GetByIdItemSubstitute(ItemSubstituteId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemSubstituteControllerImpl) GetAllItemSubstituteDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"item_substitute_id": queryValues.Get("item_substitute_id"),
	}
	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result := r.ItemSubstituteService.GetAllItemSubstitute(filterCondition, pagination)

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *ItemSubstituteControllerImpl) GetByIdItemSubstituteDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ItemSubstituteDetailIdStr := params.ByName("item_substitute_detail_id")

	ItemSubstituteDetailId, _ := strconv.Atoi(ItemSubstituteDetailIdStr)

	result := r.ItemSubstituteService.GetByIdItemSubstitute(ItemSubstituteDetailId)

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemSubstituteControllerImpl) SaveItemSubstitute(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var formRequest masteritempayloads.ItemSubstitutePayloads
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.ItemSubstituteService.SaveItemSubstitute(formRequest)

	if formRequest.ItemSubstituteId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *ItemSubstituteControllerImpl) SaveItemSubstituteDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var formRequest masteritempayloads.ItemSubstituteDetailPayloads
	helper.ReadFromRequestBody(request, &formRequest)
	var message = ""

	create := r.ItemSubstituteService.SaveItemSubstituteDetail(formRequest)

	if formRequest.ItemSubstituteId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *ItemSubstituteControllerImpl) ChangeStatusOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	ItemSubstituteId, _ := strconv.Atoi(params.ByName("item_substitute_id"))

	response := r.ItemSubstituteService.ChangeStatusItemOperation(int(ItemSubstituteId))

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// func (r *ItemSubstituteControllerImpl) ActivateItemSubstituteDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

// 	ItemSubstituteDetailId, _ := strconv.Atoi(params.ByName("item_substitute_detail_id"))

// 	response := r.ItemSubstituteService.ActivateItemSubstituteDetail(ItemSubstituteDetailId)

// 	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
// }