package mastercontroller

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ItemOperationController interface {
	GetAllItemOperation(writer http.ResponseWriter, request *http.Request)
	GetByIdItemOperation(writer http.ResponseWriter, request *http.Request)
	PostItemOperation(writer http.ResponseWriter, request *http.Request)
	DeleteItemOperation(writer http.ResponseWriter, request *http.Request)
	UpdateItemOperation(writer http.ResponseWriter, request *http.Request)
}

type ItemOperationControllerImpl struct {
	ItemOperationService masterservice.ItemOperationService
}

func NewItemOperationController(ItemOperationService masterservice.ItemOperationService) ItemOperationController {
	return &ItemOperationControllerImpl{
		ItemOperationService: ItemOperationService,
	}
}

func (r *ItemOperationControllerImpl) GetAllItemOperation(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	queryParams := map[string]string{
		"mtr_item_operation.item_operation_id": queryValues.Get("item_operation_id"),
		"mtr_item_operation.item_id":           queryValues.Get("item_id"),
		"mtr_item_operation.operation_id":      queryValues.Get("operation_id"),
		"mtr_item_operation.line_type_id":      queryValues.Get("line_type_id"),
	}

	paginate := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(queryValues, "limit"),
		Page:   utils.NewGetQueryInt(queryValues, "page"),
		SortOf: queryValues.Get("sort_of"),
		SortBy: queryValues.Get("sort_by"),
	}

	criteria := utils.BuildFilterCondition(queryParams)

	result,err := r.ItemOperationService.GetAllItemOperation(criteria,paginate)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer,result.Rows,"Get data successfully",200,result.Limit,result.Page,result.TotalRows,result.TotalPages)
}

func (r *ItemOperationControllerImpl) GetByIdItemOperation(writer http.ResponseWriter, request *http.Request){
	itemClassId, _ := strconv.Atoi(chi.URLParam(request, "item_operation_id"))

	result,err := r.ItemOperationService.GetByIdItemOperation(itemClassId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Get Data Successfully!", http.StatusOK)
}

func (r *ItemOperationControllerImpl) PostItemOperation(writer http.ResponseWriter, request *http.Request){
	var formRequest masterpayloads.ItemOperationPost
	helper.ReadFromRequestBody(request, &formRequest)
	create, err := r.ItemOperationService.PostItemOperation(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, create, "success create data", http.StatusOK)
}

func (r *ItemOperationControllerImpl) DeleteItemOperation(writer http.ResponseWriter, request *http.Request){
	itemoperationid,_ := strconv.Atoi(chi.URLParam(request,"item_operation_id"))
	delete,err := r.ItemOperationService.DeleteItemOperation(itemoperationid)
	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, delete, "success delete data", http.StatusOK)
}

func (r *ItemOperationControllerImpl) UpdateItemOperation(writer http.ResponseWriter, request *http.Request){
	var formRequest masterpayloads.ItemOperationPost
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	if err !=nil{
		helper.ReturnError(writer, request, err)
		return
	}
	itemoperationid,_ := strconv.Atoi(chi.URLParam(request,"item_operation_id"))

	update,err2 := r.ItemOperationService.UpdateItemOperation(itemoperationid,formRequest)
	if err2 != nil {
		helper.ReturnError(writer, request, err2)
		return
	}

	payloads.NewHandleSuccess(writer, update, "success update data", http.StatusOK)
}